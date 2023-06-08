package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"

	characters "github.com/mw-felker/terra-major-api/pkg/characters/handlers"
	characterModels "github.com/mw-felker/terra-major-api/pkg/characters/models"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	sandboxes "github.com/mw-felker/terra-major-api/pkg/sandboxes/handlers"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func seedDb(app *core.App) {
	fmt.Println("Seeding app database...")
	app.DB.AutoMigrate(&characterModels.Character{})
	app.DB.AutoMigrate(&sandboxModels.Sandbox{})
	app.DB.AutoMigrate(&sandboxModels.Instance{})
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

	app.Router.HandleFunc("/sandboxes", sandboxes.GetSandboxes(app)).Methods("GET")
	app.Router.HandleFunc("/sandboxes", sandboxes.CreateSandbox(app)).Methods("POST")

	corsObj := handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}))

	http.Handle("/", corsObj(app.Router))

	fmt.Println("Starting terra-major-api on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
