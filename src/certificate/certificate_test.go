package certificate

import (
	"fmt"
	"testing"
)

func TestAddNode(t *testing.T) {
	blockchain := NewBlockChain()
	blockchain.AddBlock("Hello 1")
	latest := blockchain.GetLatest()
	fmt.Printf("Data: %v\nHash: %s\nPreviousHash: %s\n", latest.Data, latest.Hash, latest.PrevHash)
	blockchain.AddBlock("Hello 2")
	latest = blockchain.GetLatest()
	fmt.Printf("Data: %v\nHash: %s\nPreviousHash: %s\n", latest.Data, latest.Hash, latest.PrevHash)
}
