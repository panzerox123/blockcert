package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/panzerox123/blockcert/src/keygen"
	"github.com/panzerox123/blockcert/src/p2p"
)

var node *p2p.P2pNode

func keygenHandler(res http.ResponseWriter, req *http.Request) {
	private_key, public_key := keygen.GenerateKeyPair(512)
	new_keypair := Keys{
		PrivateHex: keygen.EncodePrivateRSA(private_key),
		PublicHex:  keygen.EncodePublicRSA(public_key),
	}
	err := json.NewEncoder(res).Encode(new_keypair)
	if err != nil {
		panic(err)
	}
}

func httpRequestHandler(PORT int) {
	router := mux.NewRouter().StrictSlash(true)
	http.HandleFunc("/generate_keys", keygenHandler)
	fmt.Println("SERVER STARTED ON PORT", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
	if err != nil {
		panic(err)
	}
}

func StartServer(PORT int, NODE *p2p.P2pNode) {
	node = NODE
	httpRequestHandler(PORT)
}
