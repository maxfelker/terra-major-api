package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/mw-felker/terra-major-api/pkg/instances/handlers"
	"github.com/mw-felker/terra-major-api/pkg/utils"
)

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	router := mux.NewRouter()
	router.HandleFunc("/instances", handlers.GetInstances).Methods("GET")
	http.Handle("/", router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, router)
}
