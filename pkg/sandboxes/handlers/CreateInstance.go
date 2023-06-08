package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
)

func CreateInstance(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		sandboxId := vars["sandboxId"]

		if sandboxId == "" {
			http.Error(writer, "sandboxId not found", http.StatusNotFound)
			return
		}

		decoder := json.NewDecoder(request.Body)
		var newInstance models.Instance
		newInstance.SandboxId = sandboxId
		err := decoder.Decode(&newInstance)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if newInstance.CharacterId == "" {
			http.Error(writer, "characterId is required", http.StatusBadRequest)
			return
		}

		if newInstance.PrefabName == "" {
			http.Error(writer, "prefabName is required", http.StatusBadRequest)
			return
		}

		if newInstance.Position.X == nil || newInstance.Position.Y == nil || newInstance.Position.Z == nil {
			http.Error(writer, "All position fields (x, y, z) are required", http.StatusBadRequest)
			return
		}

		if newInstance.Rotation.X == nil || newInstance.Rotation.Y == nil || newInstance.Rotation.Z == nil {
			http.Error(writer, "All rotation fields (x, y, z) are required", http.StatusBadRequest)
			return
		}

		if newInstance.Health < 1 {
			http.Error(writer, "Healh must be greater than 0", http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&newInstance)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(newInstance)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
