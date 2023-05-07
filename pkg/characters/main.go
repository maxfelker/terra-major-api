package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mw-felker/centerpoint-instance-api/pkg/characters/handlers"
	"github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func registerRoutes(router *mux.Router) {
	fmt.Println("Registering routes...")
	router.HandleFunc("/characters", handlers.GetCharactersHandler).Methods("GET")
	router.HandleFunc("/characters", handlers.CreateCharacterHandler).Methods("POST")
}

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	router := mux.NewRouter()
	registerRoutes(router)
	http.Handle("/", router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, router)
}
