package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	models "github.com/mw-felker/terra-major-api/pkg/instances/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func getInstances() []models.Instance {
	var response []models.Instance
	jsonData, err := ioutil.ReadFile("./pkg/instances/handlers/instances-mock.json")
	if err != nil {
		utils.ErrorHandler(err)
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		utils.ErrorHandler(err)
	}

	return response
}

func GetInstances(writer http.ResponseWriter, request *http.Request) {
	var Instances = getInstances()
	response, e := json.Marshal(Instances)
	if e != nil {
		log.Panic(e)
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
