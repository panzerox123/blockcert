package certificate

type Certificate struct {
	TimeStamp int64
	FileHash  string
	Signature string
}

type Block struct {
	Data     Certificate
	Hash     string
	PrevHash string
	Proof    int64
}

type BlockChain struct {
	Chain []Block
}
