package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	authClient "github.com/mw-felker/terra-major-api/pkg/auth/client"
	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	"github.com/mw-felker/terra-major-api/pkg/terrains"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func CreateSandbox(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		claims, claimsError := authClient.ParseAndValidateToken(request)
		if claimsError != nil {
			utils.ReturnError(writer, claimsError.Error(), http.StatusUnauthorized)
			return
		}

		accountId := claims.AccountId

		decoder := json.NewDecoder(request.Body)
		var newSandbox models.Sandbox
		decoderError := decoder.Decode(&newSandbox)
		if decoderError != nil {
			http.Error(writer, decoderError.Error(), http.StatusBadRequest)
			return
		}

		newSandbox.AccountId = accountId

		result := app.DB.Create(&newSandbox)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "23505") {
				utils.ReturnError(writer, "A sandbox for this characterId already exists", http.StatusConflict)
			} else {
				utils.ReturnError(writer, result.Error.Error(), http.StatusInternalServerError)
			}
		}

		chunks := terrains.GenerateChunksForSandbox(newSandbox.ID)

		chunkCreateResult := app.DB.Create(&chunks)
		if chunkCreateResult.Error != nil {
			utils.ReturnError(writer, chunkCreateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(newSandbox)
		if e != nil {
			utils.ReturnError(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
