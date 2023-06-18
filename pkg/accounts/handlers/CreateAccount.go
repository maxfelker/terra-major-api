package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
)

func validatePasswordRequirements(password string) bool {
	if len(password) < 8 || !regexp.MustCompile(`[0-9]`).MatchString(password) || !regexp.MustCompile(`[A-Z]`).MatchString(password) || !regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password) {
		return false
	}
	return true
}

func CreateAccount(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newAccount models.Account
		err := decoder.Decode(&newAccount)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if newAccount.Email == "" {
			http.Error(writer, "Email is required", http.StatusBadRequest)
			return
		}

		_, err = mail.ParseAddress(newAccount.Email)
		if err != nil {
			http.Error(writer, "Invalid email format", http.StatusBadRequest)
			return
		}

		if newAccount.Password == "" {
			http.Error(writer, "Password is required", http.StatusBadRequest)
			return
		}

		if !validatePasswordRequirements(newAccount.Password) {
			http.Error(writer, "Password must be at least 8 characters long, contain at least one number, one uppercase letter, and one special character", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&newAccount)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "23505") {
				http.Error(writer, "An account with this email already exists", http.StatusConflict)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		response, e := json.Marshal(models.AccountResponse{BaseAccount: newAccount.BaseAccount})
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
