package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/panzerox123/blockcert/src/keygen"
	"github.com/panzerox123/blockcert/src/p2p"
)

var node *p2p.P2pNode

func checkCertHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(1024 << 20)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("Error Retrieving file D:"))
		fmt.Println(err.Error())
		return
	}
	file, _, err := req.FormFile("Data")
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("Error Reading file"))
		fmt.Println(err.Error())
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	pubkey := keygen.ParsePublicRSA(req.FormValue("PublicKey"))
	ret := checkCertStruct{
		node.CheckCertificate(fileBytes, pubkey),
	}
	res.WriteHeader(200)
	err = json.NewEncoder(res).Encode(ret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func newCertHandler(res http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(1024 << 20)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("Error Retrieving file D:"))
		fmt.Println("[FRONTEND_SERVER]", err.Error())
		return
	}
	file, _, err := req.FormFile("Data")
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("[FRONTEND_SERVER]Error Reading file"))
		fmt.Println(err.Error())
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	node.NewCertPublisher(context.Background(), fileBytes, req.FormValue("PrivateKey"))
	res.WriteHeader(200)
	res.Write([]byte("Success! :D"))
}

func keygenHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Generating Keys!")
	private_key, public_key := keygen.GenerateKeyPair(512)
	new_keypair := keygenStruct{
		PrivateHex: keygen.EncodePrivateRSA(private_key),
		PublicHex:  keygen.EncodePublicRSA(public_key),
	}
	err := json.NewEncoder(res).Encode(new_keypair)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("Could not generate your keypair D:"))
	}
}

func statusRequestHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Returning status")
	err := json.NewEncoder(res).Encode(node.Status.Status)
	if err != nil {
		res.WriteHeader(500)
		res.Write([]byte("Error returning Status"))
	}
}

func httpRequestHandler(PORT int) {
	http.HandleFunc("/keygen", keygenHandler)
	http.HandleFunc("/new_cert", newCertHandler)
	http.HandleFunc("/check_cert", checkCertHandler)
	http.HandleFunc("/status", statusRequestHandler)
	fmt.Println("SERVER STARTED ON PORT", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		panic(err)
	}
}

func StartServer(PORT int, NODE *p2p.P2pNode) {
	node = NODE
	httpRequestHandler(PORT)
}
