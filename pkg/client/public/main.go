package public

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyFile := path
	publicKeyBytes, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		fmt.Println("Error reading public key: ", err)
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		fmt.Println("Error parsing public key: ", err)
		return nil, err
	}

	return publicKey, nil
}
