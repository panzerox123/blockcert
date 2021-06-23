package main

import (
	"fmt"
	"os"

	"github.com/panzerox123/blockcert/src/keygen"
)

func generateKeys() {
	private_key, public_key := keygen.GenerateKeyPair(512)
	keygen.SaveHexKey("private_public.key", private_key, public_key)
}

func testSign() {
	prikey := "3082013b020100024100a1072159047761c2a3abd63b399e1cda7ba2272dc190185dca97ab82cd4fb1546df980c42fbe3b67fdd36d515e4bb11ea6dcc2cacb8977adc7d86efcad0f71a1020301000102402a00381c85e3b5a61516cf0c279d2c1d78bdf4c62484b7364f8f7bf6e42273380e7f7a78fd6c6138f58f5d4a41fac69e48b77a43707cb628f3e1638be3f47001022100cf0b679e4c2970710424e5b0e5bec2443f91d37edd63bc8898e3fbdbd9fe7ee1022100c71a4afb2dcfd1906868d94074a42ab61c90e354fe96baf876e500c83bdd0ac10221009973ad77b0a121fa5184fb4c31eb41568dfb09d2c4495089b92f7812c92e0b61022100a7c4c0ffcc077c87996318055703ea358ff68a88590a3bbc17bb39a07fc8ef4102202cd096241d9d7a0409d75f1b9b2e24c98d38ab7da5b8fef78c9f0e50c648a62e"
	pubkey := "3048024100a1072159047761c2a3abd63b399e1cda7ba2272dc190185dca97ab82cd4fb1546df980c42fbe3b67fdd36d515e4bb11ea6dcc2cacb8977adc7d86efcad0f71a10203010001"
	privateKey := keygen.ParsePrivateRSA(prikey)
	publicKey := keygen.ParsePublicRSA(pubkey)
	sign := keygen.SignData("Helo world!", privateKey)
	fmt.Println(keygen.VerifyData("Helo world!", sign, publicKey))
}

func main() {
	switch os.Args[1] {
	case "keygen":
		generateKeys()
	case "sign":
		testSign()
	}
}
