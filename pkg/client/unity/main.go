package unityClient

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	AccountId   string `json:"accountId"`
	CharacterId string `json:"characterId"`
	SandboxId   string `json:"sandboxId"`
	jwt.StandardClaims
}

func GenerateClientToken(accountId string, sandboxId string, characterId string) string {
	privateKey := loadPrivateKey()
	claims := generateClaim(accountId, sandboxId, characterId)
	return newJwt(claims, privateKey)
}

func newJwt(claims *Claims, privateKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println("Error signing token: ", err)
	}
	return tokenString
}

func generateClaim(accountId string, sandboxId string, characterId string) *Claims {
	return &Claims{
		AccountId:   accountId,
		SandboxId:   sandboxId,
		CharacterId: characterId,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "TerraMajorAPI",
			Audience:  "UnityClient",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
}

func loadPrivateKey() *rsa.PrivateKey {
	privateKeyFile := "./keys/unity-client.pem"
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
