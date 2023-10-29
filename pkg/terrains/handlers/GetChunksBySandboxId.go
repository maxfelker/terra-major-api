package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	"github.com/mw-felker/terra-major-api/pkg/core"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"

	"gorm.io/gorm"
)

func GetChunksBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var chunks []terrainModels.TerrainChunk
		result := app.DB.Find(&chunks, "sandbox_id = ?", claims.SandboxId)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				http.Error(writer, "Chunks not found with that sandboxId", http.StatusNotFound)
			} else {
				http.Error(writer, result.Error.Error(), http.StatusInternalServerError)
			}
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
