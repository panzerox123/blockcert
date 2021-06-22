package certificate

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

func NewCertificate(timeStamp time.Time, fileHash string, signature string) *Certificate {
	return &Certificate{
		TimeStamp: timeStamp,
		FileHash:  fileHash,
		Signature: signature,
	}
}

func (b *Block) mineHash(diff int) (string, int64) {
	hash := ""
	var pow int64 = 0
	diff_substr := ""
	for i := 0; i < diff; i++ {
		diff_substr += "0"
	}
	for {
		hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Data.FileHash+b.Data.Signature+fmt.Sprint(b.Data.TimeStamp)+b.PrevHash+fmt.Sprint(pow))))
		if strings.HasPrefix(hash, diff_substr) {
			break
		}
		pow++
	}
	return hash, pow
}

func NewBlock(data Certificate, prevHash string) *Block {
	new_block := Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Proof:    0,
	}
	new_block.Hash, new_block.Proof = new_block.mineHash(4)
	return &new_block
}

func (bc *BlockChain) GetLatest() *Block {
	return &bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) AddBlock(data string) {
	new_cert := NewCertificate(time.Now(), fmt.Sprintf("%x", sha256.Sum256([]byte(data))), "My Sign")
	prevHash := ""
	if len(bc.chain) != 0 {
		prevHash = bc.GetLatest().Hash
	}
	new_block := NewBlock(*new_cert, prevHash)
	bc.chain = append(bc.chain, *new_block)
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		chain: make([]Block, 0),
	}
}
