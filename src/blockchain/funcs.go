package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

func (block *DigitalCert) genHash() ([32]byte, int64) {
	hash := [32]byte{0}
	var pof int64 = 0
	for {
		hash = sha256.Sum256([]byte(block.Data + fmt.Sprint(block.TimeStamp) + fmt.Sprint(block.CertificateVersion) + fmt.Sprint(block.PrevHash) + fmt.Sprint(pof)))
		if bytes.Equal(hash[0:1], []byte{0}) {
			break
		}
		pof++
	}
	return hash, pof
}

func (block *DigitalCert) genImmHash() [32]byte {
	return sha256.Sum256([]byte(block.Data + fmt.Sprint(block.TimeStamp) + fmt.Sprint(block.CertificateVersion) + fmt.Sprint(block.PrevHash) + fmt.Sprint(block.ProofWork)))
}

func NewDigitalCert(cert_ver int64, time_stamp int64, data string, prev_hash [32]byte) *DigitalCert {
	new_block := DigitalCert{
		CertificateVersion: cert_ver,
		TimeStamp:          time_stamp,
		Data:               data,
		Hash:               [32]byte{0},
		PrevHash:           prev_hash,
		ProofWork:          0,
	}
	new_block.Hash, new_block.ProofWork = new_block.genHash()
	return &new_block
}

func (bchain *DigitalCertChain) initChain(owner string) {
	genesis_block := NewDigitalCert(-1, time.Now().Unix(), owner, [32]byte{0})
	bchain.Chain = append(bchain.Chain, *genesis_block)
}

func (bchain *DigitalCertChain) getOwner() *DigitalCert {
	return &bchain.Chain[0]
}

func (bchain *DigitalCertChain) GetLatest() *DigitalCert {
	return &bchain.Chain[len(bchain.Chain)-1]
}

func (bchain *DigitalCertChain) CheckValid() bool {
	for i := 0; i < len(bchain.Chain); i++ {
		if bchain.Chain[i].Hash != bchain.Chain[i].genImmHash() {
			return false
		}
		if i > 0 && bchain.Chain[i].PrevHash != bchain.Chain[i-1].Hash {
			return false
		}
	}
	return true
}

func (bchain *DigitalCertChain) AddCert(cert_ver int64, data string) {
	prev_cert := bchain.GetLatest()
	new_cert := NewDigitalCert(cert_ver, time.Now().Unix(), data, prev_cert.Hash)
	bchain.Chain = append(bchain.Chain, *new_cert)
}

func NewDigitalCertChain(owner string) *DigitalCertChain {
	new_chain := DigitalCertChain{
		Chain: make([]DigitalCert, 0),
	}
	new_chain.initChain(owner)
	return &new_chain
}
