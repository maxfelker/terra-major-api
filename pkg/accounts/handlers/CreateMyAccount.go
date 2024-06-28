package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	models "github.com/maxfelker/terra-major-api/pkg/accounts/models"
	authClient "github.com/maxfelker/terra-major-api/pkg/auth/client"
	"github.com/maxfelker/terra-major-api/pkg/core"
	"github.com/maxfelker/terra-major-api/pkg/utils"
)

func CreateMyAccount(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newAccount models.Account
		err := decoder.Decode(&newAccount)
		if err != nil {
			utils.ReturnError(writer, err.Error())
			return
		}

		if newAccount.Email == "" {
			utils.ReturnError(writer, "Email is required")
			return
		}

		_, err = mail.ParseAddress(newAccount.Email)
		if err != nil {
			utils.ReturnError(writer, "Invalid email format")
			return
		}

		if newAccount.Password == "" {
			utils.ReturnError(writer, "Password is required")
			return
		}

		if !validatePasswordRequirements(newAccount.Password) {
			utils.ReturnError(writer, "Password must be at least 8 characters long, contain at least one number, one uppercase letter, and one special character")
			return
		}

		result := app.DB.Create(&newAccount)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "23505") {
				utils.ReturnError(writer, "An account with this email already exists", http.StatusConflict)
			} else {
				utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		token := authClient.GenerateToken(newAccount.ID, "", "")

		response, e := json.Marshal(authClient.TokenResponse{Token: token})
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
