package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func main() {
	// Generate a new private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}

	// Encode the private key into PKCS1 ASN.1 DER format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Create a pem.Block with the private key
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// Encode the private key into PEM format
	privateKeyPem := pem.EncodeToMemory(privateKeyBlock)

	// Write the private key PEM to a file
	err = ioutil.WriteFile("unity-client.pem", privateKeyPem, 0600)
	if err != nil {
		fmt.Println("Error writing private key to file:", err)
		return
	}

	fmt.Println("Successfully generated private.pem")
}
