package api

type keygenStruct struct {
	PublicHex  string `json:"PublicKey"`
	PrivateHex string `json:"PrivateKey"`
}

type newCertStruct struct {
	PrivateKey string `json:"PrivateKey"`
	Data       []byte `json:"Data"`
}
