package p2p

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"github.com/panzerox123/blockcert/src/certificate"
)

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
		fmt.Println("Parse Address:", m_addr)
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
	node_p2p.blockchainTopic, err = node_p2p.pubsub.Join("Blockchain")
	if err != nil {
		panic(err)
	}
	node_p2p.blockchain = certificate.NewBlockChain()
	node_p2p.BlockListener(ctx)

	return &node_p2p
}

func (node_p2p *P2pNode) BlockListener(ctx context.Context) {
	subscription, err := node_p2p.blockchainTopic.Subscribe()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			msg, err := subscription.Next(ctx)
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
	return node.blockchain.ChainValid()
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
