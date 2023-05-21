package core

import (
	"github.com/gorilla/mux"
	"github.com/mw-felker/centerpoint-instance-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	Router *mux.Router
}

func createDbClient() *gorm.DB {
	var POSTGRES_HOST = utils.GetEnv("POSTGRES_HOST", "postgres")
	var POSTGRES_PORT = utils.GetEnv("POSTGRES_PORT", "5432")
	var POSTGRES_PASSWORD = utils.GetEnv("POSTGRES_PASSWORD", "postgres")
	var POSTGRES_USER = utils.GetEnv("POSTGRES_USER", "postgres")
	var POSTGRES_DATABASE = utils.GetEnv("POSTGRES_DATABASE", "terramajor")
	dsn := "host=" + POSTGRES_HOST + " user=" + POSTGRES_USER + " password=" + POSTGRES_PASSWORD + " dbname=" + POSTGRES_DATABASE + " port=" + POSTGRES_PORT + " sslmode=disable TimeZone=America/New_York"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.ErrorHandler(err)
	}
	return db
}

// Creates a new App instance and connects to db
func CreateApp() *App {

	app := &App{
		DB:     createDbClient(),
		Router: mux.NewRouter(),
	}

	return app
}
