package certificate

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/panzerox123/blockcert/src/keygen"
)

// Create new object of type Certificate
func NewCertificate(timeStamp time.Time, fileHash string, priv_key *rsa.PrivateKey) *Certificate {
	new_cert := Certificate{
		TimeStamp: timeStamp,
		FileHash:  fileHash,
		Signature: "",
	}
	new_cert.signCertificate(priv_key)
	return &new_cert
}

// Calculate hash values for timestamp and the filehash
func (c *Certificate) calcHash() string {
	hashed := sha256.Sum256([]byte(fmt.Sprint(c.TimeStamp) + c.FileHash))
	return hex.EncodeToString(hashed[:])
}

// Sign certificate with a private RSA key
func (c *Certificate) signCertificate(priv_key *rsa.PrivateKey) {
	c.Signature = keygen.SignData(c.calcHash(), priv_key)
}

// Verify the certificate hash and signature
func (c *Certificate) VerifyCertificate(pub_key *rsa.PublicKey, data string) bool {
	hashed := sha256.Sum256([]byte(data))
	return keygen.VerifyData(c.calcHash(), c.Signature, pub_key) && hex.EncodeToString(hashed[:]) == c.FileHash
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
func (bc *BlockChain) AddBlock(data string, priv_key *rsa.PrivateKey) {
	new_cert := NewCertificate(time.Now(), fmt.Sprintf("%x", sha256.Sum256([]byte(data))), priv_key)
	prevHash := ""
	if len(bc.chain) != 0 {
		prevHash = bc.GetLatest().Hash
	}
	new_block := NewBlock(*new_cert, prevHash)
	bc.chain = append(bc.chain, *new_block)
}

// Check if a given certificate exists on the blockchain
func (bc *BlockChain) CheckSignature(data string, public_key *rsa.PublicKey) bool {
	for _, x := range bc.chain {
		if x.Data.VerifyCertificate(public_key, data) {
			return true
		}
	}
	return false
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
