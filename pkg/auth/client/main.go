// pkg/client/webapp/main.go

package authClient

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type Claims struct {
	AccountId   string `json:"accountId"`
	CharacterId string `json:"characterId"`
	SandboxId   string `json:"sandboxId"`
	jwt.StandardClaims
}

func GenerateToken(accountId string, sandboxId string, characterId string) string {
	privateKey := loadPrivateKey("./keys/terra-major-client-private.pem")
	claims := generateClaim(accountId, sandboxId, characterId)
	return newJwt(claims, privateKey)
}

func ParseAndValidateToken(request *http.Request) (*Claims, error) {
	authHeader := request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	publicKey, err := loadPublicKey("./keys/terra-major-client-public.pem")
	if err != nil {
		return nil, errors.New("error loading public key: " + err.Error())
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func generateClaim(accountId string, sandboxId string, characterId string) *Claims {
	return &Claims{
		AccountId:   accountId,
		SandboxId:   sandboxId,
		CharacterId: characterId,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "terra-major-api",
			Audience:  "terra-major-client",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
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

func loadPublicKey(path string) (*rsa.PublicKey, error) {
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
