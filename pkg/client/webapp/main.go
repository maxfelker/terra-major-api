// pkg/client/webapp/main.go

package webAppClient

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	public "github.com/mw-felker/terra-major-api/pkg/client/public"
	clientUtils "github.com/mw-felker/terra-major-api/pkg/client/utils"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type WebappClaims struct {
	AccountId string `json:"accountId"`
	jwt.StandardClaims
}

func GenerateToken(accountId string) string {
	claims := generateWebappClaim(accountId)
	return clientUtils.GenerateToken(claims, "./keys/terra-major-webapp-private.pem")
}

func ParseAndValidateToken(request *http.Request) (*WebappClaims, error) {
	authHeader := request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	publicKey, err := public.LoadPublicKey("./keys/terra-major-webapp-public.pem")
	if err != nil {
		return nil, errors.New("error loading public key: " + err.Error())
	}

	claims := &WebappClaims{}

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
