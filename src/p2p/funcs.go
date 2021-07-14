package p2p

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	kdht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/multiformats/go-multiaddr"
	"github.com/panzerox123/blockcert/src/certificate"
	"github.com/panzerox123/blockcert/src/keygen"
)

func RAND_FUNC() int {
	return 2 + rand.Intn(4)
}

func NewP2pNode(ctx context.Context, addrstr string) *P2pNode {
	var node_p2p P2pNode
	node, err := libp2p.New(
		ctx,
		libp2p.Defaults,
		libp2p.NATPortMap(),
	)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error creating new node:", err.Error())
		return nil
	}
	node_p2p.pubsub, err = pubsub.NewFloodSub(ctx, node)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error creating new PUB/SUB:", err.Error())
		return nil
	}
	for _, addr := range node.Addrs() {
		node_p2p.interfaces = append(node_p2p.interfaces, fmt.Sprintf("%s/ipfs/%s\n", addr, node.ID().Pretty()))
	}
	if addrstr != "" {
		m_addr, err := multiaddr.NewMultiaddr(addrstr)
		if err != nil {
			fmt.Println(Red+"[❌]"+Reset, "Error converting multiaddr:", err.Error())
			return nil
		}
		fmt.Println("Connected to peer:", m_addr)
		peer_info, err := peerstore.AddrInfoFromP2pAddr(m_addr)
		if err != nil {
			fmt.Println(Red+"[❌]"+Reset, "Error creating converting to AddrInfo:", err.Error())
			return nil
		}
		if err := node.Connect(ctx, *peer_info); err != nil {
			fmt.Println(Red+"[❌]"+Reset, "Could not connect to given peer!", err.Error())
			return nil
		}
	}
	node_p2p.node = node
	if !DISABLE_DISCOVERY {
		node_p2p.peerDiscovery(ctx)
	}
	node_p2p.blockchainTopic, err = node_p2p.pubsub.Join("Blockchain")
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error joining topic \"Blockchain\":", err.Error())
		return nil
	}
	node_p2p.blockchainSubscription, err = node_p2p.blockchainTopic.Subscribe()
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error subscribing to topic \"Blockchain\":", err.Error())
		return nil
	}
	node_p2p.newcertTopic, err = node_p2p.pubsub.Join("Newcert")
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error joining topic \"Newcert\":", err.Error())
		return nil
	}
	node_p2p.newcertSubscription, err = node_p2p.newcertTopic.Subscribe()
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error subscribing to topic \"Newcert\":", err.Error())
		return nil
	}
	//node_p2p.blockchain = certificate.NewBlockChain()
	node_p2p.blockchain = certificate.ReadBlockChain()
	node_p2p.blockListener(ctx)
	node_p2p.blockPublisher(ctx)
	node_p2p.newCertListener(ctx)
	node_p2p.peerDiscoveryTimed(ctx)
	return &node_p2p
}

func (node_p2p *P2pNode) PrintInterfaces() {
	for i, x := range node_p2p.interfaces {
		fmt.Println(Cyan+"[interface", i, "\b]"+Reset, x)
	}
}

func (node_p2p *P2pNode) peerDiscoveryTimed(ctx context.Context) {
	go func() {
		for {
			time.Sleep(30 * time.Minute)
			fmt.Println(Blue+"[⌛]"+Reset, "Discovering new peers!")
			node_p2p.peerDiscovery(ctx)
			node_p2p.blockPublisher(ctx)
		}
	}()
}

func (node_p2p *P2pNode) peerDiscovery(ctx context.Context) {
	kaddht, err := kdht.New(ctx, node_p2p.node)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error creating new DHT:", err.Error())
	}
	if err = kaddht.Bootstrap(ctx); err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error bootstrapping peers:", err.Error())
	}
	var wg sync.WaitGroup
	for _, peerAddr := range kdht.DefaultBootstrapPeers {
		peerinfo, _ := peerstore.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node_p2p.node.Connect(ctx, *peerinfo); err != nil {
				fmt.Println(Yellow+"[!]"+Reset, "[BOOTSTRAP]Could NOT connect to peer: ", peerinfo.ID)
			} else {
				fmt.Println(Green+"[✓]"+Reset, "[BOOTSTRAP]Connected to peer:", peerinfo.ID)
			}
		}()
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(kaddht)
	discovery.Advertise(ctx, routingDiscovery, "peer_discovery")
	peers, err := routingDiscovery.FindPeers(ctx, "peer_discovery")
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error finding peers:", err.Error())
		return
	}
	for peer := range peers {
		if peer.ID == node_p2p.node.ID() {
			continue
		} else {
			err := node_p2p.node.Connect(ctx, peer)
			if err != nil {
				fmt.Println(Yellow+"[!]"+Reset, "Could NOT connect to peer:", peer.ID)
			} else {
				fmt.Println(Green+"[✓]"+Reset, "Connected to peer:", peer.ID)
			}
		}
	}
}

func (node_p2p *P2pNode) ReturnPeerList() {
	val := node_p2p.node.Peerstore()
	fmt.Println(val.PeersWithAddrs())
}

func (node_p2p *P2pNode) blockListener(ctx context.Context) {

	go func() {
		for {
			msg, err := node_p2p.blockchainSubscription.Next(ctx)
			if err != nil {
				fmt.Println(Red+"[❌]"+Reset, "Error getting new data:", err.Error())
				continue
			}
			temp_bc := certificate.NewBlockChain()
			err = json.Unmarshal(msg.GetData(), temp_bc)
			if err != nil {
				fmt.Println(Red+"[❌]"+Reset, "Error decoding JSON data:", err.Error())
				continue
			}
			if !temp_bc.ChainValid() {
				continue
			} else {
				node_p2p.LockNet.Lock()
				if !node_p2p.blockchain.ChainValid() {
					node_p2p.blockchain = temp_bc
				} else if len(temp_bc.Chain) > len(node_p2p.blockchain.Chain) {
					if node_p2p.blockchain.CompareChains(temp_bc) {
						node_p2p.blockchain = temp_bc
					} else {
						node_p2p.blockPublisher(ctx)
					}
				} else if len(temp_bc.Chain) == len(node_p2p.blockchain.Chain) && node_p2p.blockchain.CompareChains(temp_bc) {
					temp_latest := temp_bc.GetLatest()
					curr_latest := node_p2p.blockchain.GetLatest()
					if temp_latest != nil && curr_latest != nil && temp_latest.Proof > curr_latest.Proof {
						node_p2p.blockchain = temp_bc
					} else if temp_latest != nil && curr_latest != nil && temp_latest.Proof == curr_latest.Proof {
						node_p2p.LockNet.Unlock()
						continue
					} else {
						if temp_latest != nil {
							node_p2p.blockPublisher(ctx)
						}
					}
				} else {
					node_p2p.blockPublisher(ctx)
				}
				node_p2p.LockNet.Unlock()
			}
			node_p2p.blockchain.SaveBlockchainJson()
		}
	}()
}

func (node_p2p *P2pNode) NewCertPublisher(ctx context.Context, data []byte, private_key string) error {
	fmt.Println(Blue+"[📢]", "Publishing new data!", Reset)
	cert_info := NewCertPublish{
		Data:       data,
		PrivateKey: private_key,
	}
	jsoned_data, err := json.Marshal(cert_info)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error Reading JSON data:", err.Error())
		return fmt.Errorf(Red+"[❌]"+Reset, "Error Reading JSON data:", err.Error())
	}
	err = node_p2p.newcertTopic.Publish(ctx, jsoned_data)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error Publishing new data:", err.Error())
		return fmt.Errorf(Red+"[❌]"+Reset, "Error Publishing new data:", err.Error())
	}
	fmt.Println(Green+"[📢]", "Published new data!", Reset)
	return nil
}

func (node_p2p *P2pNode) newCertListener(ctx context.Context) {
	go func() {
		for {
			msg, err := node_p2p.newcertSubscription.Next(ctx)
			if err != nil {
				fmt.Println(Red+"[❌]"+Reset, "Error getting new certificate data:", err.Error())
				continue
			}
			var temp_cert NewCertPublish
			err = json.Unmarshal(msg.GetData(), &temp_cert)
			if err != nil {
				fmt.Println(Red+"[❌]"+Reset, "Error Reading JSON:", err.Error())
				continue
			}
			priv_key := keygen.ParsePrivateRSA(temp_cert.PrivateKey)
			node_p2p.LockNet.Lock()
			fmt.Println(Blue+"[💻]", "Mining newly recieved data!", Reset)
			if node_p2p.blockchain.CheckDataExists(temp_cert.Data) {
				fmt.Println(Green+"[💻]", "Data exists! Not mining!", Reset)
				node_p2p.LockNet.Unlock()
				continue
			}
			node_p2p.blockchain.AddBlock(temp_cert.Data, priv_key, RAND_FUNC())
			node_p2p.blockPublisher(ctx)
			fmt.Println(Green+"[💻]", "Block mined!", Reset)
			node_p2p.LockNet.Unlock()
		}
	}()
}

func (node_p2p *P2pNode) ShowBlocks() {
	for i, x := range node_p2p.blockchain.Chain {
		fmt.Printf("%d : %v\n", i, x)
	}
}

func (node *P2pNode) VerifyChain() bool {
	val := node.blockchain.ChainValid()
	fmt.Println(val)
	return val
}

func (node *P2pNode) CheckCertificate(data []byte, pubkey *rsa.PublicKey) bool {
	return node.blockchain.CheckSignature(data, pubkey)
}

func (node_p2p *P2pNode) blockPublisher(ctx context.Context) {
	fmt.Println(Blue+"[📢]", "Publishing blocks!", Reset)
	if !node_p2p.blockchain.ChainValid() {
		return
	}
	jsoned_bc, err := json.Marshal(node_p2p.blockchain)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error Writing JSON:", err.Error())
		return
	}
	err = node_p2p.blockchainTopic.Publish(ctx, jsoned_bc)
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error Publishing new data:", err.Error())
		return
	}
	fmt.Println(Green+"[📢]", "Blocks published!", Reset)
}

func (node_p2p *P2pNode) CloseNode() {
	err := node_p2p.node.Close()
	if err != nil {
		fmt.Println(Red+"[❌]"+Reset, "Error shutting down node:", err.Error())
	}
}
