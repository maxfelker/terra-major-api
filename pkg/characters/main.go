package main

import (
	"fmt"
	"net/http"

	handlers "github.com/mw-felker/centerpoint-instance-api/pkg/characters/handlers"
	"github.com/mw-felker/centerpoint-instance-api/pkg/characters/models"
	core "github.com/mw-felker/centerpoint-instance-api/pkg/core"
	utils "github.com/mw-felker/centerpoint-instance-api/pkg/utils"
)

func seedDb(app *core.App) {
	app.DB.AutoMigrate(&models.Character{})
}

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	app := core.CreateApp()
	seedDb(app)
	app.Router.HandleFunc("/characters", handlers.GetCharacters(app)).Methods("GET")
	app.Router.HandleFunc("/characters/{id}", handlers.GetCharacterById(app)).Methods("GET")
	app.Router.HandleFunc("/characters", handlers.CreateCharacter(app)).Methods("POST")
	app.Router.HandleFunc("/characters/{id}", handlers.UpdateCharacter(app)).Methods("PATCH")
	app.Router.HandleFunc("/characters/{id}", handlers.ArchiveCharacter(app)).Methods("DELETE")
	http.Handle("/", app.Router)
	fmt.Println("Starting API on port " + PORT)
	http.ListenAndServe(":"+PORT, app.Router)
}
