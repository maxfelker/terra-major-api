package handlers

import (
	"encoding/json"
	"net/http"

	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	"github.com/mw-felker/centerpoint-instance-api/pkg/core"
)

func CreateCharacter(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newCharacter models.Character
		err := decoder.Decode(&newCharacter)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if newCharacter.Name == "" {
			http.Error(writer, "Name is required", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&newCharacter)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(newCharacter)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
