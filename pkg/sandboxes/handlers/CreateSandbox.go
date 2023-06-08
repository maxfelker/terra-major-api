package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
)

func CreateSandbox(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newSandbox models.Sandbox
		err := decoder.Decode(&newSandbox)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if newSandbox.CharacterId == "" {
			http.Error(writer, "characterId is required", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&newSandbox)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(newSandbox)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
