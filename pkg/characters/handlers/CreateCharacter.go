package handlers

import (
	"encoding/json"
	"net/http"

	accounts "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	characters "github.com/mw-felker/terra-major-api/pkg/characters/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxes "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func CreateCharacter(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newCharacter characters.Character
		err := decoder.Decode(&newCharacter)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		if newCharacter.Name == "" {
			utils.ReturnError(writer, "Name is required")
			return
		}

		if newCharacter.AccountId == "" {
			utils.ReturnError(writer, "AccountId is required")
			return
		}

		// check if Account exists
		var account accounts.Account
		if err := app.DB.Where("id = ?", newCharacter.AccountId).First(&account).Error; err != nil {
			utils.ReturnError(writer, "Account not found")
			return
		}

		result := app.DB.Create(&newCharacter)
		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error())
			return
		}

		var newSandbox sandboxes.Sandbox
		newSandbox.CharacterId = newCharacter.ID
		newSandbox.AccountId = newCharacter.AccountId // MOVE TO JWT
		sandboxResult := app.DB.Create(&newSandbox)
		if sandboxResult.Error != nil {
			utils.ReturnError(writer, sandboxResult.Error.Error())
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
