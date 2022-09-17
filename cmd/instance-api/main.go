package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	models "github.com/mw-felker/centerpoint-instance-api/pkg/models"
	server "github.com/mw-felker/centerpoint-instance-api/pkg/server"
)

func createInstance(prefabName string, position models.Vector3) models.Instance {
	created := time.Now()
	var ownerId = 1
	var health = 100

	return models.Instance{
		PrefabName: prefabName,
		OwnerId:    ownerId,
		Created:    created,
		Modified:   created, // set modfied same as created (clean)
		Health:     health,
		Position:   position,
	}
}

func errorHandler(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func getInstances() []models.Instance {
	var response []models.Instance
	jsonData, err := ioutil.ReadFile("./data/instances-mock.json")
	if err != nil {
		errorHandler(err)
	}
	err = json.Unmarshal(jsonData, &response)
	if err != nil {
		errorHandler(err)
	}

	return response
}

func getInstancesHandler(writer http.ResponseWriter, request *http.Request) {
	var instances = getInstances()
	response, e := json.Marshal(instances)
	if e != nil {
		log.Panic(e)
	}
	server.Respond(writer, response)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	var PORT = getEnv("PORT", "80")
	var routes = server.Routes{
		{
			Path:    "/instances",
			Method:  "GET",
			Handler: getInstancesHandler,
		},
	}
	server.RegisterRoutes(routes)
	server.Start(PORT)
}
