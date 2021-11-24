package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	models "github.com/mw-felker/centerpoint-instance-api/pkg/models"
)

type Route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type Routes []Route

var routes Routes

const PORT = ":8000"

func defineRoutes() {
	routes = Routes{
		{
			path:    "/instances",
			method:  "GET",
			handler: getInstanceHandler,
		},
	}
}

func jsonResponse(writer http.ResponseWriter, response []byte) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(response)
}

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

func requestFailed(writer http.ResponseWriter, request *http.Request) {
	log.Println(request.Method + " " + request.URL.Path + " is an invalid endpoint")
	http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	writer.Write(nil)
}

func getInstanceHandler(writer http.ResponseWriter, request *http.Request) {

	if !validateRequest(request) {
		requestFailed(writer, request)
	} else {
		var instances = getInstances()
		response, e := json.Marshal(instances)

		if e != nil {
			log.Panic(e)
		}

		jsonResponse(writer, response)
	}

}

func validateRequest(request *http.Request) bool {
	for _, route := range routes {
		if request.URL.Path == route.path {
			if request.Method == route.method {
				return true
			}
		}
	}

	return false
}

func registerRoutes(routes []Route) {
	log.Println("Registering routes...")
	for _, route := range routes {
		http.HandleFunc(route.path, route.handler)
		log.Println(route.method + " " + route.path)
	}
}

func startServer() {
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func main() {
	logMessage := "Starting Centerpoint Instance API service at http://localhost" + PORT
	log.Println(logMessage)
	defineRoutes()
	registerRoutes(routes)
	startServer()
}
