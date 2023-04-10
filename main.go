package main

import (
	"context"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	encryption "github.com/benni347/encryption"
)

func main() {
	log.SetFlags(0)

	encrypt()
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func encrypt() {
	Verbose := true

	msg := []byte("Profil. Berufsvorbereitung")

	curve := elliptic.P256()
	privateKeyEcc, publicKeyEcc, err := encryption.GenerateECCKeyPair(curve)
	if err != nil {
		encryption.PrintError("Generating key pair for ecc the error is", err)
		return
	}

	hash := encryption.CalculateHash(msg)
	signature, err := encryption.SignEcc(privateKeyEcc, hash)
	if err != nil {
		encryption.PrintError("During signing the message with ECC the error is", err)
		return
	}

	valid := encryption.VerifyEcc(publicKeyEcc, hash, signature)
	if valid {
		encryption.PrintInfo("Signature is valid!", Verbose)
	} else {
		encryption.PrintInfo("Signature is not valid!", Verbose)
	}

	modeName := "Dilithium5-AES"
	// Generate Dilithium key pair
	publicKey, privateKey, err := encryption.GenerateDilithiumKeyPair(modeName)
	if err != nil {
		encryption.PrintError("During generating Dilithium key pair the error is", err)
		return
	}

	encryption.PrintInfo(
		"CRYSTAl-Dilithium Public Key: "+hex.EncodeToString(publicKey.Bytes()),
		Verbose,
	)

	// Pack and unpack Dilithium keys
	packedPublicKey, packedPrivateKey := encryption.PackDilithiumKeys(publicKey, privateKey)
	publicKey2, privateKey2 := encryption.UnpackDilithiumKeys(
		modeName,
		packedPublicKey,
		packedPrivateKey,
	)

	// Sign and verify with Dilithium keys
	signature, signatureLength, err := encryption.SignDilithium(privateKey2, msg, modeName)
	if err != nil {
		encryption.PrintError("During signing the message with Dilithium the error is", err)
	}

	encryption.PrintInfo(fmt.Sprintf("Signature length is: %d", signatureLength), Verbose)

	encryption.PrintInfo("CRYSTAl-Dilithium Signature: "+hex.EncodeToString(signature), Verbose)

	valid, err = encryption.VerifyDilithium(publicKey2, msg, signature, modeName)
	if err != nil {
		encryption.PrintError("During verifying the message with Dilithium the error is", err)
	}

	if valid {
		encryption.PrintInfo("CRYSTAl-Dilithium Signature is valid!", Verbose)
	} else {
		encryption.PrintInfo("CRYSTAl-Dilithium Signature is not valid!", Verbose)
	}
}

// run initializes the chatServer and then
// starts a http.Server for the passed in address.
func run() error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	l, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		return err
	}
	log.Printf("listening on http://%v", l.Addr())

	server := newServer()
	s := &http.Server{
		Handler:      server,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	errc := make(chan error, 1)
	go func() {
		errc <- s.Serve(l)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return s.Shutdown(ctx)
}
