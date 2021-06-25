package p2p

import (
	"context"
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
	node_p2p.blockchain = certificate.NewBlockChain()

	return &node_p2p
}

func (node_p2p *P2pNode) BlockChainListener(ctx context.Context) {
	topic, err := node_p2p.pubsub.Join("Blockchain")
	if err != nil {
		panic(err)
	}
	subscription, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			msg, err := subscription.Next(ctx)
			if err != nil {
				panic(err)
			}
			_ = msg.GetData()
		}
	}()
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
