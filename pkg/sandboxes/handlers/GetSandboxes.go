package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mw-felker/terra-major-api/pkg/core"
	models "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func selectSandboxes() []models.Sandbox {
	var response []models.Sandbox
	jsonData, err := ioutil.ReadFile("./pkg/instances/handlers/mock-sandboxes.json")
	if err != nil {
		utils.ErrorHandler(err)
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		utils.ErrorHandler(err)
	}

	return response
}

func GetSandboxes(app *core.App) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var Sandboxes = selectSandboxes()
		response, e := json.Marshal(Sandboxes)
		if e != nil {
			log.Panic(e)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
