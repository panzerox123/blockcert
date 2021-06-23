package keygen

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return privateKey, &privateKey.PublicKey
}

func SaveHexKey(filename string, private *rsa.PrivateKey, public *rsa.PublicKey) {
	out_file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out_file.Close()
	private_bytes := x509.MarshalPKCS1PrivateKey(private)
	public_bytes := x509.MarshalPKCS1PublicKey(public)
	private_hex := hex.EncodeToString(private_bytes)
	public_hex := hex.EncodeToString(public_bytes)
	fmt.Printf("Private: %s\nPublic: %s\n", private_hex, public_hex)
	fmt.Fprintf(out_file, "Private: %s\nPublic: %s\n", private_hex, public_hex)
}

func ParsePrivateRSA(input string) *rsa.PrivateKey {
	private_hex, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	ret, err := x509.ParsePKCS1PrivateKey(private_hex)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParsePublicRSA(input string) *rsa.PublicKey {
	public_hex, err := hex.DecodeString(input)
	if err != nil {
		panic(err)
	}
	ret, err := x509.ParsePKCS1PublicKey(public_hex)
	if err != nil {
		panic(err)
	}
	return ret
}

func SignData(data string, private *rsa.PrivateKey) string {
	hashed := sha256.Sum256([]byte(data))
	ret, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, hashed[:])
	if err != nil {
		log.Fatal("Your private key maybe invalid! Try again, or generate a new keypair!", err)
	}
	return hex.EncodeToString(ret)
}

func VerifyData(data string, signature string, public *rsa.PublicKey) bool {
	decoded_signature, err := hex.DecodeString(signature)
	if err != nil {
		log.Fatal("Your public key may be invalid! Try again!", err)
	}
	hashed := sha256.Sum256([]byte(data))
	err = rsa.VerifyPKCS1v15(public, crypto.SHA256, hashed[:], decoded_signature)
	if err == nil {
		return true
	} else {
		return false
	}
}
