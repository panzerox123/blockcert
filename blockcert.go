package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/panzerox123/blockcert/src/api"
	"github.com/panzerox123/blockcert/src/certificate"
	"github.com/panzerox123/blockcert/src/keygen"
	"github.com/panzerox123/blockcert/src/p2p"
)

func generateKeys(filename string) {
	private_key, public_key := keygen.GenerateKeyPair(512)
	keygen.SaveHexKey(filename, private_key, public_key)
}

func shell(ctx context.Context, node *p2p.P2pNode) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\033[35mDEBUG>\033[0m ")
		cli, _ := reader.ReadString('\n')
		cli = strings.Replace(cli, "\n", "", -1)
		cli = strings.Replace(cli, "\r", "", -1)
		cli_args := strings.Split(cli, " ")
		switch cli_args[0] {
		case "addcert":
			//privateKey := keygen.ParsePrivateRSA(cli_args[2])
			data := certificate.FileByteOut(cli_args[1])
			node.NewCertPublisher(ctx, data, cli_args[2])
			fmt.Println("Certificate successfully added!")
		case "showallcerts":
			node.ShowBlocks()
		case "peers":
			node.ReturnPeerList()
		case "verifyallcerts":
			ret := node.VerifyChain()
			if ret {
				fmt.Println("Blockchain VALID!")
			} else {
				fmt.Println("Blockchain INVALID! Rebuilding!")
			}
		case "interfaces":
			node.PrintInterfaces()
		case "checkcert":
			pubKey := keygen.ParsePublicRSA(cli_args[2])
			ret := node.CheckCertificate(certificate.FileByteOut(cli_args[1]), pubKey)
			if ret {
				fmt.Println("Certificate VERIFIED!")
			} else {
				fmt.Println("Certificate INVALID! Please try again!")
			}
		}
	}
}

func main() {
	startDebug := flag.Bool("debug", false, "Start the debug shell")
	generateKey := flag.Bool("keygen", false, "Generate a Private and Public RSA Key")
	key_output := flag.String("o", "certificate.key", "Output for the generated keys!")
	flag.Parse()

	if *generateKey {
		generateKeys(*key_output)
		os.Exit(0)
	}
	ctx := context.Background()
	node := p2p.NewP2pNode(ctx, "")
	if *startDebug {
		go shell(ctx, node)
	} else {
		go api.StartServer(8080, node)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Printf("Shutting down...")
	node.CloseNode()
}
