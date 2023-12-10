package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxesModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetChunksBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		vars := mux.Vars(request)
		sandboxId := vars["sandboxId"]

		query := request.URL.Query()

		xQuery, _ := strconv.ParseFloat(query.Get("x"), 64)
		x := float32(xQuery)

		zQuery, _ := strconv.ParseFloat(query.Get("z"), 64)
		z := float32(zQuery)

		y := float32(0) // all chunks have a 0 y value

		var queryPosition = sandboxesModels.Vector3{
			X: &x,
			Y: &y,
			Z: &z,
		}

		if queryPosition.X == nil || queryPosition.Z == nil {
			http.Error(writer, "x and z coordinates are required", http.StatusBadRequest)
			return
		}

		radius, err := strconv.ParseFloat(query.Get("radius"), 32)
		if err != nil {
			http.Error(writer, "Invalid radius", http.StatusBadRequest)
			return
		}

		var chunks []terrainModels.TerrainChunk
		result := app.DB.Raw(`
			SELECT * FROM terrain_chunks
			WHERE 
			SQRT(
				POWER(CAST(position->>'x' AS FLOAT) - ?, 2) +
				POWER(CAST(position->>'z' AS FLOAT) - ?, 2)
			) <= ? AND sandbox_id = ?
		`, *queryPosition.X, *queryPosition.Z, radius, sandboxId).Scan(&chunks)

		if result.Error != nil {
			http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(chunks)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
