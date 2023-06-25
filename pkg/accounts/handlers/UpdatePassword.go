package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordUpdate struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdatePassword(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		accountId := vars["id"]

		var passwordUpdate PasswordUpdate
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&passwordUpdate)
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

		if passwordUpdate.CurrentPassword == "" {
			utils.ReturnError(writer, "currentPassword is required")
			return
		}

		mismatched := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(passwordUpdate.CurrentPassword))
		if mismatched != nil {
			utils.ReturnError(writer, "currentPassword is incorrect", http.StatusUnauthorized)
			return
		}

		if passwordUpdate.NewPassword == "" {
			utils.ReturnError(writer, "newPassword is required")
			return
		}

		if !validatePasswordRequirements(passwordUpdate.NewPassword) {
			utils.ReturnError(writer, "newPassword must be at least 8 characters long, contain at least one number, one uppercase letter, and one special character")
			return
		}

		existingAccount.Password = models.GeneratePassword(passwordUpdate.NewPassword)
		result := app.DB.Save(&existingAccount)

		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
	}
}
