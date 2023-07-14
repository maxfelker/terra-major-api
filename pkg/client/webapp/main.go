// pkg/client/webapp/main.go

package webAppClient

import (
	"time"

	"github.com/golang-jwt/jwt"
	clientUtils "github.com/mw-felker/terra-major-api/pkg/client/utils"
)

type WebappClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
}

func GenerateToken(accountId string) string {
	claims := generateWebappClaim(accountId)
	return clientUtils.GenerateToken(claims, "./keys/terra-major-webapp-private.pem")
}

func generateWebappClaim(accountId string) *WebappClaims {
	return &WebappClaims{
		AccountId: accountId,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "terra-major-api",
			Audience:  "terra-major-webapp",
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
}
