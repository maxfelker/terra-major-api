package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	webAppClient "github.com/mw-felker/terra-major-api/pkg/client/webapp"
	"github.com/mw-felker/terra-major-api/pkg/core"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type PasswordUpdate struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdatePassword(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := webAppClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var passwordUpdate PasswordUpdate
		decoder := json.NewDecoder(request.Body)
		err = decoder.Decode(&passwordUpdate)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		var existingAccount models.Account
		app.DB.First(&existingAccount, "id = ?", claims.AccountId)

		mismatched := bcrypt.CompareHashAndPassword([]byte(existingAccount.Password), []byte(passwordUpdate.CurrentPassword))
		if mismatched != nil {
			utils.ReturnError(writer, "currentPassword is incorrect", http.StatusUnauthorized)
			return
		}

		if !validatePasswordRequirements(passwordUpdate.NewPassword) {
			utils.ReturnError(writer, "newPassword must be at least 8 characters long, contain at least one number, one uppercase letter, and one special character")
			return
		}

		existingAccount.Password = models.GeneratePassword(passwordUpdate.NewPassword)
		app.DB.Save(&existingAccount)

		token := webAppClient.GenerateToken(existingAccount.ID)

		response, e := json.Marshal(webAppClient.TokenResponse{Token: token})
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
