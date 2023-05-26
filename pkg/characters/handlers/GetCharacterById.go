package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	"github.com/mw-felker/centerpoint-instance-api/pkg/core"
	"gorm.io/gorm"
)

func GetCharacterById(app *core.App) http.HandlerFunc {
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

		response, e := json.Marshal(character)
		if e != nil {
			log.Panic(e)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
