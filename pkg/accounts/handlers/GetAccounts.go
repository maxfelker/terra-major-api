package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/maxfelker/terra-major-api/pkg/accounts/models"
	"github.com/maxfelker/terra-major-api/pkg/core"
	"github.com/maxfelker/terra-major-api/pkg/utils"
)

func GetAccounts(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email := request.URL.Query().Get("email")

		var accounts []models.Account
		query := app.DB

		if email != "" {
			query = query.Where("email LIKE ?", "%"+strings.TrimSpace(email)+"%")
		}

		result := query.Find(&accounts)
		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		baseAccounts := []models.BaseAccount{}
		for _, account := range accounts {
			baseAccounts = append(baseAccounts, models.BaseAccount{
				ID:      account.ID,
				Email:   account.Email,
				Created: account.Created,
				Updated: account.Updated,
			})
		}

		response, err := json.Marshal(baseAccounts)

		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
