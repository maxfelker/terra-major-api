package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"

	characters "github.com/mw-felker/terra-major-api/pkg/characters/handlers"
	"github.com/mw-felker/terra-major-api/pkg/characters/models"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func seedDb(app *core.App) {
	app.DB.AutoMigrate(&models.Character{})
}

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	app := core.CreateApp()
	seedDb(app)
	app.Router.HandleFunc("/characters", characters.GetCharacters(app)).Methods("GET")
	app.Router.HandleFunc("/characters/{id}", characters.GetCharacterById(app)).Methods("GET")
	app.Router.HandleFunc("/characters", characters.CreateCharacter(app)).Methods("POST")
	app.Router.HandleFunc("/characters/{id}", characters.UpdateCharacter(app)).Methods("PATCH")
	app.Router.HandleFunc("/characters/{id}", characters.ArchiveCharacter(app)).Methods("DELETE")

	corsObj := handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}))

	http.Handle("/", corsObj(app.Router))

	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
