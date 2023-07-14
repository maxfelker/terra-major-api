package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a file path and a key name.")
		return
	}

	filePath := os.Args[1]
	keyName := os.Args[2]

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return
	}

	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	privateKeyPem := pem.EncodeToMemory(privateKeyBlock)
	publicKeyPem := pem.EncodeToMemory(publicKeyBlock)

	privateKeyPath := filepath.Join(filePath, keyName+"-private.pem")
	err = ioutil.WriteFile(privateKeyPath, privateKeyPem, 0600)
	if err != nil {
		fmt.Println("Error writing private key to file:", err)
		return
	}

	publicKeyPath := filepath.Join(filePath, keyName+"-public.pem")
	err = ioutil.WriteFile(publicKeyPath, publicKeyPem, 0600)
	if err != nil {
		fmt.Println("Error writing public key to file:", err)
		return
	}

	fmt.Println("Successfully generated", privateKeyPath, "and", publicKeyPath)
}
