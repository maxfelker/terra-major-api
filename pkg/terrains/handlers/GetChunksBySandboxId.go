package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetChunksBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		var chunks []models.TerrainChunk
		result := app.DB.Where("sandbox_id = ?", claims.SandboxId).Find(&chunks)

		if result.Error != nil {
			utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		if err := json.NewEncoder(gz).Encode(chunks); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := gz.Close(); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Content-Encoding", "gzip")
		writer.Write(buf.Bytes())
	}
}
