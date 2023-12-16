package handlers

import (
	"encoding/json"
	"fmt"
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
		//characterId := claims.CharacterId

		decoder := json.NewDecoder(request.Body)
		var newSandbox models.Sandbox
		decoderError := decoder.Decode(&newSandbox)
		if decoderError != nil {
			http.Error(writer, decoderError.Error(), http.StatusBadRequest)
			return
		}

		if newSandbox.CharacterId == "" {
			http.Error(writer, "characterId is required", http.StatusBadRequest)
			return
		}

		newSandbox.AccountId = accountId

		result := app.DB.Create(&newSandbox)
		if result.Error != nil {
			if strings.Contains(result.Error.Error(), "23505") {
				http.Error(writer, "A sandbox for this characterId already exists", http.StatusConflict)
			} else {
			}
			return
		}

		chunks := terrains.GenerateChunksForSandbox(newSandbox.ID)

		// Print the chunks to the command line
		fmt.Println(chunks)

		chunkCreateResult := app.DB.Create(&chunks)
		if chunkCreateResult.Error != nil {
			http.Error(writer, chunkCreateResult.Error.Error(), http.StatusInternalServerError)
			return
		}

		response, e := json.Marshal(newSandbox)
		if e != nil {
			http.Error(writer, e.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		writer.Write(response)
	}
}
