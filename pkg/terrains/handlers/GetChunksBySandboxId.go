package handlers

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
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

		containerClient, err := app.NoSQL.NewContainer("chunks")
		if err != nil {
			utils.ReturnError(writer, "Failed to create a container client: "+err.Error(), http.StatusInternalServerError)
			return
		}

		pk := azcosmos.NewPartitionKeyString(claims.SandboxId)
		opt := azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@sandboxId", Value: claims.SandboxId},
			},
		}
		queryPager := containerClient.NewQueryItemsPager("SELECT * FROM c WHERE c.sandboxId = @sandboxId", pk, &opt)

		var chunks []models.TerrainChunk
		for queryPager.More() {
			queryResponse, err := queryPager.NextPage(context.Background())
			if err != nil {
				utils.ReturnError(writer, "Failed to get next page of query results: "+err.Error(), http.StatusInternalServerError)
				return
			}

			for _, item := range queryResponse.Items {
				var chunk models.TerrainChunk
				if err := json.Unmarshal(item, &chunk); err != nil {
					utils.ReturnError(writer, "Failed to decode item: "+err.Error(), http.StatusInternalServerError)
					return
				}
				chunks = append(chunks, chunk)
			}
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
