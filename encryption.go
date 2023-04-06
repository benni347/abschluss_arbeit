package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"math/big"

	"golang.org/x/crypto/blake2b"
)

func generateECCKeyPair(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func calculateHash(message []byte) []byte {
	hash, err := blake2b.New256(nil)
	if err != nil {
		fmt.Printf("Error creating hash: %v\n", err)
		return nil
	}
	hash.Write(message)
	return hash.Sum(nil)
}

func sign(privateKey *ecdsa.PrivateKey, messageHash []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, messageHash)
	if err != nil {
		return nil, err
	}

	curveBits := privateKey.PublicKey.Curve.Params().BitSize
	keyBytes := (curveBits + 7) / 8

	signature := make([]byte, keyBytes*2)
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	copy(signature[keyBytes-len(rBytes):], rBytes)
	copy(signature[keyBytes*2-len(sBytes):], sBytes)

	return signature, nil
}

func verify(publicKey *ecdsa.PublicKey, messageHash []byte, signature []byte) bool {
	curveBits := publicKey.Curve.Params().BitSize
	keyBytes := (curveBits + 7) / 8

	r := new(big.Int).SetBytes(signature[:keyBytes])
	s := new(big.Int).SetBytes(signature[keyBytes:])

	return ecdsa.Verify(publicKey, messageHash, r, s)
}

func main() {
	curve := elliptic.P256()
	privateKey, publicKey, err := generateECCKeyPair(curve)
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}

	msg := []byte("Hello World!")

	hash := calculateHash(msg)
	fmt.Printf("Hash: %x\n", hash)
	signature, err := sign(privateKey, hash)
	if err != nil {
		fmt.Printf("Error signing message: %v\n", err)
		return
	}

	fmt.Printf("Signature: %x\n", signature)

	valid := verify(publicKey, hash, signature)
	if valid {
		fmt.Println("Signature is valid!")
	} else {
		fmt.Println("Signature is not valid!")
	}
}
