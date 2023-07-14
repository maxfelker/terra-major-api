package clientUtils

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

type JWTClaim interface {
	jwt.Claims
}

func GenerateToken(claims JWTClaim, privateKeyFile string) string {
	privateKey := loadPrivateKey(privateKeyFile)
	return newJwt(claims, privateKey)
}

func newJwt(claims jwt.Claims, privateKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Error signing token: ", err)
	}
	return tokenString
}

func loadPrivateKey(privateKeyFile string) *rsa.PrivateKey {
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		fmt.Println("Error reading private key: ", err)
		return nil
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		fmt.Println("Error parsing private key: ", err)
		return nil
	}

	return privateKey
}
