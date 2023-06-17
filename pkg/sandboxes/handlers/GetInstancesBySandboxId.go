package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

func GetInstancesBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		sandboxId := vars["sandboxId"]
		var instances []models.Instance
		result := app.DB.Find(&instances, "sandbox_id = ?", sandboxId)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Sandbox not found", http.StatusNotFound)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		response, e := json.Marshal(instances)
		if e != nil {
			log.Panic(e)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
