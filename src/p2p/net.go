package p2p

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
)

func StartNode() {
	ctx := context.Background()
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/8080"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Listen Addresses: ", node.Addrs())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Recieved shutdown signal... shutting down")
	if err := node.Close(); err != nil {
		panic(err)
	}
}
