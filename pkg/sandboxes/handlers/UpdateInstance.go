package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"gorm.io/gorm"
)

func UpdateInstance(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		vars := mux.Vars(request)
		instanceId := vars["instanceId"]
		sandboxId := vars["sandboxId"]

		if sandboxId == "" {
			http.Error(writer, "Sandbox not found", http.StatusNotFound)
			return
		}

		var updatedInstance models.Instance
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&updatedInstance)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		var existingInstance models.Instance
		findResult := app.DB.Where("id = ? AND sandbox_id = ?", instanceId, sandboxId).First(&existingInstance)
		if findResult.Error != nil {
			if errors.Is(findResult.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Instance not found", http.StatusNotFound)
			} else {
				http.Error(writer, findResult.Error.Error(), http.StatusInternalServerError)
			}
			return
		}

		if updatedInstance.Position.X != nil && updatedInstance.Position.Y != nil && updatedInstance.Position.Z != nil {
			existingInstance.Position = updatedInstance.Position
		} else if updatedInstance.Position.X != nil || updatedInstance.Position.Y != nil || updatedInstance.Position.Z != nil {
			http.Error(writer, "All position fields (x, y, z) are required", http.StatusBadRequest)
			return
		}

		if updatedInstance.Rotation.X != nil && updatedInstance.Rotation.Y != nil && updatedInstance.Rotation.Z != nil {
			existingInstance.Rotation = updatedInstance.Rotation
		} else if updatedInstance.Rotation.X != nil || updatedInstance.Rotation.Y != nil || updatedInstance.Rotation.Z != nil {
			http.Error(writer, "All rotation fields (x, y, z) are required", http.StatusBadRequest)
			return
		}

		if updatedInstance.PrefabName != "" {
			existingInstance.PrefabName = updatedInstance.PrefabName
		}

		result := app.DB.Save(&existingInstance)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(existingInstance)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
