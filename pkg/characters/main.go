package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	handlers "github.com/mw-felker/centerpoint-instance-api/pkg/characters/handlers"
	"github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	router := mux.NewRouter()
	router.HandleFunc("/characters", handlers.GetCharacters).Methods("GET")
	router.HandleFunc("/characters/{id}", handlers.GetCharacterById).Methods("GET")
	router.HandleFunc("/characters", handlers.CreateCharacter).Methods("POST")
	router.HandleFunc("/characters/{id}", handlers.UpdateCharacter).Methods("PATCH")
	router.HandleFunc("/characters/{id}/archive", handlers.ArchiveCharacter).Methods("PATCH")
	http.Handle("/", router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, router)
}
