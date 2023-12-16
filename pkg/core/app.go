package core

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gorilla/mux"
	"github.com/mw-felker/terra-major-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	DB     *gorm.DB
	NoSQL  *azcosmos.DatabaseClient
	Router *mux.Router
}

func generateDsn() string {
	var POSTGRES_HOST = utils.GetEnv("POSTGRES_HOST", "postgres")
	var POSTGRES_PORT = utils.GetEnv("POSTGRES_PORT", "5432")
	var POSTGRES_PASSWORD = utils.GetEnv("POSTGRES_PASSWORD", "postgres")
	var POSTGRES_USER = utils.GetEnv("POSTGRES_USER", "postgres")
	var POSTGRES_SSL = utils.GetEnv("POSTGRES_SSL", "require")
	return "host=" + POSTGRES_HOST + " user=" + POSTGRES_USER + " password=" + POSTGRES_PASSWORD + " port=" + POSTGRES_PORT + " sslmode=" + POSTGRES_SSL + " TimeZone=America/New_York"
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

func connectToCosmos() (*azcosmos.DatabaseClient, error) {
	COSMOS_DB_PORT := utils.GetEnv("COSMOS_DB_PORT", "8081")
	COSMOS_DB_HOST := utils.GetEnv("COSMOS_DB_HOST", "cosmos")
	url := "https://" + COSMOS_DB_HOST + ":" + COSMOS_DB_PORT
	COSMOS_DB_PRIMARY_KEY := utils.GetEnv("COSMOS_DB_PRIMARY_KEY")
	COSMOS_DB_NAME := utils.GetEnv("COSMOS_DB_NAME", "terramajor")

	cred, err := azcosmos.NewKeyCredential(COSMOS_DB_PRIMARY_KEY)
	if err != nil {
		log.Println("Failed to create Cosmos DB key credential:", err)
		return nil, err
	}

	client, err := azcosmos.NewClientWithKey(url, cred, nil)
	if err != nil {
		log.Println("Failed to create Cosmos client:", err)
		return nil, err
	}

	database, err := client.NewDatabase(COSMOS_DB_NAME)
	if err != nil {
		log.Println("Failed to create Cosmos DB database client:", err)
		return nil, err
	}

	return database, nil
}

func CreateApp() *App {
	fmt.Println("Starting up app...")
	var dbConnection = connectToDb("postgres")
	var POSTGRES_DATABASE = utils.GetEnv("POSTGRES_DATABASE", "terramajor")
	checkDbExistsAndCreate(dbConnection, POSTGRES_DATABASE)

	db := connectToDb(POSTGRES_DATABASE)
	db.Logger = logger.Default.LogMode(logger.Silent)
	//db.Logger.LogMode(logger.Silent)

	noSqlClient, error := connectToCosmos()
	if error != nil {
		log.Println("Failed to get Cosmos DB database client:", error)
		return nil
	}

	return &App{
		DB:     connectToDb(POSTGRES_DATABASE),
		NoSQL:  noSqlClient,
		Router: mux.NewRouter(),
	}
}
