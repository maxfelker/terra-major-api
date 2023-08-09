package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	characters "github.com/mw-felker/terra-major-api/pkg/characters/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxes "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"github.com/mw-felker/terra-major-api/pkg/utils"
	"gorm.io/gorm"
)

type TokenPayload struct {
	CharacterId string `json:"characterId"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func CreateUnityClientToken(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(request.Body)
		var tokenPayload TokenPayload
		err = decoder.Decode(&tokenPayload)
		if err != nil {
			utils.ReturnError(writer, err.Error())
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

		var sandbox sandboxes.Sandbox
		sandboxResult := app.DB.First(&sandbox, "character_id = ?", tokenPayload.CharacterId)

		if sandboxResult.Error != nil {
			if errors.Is(sandboxResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Sandbox not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, sandboxResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		token := authClient.GenerateToken(claims.AccountId, sandbox.ID, tokenPayload.CharacterId)
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
