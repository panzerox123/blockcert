package certificate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func FileByteOut(srcFile string) []byte {
	data, err := ioutil.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}
	return data
}

func (bc *BlockChain) SaveBlockchainJson() {
	jsoned_data, err := json.Marshal(bc)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("blockdata.json", jsoned_data, 0644)
	if err != nil {
		panic(err)
	}
}

func ReadBlockChain() *BlockChain {
	reader, err := os.Open("blockdata.json")
	bc := NewBlockChain()
	if err != nil {
		fmt.Println(err)
		return bc
	}
	err = json.NewDecoder(reader).Decode(&bc)
	if err != nil {
		fmt.Println(err)
		return bc
	}
	return bc
}
