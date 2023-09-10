package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	"github.com/mw-felker/terra-major-api/pkg/core"
	"github.com/mw-felker/terra-major-api/pkg/terrains"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetChunksBySandboxId(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := authClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		seed := int64(3)
		chunkNeighborhood := terrains.CreateChunkNeighborhood(seed)
		chunks := terrains.FlattenChunksArray(chunkNeighborhood)

		for _, chunk := range chunks {
			chunk.SandboxId = claims.SandboxId
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
