// pkg/sandboxes/handlers/sandboxes.go

package handlers

import (
	"encoding/json"
	"net/http"

	authClient "github.com/maxfelker/terra-major-api/pkg/auth/client"
	core "github.com/maxfelker/terra-major-api/pkg/core"
	models "github.com/maxfelker/terra-major-api/pkg/sandboxes/models"
	utils "github.com/maxfelker/terra-major-api/pkg/utils"
)

func GetMySandboxes(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var sandboxes []models.Sandbox
		result := app.DB.Where("account_id = ?", claims.AccountId).Find(&sandboxes)

		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(sandboxes)

		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
