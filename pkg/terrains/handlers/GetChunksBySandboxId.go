package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mw-felker/terra-major-api/pkg/core"
	sandboxesModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetChunksBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()

		xQuery, _ := strconv.ParseFloat(query.Get("x"), 64)
		x := float32(xQuery)

		yQuery, _ := strconv.ParseFloat(query.Get("y"), 64)
		y := float32(yQuery)

		zQuery, _ := strconv.ParseFloat(query.Get("z"), 64)
		z := float32(zQuery)

		var queryPosition = sandboxesModels.Vector3{
			X: &x,
			Y: &y,
			Z: &z,
		}

		if queryPosition.X == nil || queryPosition.Y == nil || queryPosition.Z == nil {
			http.Error(writer, "All position fields (x, y, z) are required", http.StatusBadRequest)
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
			POWER(CAST(position->>'x' AS FLOAT) - ?, 2) +
			POWER(CAST(position->>'y' AS FLOAT) - ?, 2) +
			POWER(CAST(position->>'z' AS FLOAT) - ?, 2) 
			<= POWER(?, 2)
		`, queryPosition.X, queryPosition.Y, queryPosition.Z, radius).Scan(&chunks)

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
