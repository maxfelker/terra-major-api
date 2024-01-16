package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"github.com/mw-felker/terra-major-api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func Login(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var suppliedAccount models.Account
		err := decoder.Decode(&suppliedAccount)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		if suppliedAccount.Email == "" {
			utils.ReturnError(writer, "Email is required")
			return
		}

		_, err = mail.ParseAddress(suppliedAccount.Email)
		if err != nil {
			utils.ReturnError(writer, "Invalid email format")
			return
		}

		if suppliedAccount.Password == "" {
			utils.ReturnError(writer, "Password is required")
			return
		}

		var accountInDB models.Account
		if result := app.DB.Where("email = ?", strings.TrimSpace(suppliedAccount.Email)).First(&accountInDB); result.Error != nil {
			utils.ReturnError(writer, "No account with this email", http.StatusNotFound)
			return
		}

		var sandboxInDB sandboxModels.Sandbox
		if querySandboxResult := app.DB.Where("account_id = ?", accountInDB.ID).First(&sandboxInDB); querySandboxResult.Error != nil {
			utils.ReturnError(writer, "No sandbox for this account", http.StatusNotFound)
			return
		}

		userPass := []byte(strings.TrimSpace(suppliedAccount.Password))
		passInDb := []byte(accountInDB.Password)
		mismatched := bcrypt.CompareHashAndPassword(passInDb, userPass)
		if mismatched != nil {
			utils.ReturnError(writer, "Incorrect password", http.StatusUnauthorized)
			return
		}

		token := authClient.GenerateToken(accountInDB.ID, sandboxInDB.ID, "")

		response, e := json.Marshal(authClient.TokenResponse{Token: token})
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
