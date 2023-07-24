package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"

	webAppClient "github.com/mw-felker/terra-major-api/pkg/client/webapp"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	terrains "github.com/mw-felker/terra-major-api/pkg/terrains"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func GetTerrain(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		claims, err := webAppClient.ParseAndValidateToken(request)
		if err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusUnauthorized)
			return
		}

		fmt.Println("Generating code for " + claims.AccountId)

		seed := int64(42)
		world := terrains.NewWorld(4, 128, 32, seed)

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
