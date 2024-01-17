package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	models "github.com/maxfelker/terra-major-api/pkg/accounts/models"
	authClient "github.com/maxfelker/terra-major-api/pkg/auth/client"
	"github.com/maxfelker/terra-major-api/pkg/core"
	sandboxModels "github.com/maxfelker/terra-major-api/pkg/sandboxes/models"
	utils "github.com/maxfelker/terra-major-api/pkg/utils"
	"gorm.io/gorm"
)

func GetMyAccount(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		accountId := claims.AccountId

		var account models.Account
		result := app.DB.First(&account, "id = ?", accountId)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Account not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		var sandbox sandboxModels.Sandbox
		sandboxResult := app.DB.First(&sandbox, "account_id = ?", accountId)

		if sandboxResult.Error != nil {
			if errors.Is(sandboxResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Sandbox not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, sandboxResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		var accountResponse models.AccountResponse
		accountResponse.BaseAccount = account.BaseAccount
		accountResponse.Sandbox = sandbox

		response, e := json.Marshal(accountResponse)
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
