package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	accounts "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	characters "github.com/mw-felker/terra-major-api/pkg/characters/models"
	unityClient "github.com/mw-felker/terra-major-api/pkg/client/unity"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxes "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"github.com/mw-felker/terra-major-api/pkg/utils"
	"gorm.io/gorm"
)

type TokenPayload struct {
	AccountId   string `json:"accountId"`
	CharacterId string `json:"characterId"`
	SandboxId   string `json:"sandboxId"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func CreateUnityClientToken(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var tokenPayload TokenPayload
		err := decoder.Decode(&tokenPayload)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		if tokenPayload.AccountId == "" {
			utils.ReturnError(writer, "accountId is required")
			return
		}

		var account accounts.Account
		accountResult := app.DB.First(&account, "id = ?", tokenPayload.AccountId)

		if accountResult.Error != nil {
			if errors.Is(accountResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Account not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, accountResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if tokenPayload.CharacterId == "" {
			utils.ReturnError(writer, "characterId is required")
			return
		}

		var character characters.Character
		characterResult := app.DB.First(&character, "id = ?", tokenPayload.CharacterId)

		if characterResult.Error != nil {
			if errors.Is(characterResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Character not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, characterResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if tokenPayload.SandboxId == "" {
			utils.ReturnError(writer, "sandboxId is required")
			return
		}

		var sandbox sandboxes.Sandbox
		sandboxResult := app.DB.First(&sandbox, "id = ?", tokenPayload.SandboxId)

		if sandboxResult.Error != nil {
			if errors.Is(sandboxResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Sandbox not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, sandboxResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		token := unityClient.GenerateClientToken(tokenPayload.AccountId, tokenPayload.SandboxId, tokenPayload.CharacterId)
		response, err := json.Marshal(TokenResponse{
			Token: token,
		})

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("Error marshalling response to JSON"))
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
