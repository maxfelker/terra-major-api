package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Instance struct {
	PrefabId string    `json:"prefabId"`
	OwnerId  int       `json:"ownerId"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Health   int       `json:"health"`
	Position Vector3   `json:"position"`
}

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

func createInstance(prefabName string, position Vector3) Instance {
	created := time.Now()
	var ownerId = 1
	var health = 100

	return Instance{
		prefabName,
		ownerId,
		created,
		created, // set modfied same as created (clean)
		health,
		position,
	}
}

func getInstances() []Instance {
	return []Instance{
		createInstance("Rock5", Vector3{1, 10, 1}),
		createInstance("Rock1", Vector3{10, 1, 1}),
		createInstance("Rock1", Vector3{0, 1, 10}),
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
