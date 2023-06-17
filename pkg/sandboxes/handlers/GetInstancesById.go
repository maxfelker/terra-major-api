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

func GetInstanceById(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		//sandboxId := vars["sandboxId"]
		instanceId := vars["instanceId"]
		var instance models.Instance
		result := app.DB.First(&instance, "id = ?", instanceId)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Instance not found", http.StatusNotFound)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		response, e := json.Marshal(instance)
		if e != nil {
			log.Panic(e)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
