package handlers

import (
	"encoding/json"
	"net/http"

	core "github.com/maxfelker/terra-major-api/pkg/core"
	models "github.com/maxfelker/terra-major-api/pkg/sandboxes/models"
)

func GetSandboxes(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var sandboxes []models.Sandbox
		result := app.DB.Find(&sandboxes)

		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(sandboxes)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
