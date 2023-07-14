package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/mw-felker/terra-major-api/pkg/client/public"
	webAppClient "github.com/mw-felker/terra-major-api/pkg/client/webapp"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
)

func GetMySandboxes(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		publicKey, err := public.LoadPublicKey("./keys/terra-major-webapp-public.pem")
		if err != nil {
			http.Error(writer, "Error loading public key: "+err.Error(), http.StatusInternalServerError)
			return
		}

		claims := &webAppClient.WebappClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})

		if err != nil {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		var sandboxes []models.Sandbox
		result := app.DB.Where("account_id = ?", claims.AccountId).Find(&sandboxes)

		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(sandboxes)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
