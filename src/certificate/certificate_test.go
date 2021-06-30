package certificate

import (
	"fmt"
	"testing"

	"github.com/panzerox123/blockcert/src/keygen"
)

func TestAddNode(t *testing.T) {
	private_key, _ := keygen.GenerateKeyPair(512)
	blockchain := NewBlockChain()
	blockchain.AddBlock([]byte("Hello 1"), private_key, 5)
	latest := blockchain.GetLatest()
	fmt.Printf("Data: %v\nHash: %s\nPreviousHash: %s\n", latest.Data, latest.Hash, latest.PrevHash)
	blockchain.AddBlock([]byte("Hello 2"), private_key, 5)
	latest = blockchain.GetLatest()
	fmt.Printf("Data: %v\nHash: %s\nPreviousHash: %s\n", latest.Data, latest.Hash, latest.PrevHash)
}
