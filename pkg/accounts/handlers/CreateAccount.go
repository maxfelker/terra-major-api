package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	models "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrains "github.com/mw-felker/terra-major-api/pkg/terrains"
	"github.com/mw-felker/terra-major-api/pkg/utils"
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

		creatAccountResult := app.DB.Create(&newAccount)

		if creatAccountResult.Error != nil {
			if strings.Contains(creatAccountResult.Error.Error(), "23505") {
				utils.ReturnError(writer, "An account with this email already exists", http.StatusConflict)
			} else {
				utils.ReturnError(writer, creatAccountResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		var newSandbox sandboxModels.Sandbox
		newSandbox.AccountId = newAccount.ID

		sandboxCreateResult := app.DB.Create(&newSandbox)
		if sandboxCreateResult.Error != nil {
			if strings.Contains(sandboxCreateResult.Error.Error(), "23505") {
				utils.ReturnError(writer, "A sandbox for this account already exists", http.StatusConflict)
			} else {
				utils.ReturnError(writer, sandboxCreateResult.Error.Error(), http.StatusInternalServerError)
			}
		}

		chunks := terrains.GenerateChunksForSandbox(newSandbox.ID)

		chunkCreateResult := app.DB.Create(&chunks)
		if chunkCreateResult.Error != nil {
			utils.ReturnError(writer, chunkCreateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		var accountResponse models.AccountResponse
		accountResponse.BaseAccount = newAccount.BaseAccount
		accountResponse.Sandbox = newSandbox

		response, e := json.Marshal(accountResponse)
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
