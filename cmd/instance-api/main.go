package main

import (
	"encoding/json"
	"log"
	"net/http"
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

func getInstances() []models.Instance {

	return []models.Instance{
		createInstance("Rock5", models.Vector3{X: 1, Y: 10, Z: 1}),
		createInstance("Rock1", models.Vector3{X: 10, Y: 1, Z: 1}),
		createInstance("Rock1", models.Vector3{X: 0, Y: 1, Z: 10}),
	}
}

func getInstancesHandler(writer http.ResponseWriter, request *http.Request) {
	var instances = getInstances()
	response, e := json.Marshal(instances)
	if e != nil {
		log.Panic(e)
	}
	server.Respond(writer, response)
}

func main() {
	const PORT = ":8000"
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
