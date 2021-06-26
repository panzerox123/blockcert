package certificate

import "time"

type Certificate struct {
	TimeStamp time.Time
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
