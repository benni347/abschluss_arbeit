package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func generateECCKeyPair(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func main() {
	curve := elliptic.P256()
	privateKey, publicKey, err := generateECCKeyPair(curve)
	if err != nil {
		fmt.Printf("Error generating key pair: %v\n", err)
		return
	}

	fmt.Printf("Private Key: %x\n", privateKey.D.Bytes())
	fmt.Printf("Public Key (X): %x\n", publicKey.X.Bytes())
	fmt.Printf("Public Key (Y): %x\n", publicKey.Y.Bytes())
}
