package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maxfelker/terra-major-api/pkg/core"
	models "github.com/maxfelker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

func ArchiveSandbox(app *core.App) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		sandboxId := vars["id"]

		var sandbox models.Sandbox
		result := app.DB.First(&sandbox, "id = ?", sandboxId)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Sandbox not found", http.StatusNotFound)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		result = app.DB.Delete(&sandbox)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(sandbox)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
