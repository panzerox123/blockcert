package api

type keygenStruct struct {
	PublicHex  string `json:"PublicKey"`
	PrivateHex string `json:"PrivateKey"`
}

type checkCertStruct struct {
	Verification bool `json:"Verification"`
}

/*
type newCertStruct struct {
	PrivateKey string `json:"PrivateKey"`
	Data       []byte `json:"Data"`
}
*/
