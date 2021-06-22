package main

import (
	"fmt"

	"github.com/panzerox123/blockcert/src/blockchain"
)

func main() {
	chain := blockchain.NewDigitalCertChain("Kunal")
	chain.AddCert(0, "My name is kunal bhat. Hopefully you can read my data fine")
	fmt.Println(chain.GetLatest().Data)
	fmt.Println(chain.CheckValid())
	chain.AddCert(1, "This is my second piece of data")
	fmt.Println(chain.GetLatest().Data)
	fmt.Println(chain.CheckValid())

}
