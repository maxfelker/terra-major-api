package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	models "github.com/mw-felker/centerpoint-instance-api/pkg/models"
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

	//update content type
	writer.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func registerRoutes(router *mux.Router) {
	fmt.Println("Registering routes...")
	router.HandleFunc("/instances", getInstancesHandler).Methods("GET")
}

func main() {
	var PORT = getEnv("PORT", "8000")
	router := mux.NewRouter()
	registerRoutes(router)
	http.Handle("/", router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, router)
}
