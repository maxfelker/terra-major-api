package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path.")
		return
	}

	filePath := os.Args[1]

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	privateKeyPem := pem.EncodeToMemory(privateKeyBlock)

	err = ioutil.WriteFile(filePath, privateKeyPem, 0600)
	if err != nil {
		fmt.Println("Error writing private key to file:", err)
		return
	}

	fmt.Println("Successfully generated", filePath)
}
