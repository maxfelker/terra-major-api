package main

import (
	"fmt"
	"net/http"

	handlers "github.com/mw-felker/centerpoint-instance-api/pkg/characters/handlers"
	core "github.com/mw-felker/centerpoint-instance-api/pkg/core"
	"github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	app := core.CreateApp()
	app.Router.HandleFunc("/characters", handlers.GetCharacters(app)).Methods("GET")
	app.Router.HandleFunc("/characters/{id}", handlers.GetCharacterById).Methods("GET")
	app.Router.HandleFunc("/characters", handlers.CreateCharacter).Methods("POST")
	app.Router.HandleFunc("/characters/{id}", handlers.UpdateCharacter).Methods("PATCH")
	app.Router.HandleFunc("/characters/{id}/archive", handlers.ArchiveCharacter).Methods("PATCH")
	http.Handle("/", app.Router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, app.Router)
}
