package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
	"gorm.io/gorm"
)

func UpdateAccount(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		accountId := vars["id"]

		var updatedAccount models.Account
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&updatedAccount)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		var existingAccount models.Account
		findResult := app.DB.First(&existingAccount, "id = ?", accountId)
		if findResult.Error != nil {
			if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
				utils.ReturnError(writer, "Account not found", http.StatusNotFound)
			} else {
				utils.ReturnError(writer, findResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		updatedFields := make(map[string]interface{})

		if updatedAccount.Email != "" {
			_, err := mail.ParseAddress(updatedAccount.Email)
			if err != nil {
				utils.ReturnError(writer, "Invalid email format")
				return
			}
			updatedFields["Email"] = updatedAccount.Email
		}

		if updatedAccount.Password != "" {
			if !validatePasswordRequirements(updatedAccount.Password) {
				utils.ReturnError(writer, "Password must be at least 8 characters long, contain at least one number, one uppercase letter, and one special character")
				return
			}
			updatedFields["Password"] = updatedAccount.Password
		}

		result := app.DB.Model(&existingAccount).Updates(updatedFields)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "23505") {
				utils.ReturnError(writer, "An account with this email already exists", http.StatusConflict)
			} else {
				utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		response, e := json.Marshal(models.AccountResponse{BaseAccount: existingAccount.BaseAccount})
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
