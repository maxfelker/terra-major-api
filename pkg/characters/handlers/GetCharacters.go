package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/mw-felker/terra-major-api/pkg/characters/models"
	core "github.com/mw-felker/terra-major-api/pkg/core"
)

func GetCharacters(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		var characters []models.Character
		result := app.DB.Find(&characters)

		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(characters)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
