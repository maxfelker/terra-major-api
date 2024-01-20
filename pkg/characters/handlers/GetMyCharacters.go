// pkg/characters/handlers/characters.go

package handlers

import (
	"encoding/json"
	"net/http"

	authClient "github.com/maxfelker/terra-major-api/pkg/auth/client"
	models "github.com/maxfelker/terra-major-api/pkg/characters/models"
	core "github.com/maxfelker/terra-major-api/pkg/core"
	utils "github.com/maxfelker/terra-major-api/pkg/utils"
)

func GetMyCharacters(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var characters []models.Character
		result := app.DB.Where("account_id = ?", claims.AccountId).Find(&characters)

		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(characters)

		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
