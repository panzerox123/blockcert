package main

import (
	"os"

	"github.com/panzerox123/blockcert/src/certificate"
	"github.com/panzerox123/blockcert/src/keygen"
)

func generateKeys() {
	private_key, public_key := keygen.GenerateKeyPair(512)
	keygen.SaveHexKey("private_public.key", private_key, public_key)
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
}

func main() {
	switch os.Args[1] {
	case "keygen":
		generateKeys()
	case "sign":
		sign("3082013c020100024100d49ec03ffdb560e7f6fa16d65d2472b74ceeec96940f06ae3b8d060c16d58ae512478de038cf05754ae5bb51d29c4b6c14fbf4a5bb838a5d42d59a39b21d03bf0203010001024100bc1a66833675ccf1eb727dd9d0357ab7e7fc489b3f09bc2350d406d1933200d9e36e896c9a1c33e79d004320e29ad187ba4b085d69085ed13643fad664309001022100e9bcbc9a9cc0103bb83ff4774d5d6cd2fe2dcb09ed0a7524649906c524fd3201022100e8df1c6399b93d28f79e17f2f4be10ba225370ba83679dd7b9f4835f80dab5bf02203322e99861e6db26559f185ae9802108e037208ea15f82555df4e4b848e9640102210095e52f77e9367468cf62e3158f765c7c03a664149a8af2ee2e937690ddf76a2f022100d88ba95c4b998a8947134a6fd3f0420113376fdfffba6008871784ceb61818c8",
			"3048024100d49ec03ffdb560e7f6fa16d65d2472b74ceeec96940f06ae3b8d060c16d58ae512478de038cf05754ae5bb51d29c4b6c14fbf4a5bb838a5d42d59a39b21d03bf0203010001")
	}
}
