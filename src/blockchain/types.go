package blockchain

type DigitalCert struct {
	CertificateVersion int64    // Certificate Version
	Data               string   // The Actual Data
	Hash               [32]byte // Hash of current block
	PrevHash           [32]byte // Hash of previous Block
	TimeStamp          int64    // Unix Time Stamp value
	ProofWork          int64    // Store Proof of work value
}

type DigitalCertChain struct {
	Chain []DigitalCert
}
