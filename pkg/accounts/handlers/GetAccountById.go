package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/maxfelker/terra-major-api/pkg/accounts/models"
	"github.com/maxfelker/terra-major-api/pkg/core"
	utils "github.com/maxfelker/terra-major-api/pkg/utils"
	"gorm.io/gorm"
)

func GetAccountById(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		accountId := vars["id"]
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

		response, e := json.Marshal(models.AccountResponse{BaseAccount: account.BaseAccount})
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
