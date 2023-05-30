package core

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

func generateDsn() string {
	var POSTGRES_HOST = utils.GetEnv("POSTGRES_HOST", "postgres")
	var POSTGRES_PORT = utils.GetEnv("POSTGRES_PORT", "5432")
	var POSTGRES_PASSWORD = utils.GetEnv("POSTGRES_PASSWORD", "postgres")
	var POSTGRES_USER = utils.GetEnv("POSTGRES_USER", "postgres")
	var POSTGRES_DATABASE = utils.GetEnv("POSTGRES_DATABASE", "terramajor")
	return "host=" + POSTGRES_HOST + " user=" + POSTGRES_USER + " password=" + POSTGRES_PASSWORD + " port=" + POSTGRES_PORT + " database=" + POSTGRES_DATABASE + " sslmode=disable TimeZone=America/New_York"
}

func createDbClient() *gorm.DB {
	dsn := generateDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.ErrorHandler(err)
	}
	return db
}

func seedDb(app *App) {
	fmt.Printf("Running initial database setup.......")
	app.DB.Exec("DROP DATABASE IF EXISTS terramajor")
	app.DB.Exec("CREATE DATABASE terramajor")
}

// Creates a new App instance and connects to db
func CreateApp() *App {
	app := &App{
		DB:     createDbClient(),
		Router: mux.NewRouter(),
	}
	//seedDb(app)
	return app
}
