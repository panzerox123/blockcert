package certificate

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Create new object of type Certificate
func NewCertificate(timeStamp time.Time, fileHash string, signature string) *Certificate {
	return &Certificate{
		TimeStamp: timeStamp,
		FileHash:  fileHash,
		Signature: signature,
	}
}

// Hash Mining function
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

// Calculate Hash with given Proof of Work
func (b *Block) calcHash() string {
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(b.Data.FileHash+b.Data.Signature+fmt.Sprint(b.Data.TimeStamp)+b.PrevHash+fmt.Sprint(b.Proof))))
	return hash
}

// Create new object of type Block
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

// Get the last block in the BlockChain
func (bc *BlockChain) GetLatest() *Block {
	return &bc.chain[len(bc.chain)-1]
}

// Add a block to the BlockChain
func (bc *BlockChain) AddBlock(data string) {
	new_cert := NewCertificate(time.Now(), fmt.Sprintf("%x", sha256.Sum256([]byte(data))), "My Sign")
	prevHash := ""
	if len(bc.chain) != 0 {
		prevHash = bc.GetLatest().Hash
	}
	new_block := NewBlock(*new_cert, prevHash)
	bc.chain = append(bc.chain, *new_block)
}

// TO BE DELETED AFTER TESTS
func (bc *BlockChain) AlterChainTest() {
	bc.chain[2].Hash = "90"
}

// Check if the BlockChain is valid
func (bc *BlockChain) ChainValid() bool {
	for i := 0; i < len(bc.chain); i++ {
		if bc.chain[i].calcHash() != bc.chain[i].Hash {
			return false
		}
		if i > 0 && bc.chain[i].PrevHash != bc.chain[i-1].Hash {
			return false
		}
	}
	return true
}

// Create new object of type BlockChain
func NewBlockChain() *BlockChain {
	return &BlockChain{
		chain: make([]Block, 0),
	}
}
