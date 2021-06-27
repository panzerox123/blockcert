package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/panzerox123/blockcert/src/certificate"
	"github.com/panzerox123/blockcert/src/keygen"
	"github.com/panzerox123/blockcert/src/p2p"
)

func generateKeys(filename string) {
	private_key, public_key := keygen.GenerateKeyPair(512)
	keygen.SaveHexKey(filename, private_key, public_key)
}

func sign(prikey string, pubkey string) {
	privateKey := keygen.ParsePrivateRSA(prikey)
	publicKey := keygen.ParsePublicRSA(pubkey)
	blockchain := certificate.NewBlockChain()
	blockchain.AddBlock("Hello World!", privateKey)
	println(blockchain.ChainValid())
	blockchain.AddBlock("My Name is Kunal Bhat!", privateKey)
	println(blockchain.ChainValid())
	blockchain.AddBlock("This is the third block", privateKey)
	println(blockchain.ChainValid())
	println("Checks")
	println(blockchain.CheckSignature("Hello World!", publicKey))
	println(blockchain.CheckSignature("My Name is Kunal Bhat!", publicKey))
	println(blockchain.CheckSignature("This is the third block", publicKey))
	fmt.Printf("%v", blockchain.GetLatest())
}

func shell(ctx context.Context, node *p2p.P2pNode) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("->")
		cli, _ := reader.ReadString('\n')
		cli = strings.Replace(cli, "\n", "", -1)
		cli_args := strings.Split(cli, " ")
		switch cli_args[0] {
		case "addcert":
			privateKey := keygen.ParsePrivateRSA(cli_args[2])
			node.AddBlock(ctx, cli_args[1], privateKey)
			fmt.Println("Certificate successfully added!")
		case "showallcerts":
			node.ShowBlocks()
		case "verifyallcerts":
			ret := node.VerifyChain()
			if ret {
				fmt.Println("Blockchain VALID!")
			} else {
				fmt.Println("Blockchain INVALID! Rebuilding!")
			}
		case "checkcert":
			pubKey := keygen.ParsePublicRSA(cli_args[2])
			ret := node.CheckCertificate(cli_args[1], pubKey)
			if ret {
				fmt.Println("Certificate VERIFIED!")
			} else {
				fmt.Println("Certificate INVALID! Please try again!")
			}
		}
	}
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "keygen":
			if len(os.Args) < 2 {
				generateKeys("certificate.key")
				return
			}
			generateKeys(os.Args[2])
		case "sign":
			sign("3082013c020100024100d49ec03ffdb560e7f6fa16d65d2472b74ceeec96940f06ae3b8d060c16d58ae512478de038cf05754ae5bb51d29c4b6c14fbf4a5bb838a5d42d59a39b21d03bf0203010001024100bc1a66833675ccf1eb727dd9d0357ab7e7fc489b3f09bc2350d406d1933200d9e36e896c9a1c33e79d004320e29ad187ba4b085d69085ed13643fad664309001022100e9bcbc9a9cc0103bb83ff4774d5d6cd2fe2dcb09ed0a7524649906c524fd3201022100e8df1c6399b93d28f79e17f2f4be10ba225370ba83679dd7b9f4835f80dab5bf02203322e99861e6db26559f185ae9802108e037208ea15f82555df4e4b848e9640102210095e52f77e9367468cf62e3158f765c7c03a664149a8af2ee2e937690ddf76a2f022100d88ba95c4b998a8947134a6fd3f0420113376fdfffba6008871784ceb61818c8",
				"3048024100d49ec03ffdb560e7f6fa16d65d2472b74ceeec96940f06ae3b8d060c16d58ae512478de038cf05754ae5bb51d29c4b6c14fbf4a5bb838a5d42d59a39b21d03bf0203010001")
		case "node":
			ctx := context.Background()
			var node *p2p.P2pNode
			if len(os.Args) > 2 {
				node = p2p.NewP2pNode(ctx, os.Args[2])
			} else {
				node = p2p.NewP2pNode(ctx, "")
			}
			go shell(ctx, node)
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			<-ch
			fmt.Printf("Shutting down...")
			node.CloseNode()

		}
	}
}
