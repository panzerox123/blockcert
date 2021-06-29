package p2p

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	discovery "github.com/libp2p/go-libp2p-discovery"
	kdht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	disc "github.com/libp2p/go-libp2p/p2p/discovery"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"

	"github.com/multiformats/go-multiaddr"
	"github.com/panzerox123/blockcert/src/certificate"
)

var DISABLE_DISCOVERY bool = true

func NewP2pNode(ctx context.Context, addrstr string) *P2pNode {
	var node_p2p P2pNode
	node, err := libp2p.New(
		ctx,
		libp2p.Defaults,
	)
	if err != nil {
		panic(err)
	}
	node_p2p.pubsub, err = pubsub.NewFloodSub(ctx, node)
	if err != nil {
		panic(err)
	}
	fmt.Println("Available interfaces:")
	for i, addr := range node.Addrs() {
		fmt.Printf("%d: %s/ipfs/%s\n", i, addr, node.ID().Pretty())
	}
	if addrstr != "" {
		m_addr, err := multiaddr.NewMultiaddr(addrstr)
		if err != nil {
			panic(err)
		}
		fmt.Println("Connected to peer:", m_addr)
		peer_info, err := peerstore.AddrInfoFromP2pAddr(m_addr)
		if err != nil {
			panic(err)
		}
		if err := node.Connect(ctx, *peer_info); err != nil {
			fmt.Println("Could not connect to peer!")
			panic(err)
		}
	}
	node_p2p.node = node
	if !DISABLE_DISCOVERY {
		node_p2p.peerDiscovery(ctx)
	}
	node_p2p.blockchainTopic, err = node_p2p.pubsub.Join("Blockchain")
	if err != nil {
		panic(err)
	}
	node_p2p.blockchainSubscription, err = node_p2p.blockchainTopic.Subscribe()
	if err != nil {
		panic(err)
	}
	node_p2p.blockchain = certificate.NewBlockChain()
	node_p2p.BlockListener(ctx)

	return &node_p2p
}

type discoveryNotifee struct {
	h host.Host
	c context.Context
}

func (d *discoveryNotifee) HandlePeerFound(pi peerstore.AddrInfo) {
	fmt.Printf("Discovered new peer: %s\n", pi.ID.Pretty())
	err := d.h.Connect(d.c, pi)
	if err != nil {
		fmt.Printf("!! Could not connect to peer : %s\n", pi.ID.Pretty())
	}
}

func (node_p2p *P2pNode) LocalPeerDiscovery(ctx context.Context) {
	serv, err := disc.NewMdnsService(ctx, node_p2p.node, time.Hour, "blockchain_pubsub")
	if err != nil {
		fmt.Println(err)
		return
	}
	n := discoveryNotifee{h: node_p2p.node, c: ctx}
	serv.RegisterNotifee(&n)
}

func (node_p2p *P2pNode) peerDiscovery(ctx context.Context) {
	kaddht, err := kdht.New(ctx, node_p2p.node)
	if err != nil {
		println(err)
	}
	if err = kaddht.Bootstrap(ctx); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range kdht.DefaultBootstrapPeers {
		peerinfo, _ := peerstore.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node_p2p.node.Connect(ctx, *peerinfo); err != nil {
				fmt.Println("(Bootstrap)!! Could not connect to peer: ", peerinfo.ID)
			} else {
				fmt.Println("(Bootstrap)Connected to peer: ", peerinfo.ID)
			}
		}()
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(kaddht)
	discovery.Advertise(ctx, routingDiscovery, "peer_discovery")
	peers, err := routingDiscovery.FindPeers(ctx, "peer_discovery")
	if err != nil {
		panic(err)
	}
	for peer := range peers {
		if peer.ID == node_p2p.node.ID() {
			continue
		} else {
			err := node_p2p.node.Connect(ctx, peer)
			if err != nil {
				fmt.Println("!! Could not connect to peer: ", peer.ID)
			}
			fmt.Println("Connected to peer:", peer.ID)
		}
	}
}

func (node_p2p *P2pNode) ReturnPeerList() {
	val := node_p2p.node.Peerstore()
	fmt.Println(val.PeersWithAddrs())
}

func (node_p2p *P2pNode) BlockListener(ctx context.Context) {

	go func() {
		for {
			msg, err := node_p2p.blockchainSubscription.Next(ctx)
			if err != nil {
				panic(err)
			}
			temp_bc := certificate.NewBlockChain()
			err = json.Unmarshal(msg.GetData(), temp_bc)
			if err != nil {
				panic(err)
			}
			node_p2p.blockchain = temp_bc
		}
	}()
}

func (node_p2p *P2pNode) AddBlock(ctx context.Context, data string, prikey *rsa.PrivateKey) {
	node_p2p.blockchain.AddBlock(data, prikey)
	node_p2p.BlockPublisher(ctx)
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

func (node *P2pNode) CheckCertificate(data string, pubkey *rsa.PublicKey) bool {
	return node.blockchain.CheckSignature(data, pubkey)
}

func (node_p2p *P2pNode) BlockPublisher(ctx context.Context) {
	jsoned_bc, err := json.Marshal(node_p2p.blockchain)
	if err != nil {
		panic(err)
	}
	err = node_p2p.blockchainTopic.Publish(ctx, jsoned_bc)
	if err != nil {
		panic(err)
	}
}

func (node_p2p *P2pNode) CloseNode() {
	err := node_p2p.node.Close()
	if err != nil {
		panic(err)
	}
}

// TODO: Remove below functions

func StartNode(ctx context.Context) *host.Host {
	node, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}
	return &node
}

func ConnectToNode(ctx context.Context, node *host.Host, host string) {
	ping_service := &ping.PingService{
		Host: *node,
	}
	(*node).SetStreamHandler(ping.ID, ping_service.PingHandler)

	addr, err := multiaddr.NewMultiaddr(host)
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := (*node).Connect(ctx, *peer); err != nil {
		panic(err)
	}
	ch := ping_service.Ping(ctx, peer.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("Pinged ", host, "in", res.RTT)
	}
}

func NodeInfo(node *host.Host) []multiaddr.Multiaddr {
	peer_info := &peerstore.AddrInfo{
		ID:    (*node).ID(),
		Addrs: (*node).Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(peer_info)
	if err != nil {
		panic(err)
	}
	return addrs
}

func CloseNode(node *host.Host) {
	if err := (*node).Close(); err != nil {
		panic(err)
	}
}
