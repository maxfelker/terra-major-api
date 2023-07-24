package handlers

import (
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

		fmt.Println(claims)

		seed := int64(42)
		world := terrains.NewWorld(4, 128, 32, seed)

		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(writer)
		defer gz.Close()

		if err := json.NewEncoder(gz).Encode(world); err != nil {
			utils.ReturnError(writer, err.Error(), http.StatusInternalServerError)
		}
	}
}
