package handlers

import (
	"encoding/json"
	"net/http"

	accounts "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	characters "github.com/mw-felker/terra-major-api/pkg/characters/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func CreateCharacter(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(request.Body)
		var newCharacter characters.Character
		err = decoder.Decode(&newCharacter)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		if newCharacter.Name == "" {
			utils.ReturnError(writer, "Name is required")
			return
		}

		var account accounts.Account
		if err := app.DB.Where("id = ?", claims.AccountId).First(&account).Error; err != nil {
			utils.ReturnError(writer, "Account not found")
			return
		}

		newCharacter.AccountId = claims.AccountId
		newCharacter.SandboxId = claims.SandboxId

		result := app.DB.Create(&newCharacter)
		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error())
			return
		}

		response, e := json.Marshal(newCharacter)
		if e != nil {
			utils.ReturnError(writer, e.Error())
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
