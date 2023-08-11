package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mw-felker/terra-major-api/pkg/core"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
)

func CreateTerrainChunk(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var newTerrainChunk terrainModels.TerrainChunk
		err := decoder.Decode(&newTerrainChunk)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		result := app.DB.Create(&newTerrainChunk)
		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		/*response, e := json.Marshal(newTerrainChunk)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}*/

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
	}
}
