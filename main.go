package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"

	accounts "github.com/mw-felker/terra-major-api/pkg/accounts/handlers"
	accountModels "github.com/mw-felker/terra-major-api/pkg/accounts/models"
	auth "github.com/mw-felker/terra-major-api/pkg/auth/handlers"
	characters "github.com/mw-felker/terra-major-api/pkg/characters/handlers"
	characterModels "github.com/mw-felker/terra-major-api/pkg/characters/models"
	core "github.com/mw-felker/terra-major-api/pkg/core"
	sandboxes "github.com/mw-felker/terra-major-api/pkg/sandboxes/handlers"
	sandboxModels "github.com/mw-felker/terra-major-api/pkg/sandboxes/models"
	terrains "github.com/mw-felker/terra-major-api/pkg/terrains/handlers"
	terrainModels "github.com/mw-felker/terra-major-api/pkg/terrains/models"
	utils "github.com/mw-felker/terra-major-api/pkg/utils"
)

func seedDb(app *core.App) {
	fmt.Println("Seeding app database...")
	app.DB.AutoMigrate(&accountModels.Account{})
	app.DB.AutoMigrate(&characterModels.Character{})
	app.DB.AutoMigrate(&sandboxModels.Sandbox{})
	app.DB.AutoMigrate(&sandboxModels.Instance{})
	app.DB.AutoMigrate(&terrainModels.TerrainChunk{})
}

func main() {
	var PORT = utils.GetEnv("PORT", "8000")
	app := core.CreateApp()
	seedDb(app)

	// Auth
	app.Router.HandleFunc("/login", accounts.Login(app)).Methods("POST")
	app.Router.HandleFunc("/tokens", auth.CreateUnityClientToken(app)).Methods("POST")

	// User-centric (me) routes
	app.Router.HandleFunc("/signup", accounts.CreateMyAccount(app)).Methods("POST")
	app.Router.HandleFunc("/me", accounts.GetMyAccount(app)).Methods("GET")
	app.Router.HandleFunc("/my/password", accounts.UpdatePassword(app)).Methods("PATCH")
	app.Router.HandleFunc("/my/sandboxes", sandboxes.GetMySandboxes(app)).Methods("GET")
	app.Router.HandleFunc("/my/characters", characters.GetMyCharacters(app)).Methods("GET")

	// Accounts
	app.Router.HandleFunc("/accounts", accounts.CreateAccount(app)).Methods("POST")
	app.Router.HandleFunc("/accounts/{id}", accounts.GetAccountById(app)).Methods("GET")
	app.Router.HandleFunc("/accounts/{id}", accounts.UpdateAccount(app)).Methods("PATCH")

	// Characters
	app.Router.HandleFunc("/characters", characters.CreateCharacter(app)).Methods("POST")
	app.Router.HandleFunc("/characters/{id}", characters.GetCharacterById(app)).Methods("GET")
	app.Router.HandleFunc("/characters/{id}", characters.UpdateCharacter(app)).Methods("PATCH")
	app.Router.HandleFunc("/characters/{id}", characters.ArchiveCharacter(app)).Methods("DELETE")

	// Sandboxes
	app.Router.HandleFunc("/sandboxes", sandboxes.GetSandboxes(app)).Methods("GET")
	app.Router.HandleFunc("/sandboxes", sandboxes.CreateSandbox(app)).Methods("POST")
	app.Router.HandleFunc("/sandboxes/{id}", sandboxes.GetSandboxById(app)).Methods("GET")
	app.Router.HandleFunc("/sandboxes/{id}", sandboxes.ArchiveSandbox(app)).Methods("DELETE")

	// Instances
	app.Router.HandleFunc("/sandboxes/{sandboxId}/instances", sandboxes.GetInstancesBySandboxId(app)).Methods("GET")
	app.Router.HandleFunc("/sandboxes/{sandboxId}/instances/{instanceId}", sandboxes.GetInstanceById(app)).Methods("GET")
	app.Router.HandleFunc("/sandboxes/{sandboxId}/instances", sandboxes.CreateInstance(app)).Methods("POST")
	app.Router.HandleFunc("/sandboxes/{sandboxId}/instances/{instanceId}", sandboxes.UpdateInstance(app)).Methods("PATCH")
	app.Router.HandleFunc("/sandboxes/{sandboxId}/instances/{instanceId}", sandboxes.ArchiveInstance(app)).Methods("DELETE")

	app.Router.HandleFunc("/terrain", terrains.GetTerrain(app)).Methods("GET")
	app.Router.HandleFunc("/terrain", terrains.CreateTerrainChunk(app)).Methods("POST")

	corsObj := handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}))

	http.Handle("/", corsObj(app.Router))

	fmt.Println("Starting terra-major-api on port " + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
