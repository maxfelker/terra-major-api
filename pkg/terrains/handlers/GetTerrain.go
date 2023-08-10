package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	terrains "github.com/mw-felker/terra-major-api/pkg/terrains"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetTerrain(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}
		chunkCount := 6
		chunkDimension := 128
		chunkHeight := 32
		seed := int64(7)
		world := terrains.NewWorld(chunkCount, chunkDimension, chunkHeight, seed)

		// first, encode to a buffer
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)

		if err := json.NewEncoder(gz).Encode(world); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := gz.Close(); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// only write to the ResponseWriter if there are no errors
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Content-Encoding", "gzip")
		if _, err := writer.Write(buf.Bytes()); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}
