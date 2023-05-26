package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	"github.com/mw-felker/centerpoint-instance-api/pkg/core"
	"gorm.io/gorm"
)

func ArchiveCharacter(app *core.App) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		characterId := vars["id"]

		var character models.Character
		result := app.DB.First(&character, "id = ?", characterId)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Character not found", http.StatusNotFound)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		result = app.DB.Delete(&character)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(character)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
