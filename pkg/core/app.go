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
	return "host=" + POSTGRES_HOST + " user=" + POSTGRES_USER + " password=" + POSTGRES_PASSWORD + " port=" + POSTGRES_PORT + " sslmode=disable TimeZone=America/New_York"
}

func generateDatabaseDsn(dbName string) string {
	return generateDsn() + " dbname=" + dbName
}

func checkDbExistsAndCreate(db *gorm.DB, dbName string) {
	var count int
	db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
	if count == 0 {
		fmt.Println("Database " + dbName + " does not exist, creating database...")
		db.Exec("CREATE DATABASE " + dbName)
	} else {
		fmt.Println("Database " + dbName + " exists, skipping.")
	}

}

func connectToDb(dbName string) *gorm.DB {
	fmt.Println("Connecting to database: " + dbName)
	dsn := generateDatabaseDsn(dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.ErrorHandler(err)
	}
	return db
}

func CreateApp() *App {
	fmt.Println("Starting up app...")
	var dbConnection = connectToDb("postgres")
	var POSTGRES_DATABASE = utils.GetEnv("POSTGRES_DATABASE", "terramajor")
	checkDbExistsAndCreate(dbConnection, POSTGRES_DATABASE)

	return &App{
		DB:     connectToDb(POSTGRES_DATABASE),
		Router: mux.NewRouter(),
	}
}
