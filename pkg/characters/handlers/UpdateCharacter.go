package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/terra-major-api/pkg/characters/models"
	"github.com/mw-felker/terra-major-api/pkg/core"
	"gorm.io/gorm"
)

func UpdateCharacter(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		vars := mux.Vars(request)
		characterId := vars["id"]

		var updatedCharacter models.Character
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&updatedCharacter)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var existingCharacter models.Character
		findResult := app.DB.First(&existingCharacter, "id = ?", characterId)
		if findResult.Error != nil {
			if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Character not found", http.StatusNotFound)
			} else {
				http.Error(writer, findResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if updatedCharacter.Name != "" {
			existingCharacter.Name = updatedCharacter.Name
		}

		if updatedCharacter.Bio != "" {
			existingCharacter.Bio = updatedCharacter.Bio
		}

		result := app.DB.Save(&existingCharacter)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(existingCharacter)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
