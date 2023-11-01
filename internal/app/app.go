package app

import (
	"gophermap/internal/db"
	"gophermap/pkg/logger"
	"os"

	"github.com/go-chi/chi/v5"
)

type AppConfig struct {
	AppName        string
	Version        string
	PersistentType string
	Database       DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	SSLMode  string
}

// App contains the top level components of the application
// includes the router, database, and other services
type App struct {
	router        *chi.Mux
	db            db.IDatabase
	dataStoreFile *os.File
	logger        logger.ILogger
}

func NewApp(cfg *AppConfig) *App {
	app := &App{}
	app.initRouter()
	app.initLogger()

	if cfg.PersistentType == "file" {
		app.initDataStoreFile()
	} else {
		app.db = db.NewDBService()
		err := app.db.Connect()
		if err != nil {
			app.logger.Fatal("Failed to connect to database")
		}
	}

	
	return app
}

func (a *App) initRouter() {
	a.router = chi.NewRouter()
}

func (a *App) initLogger() {
	a.logger = logger.GetInstance()
}

func (a *App) initDataStoreFile() {
	dataStoreFile, err := os.OpenFile("ds.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		a.logger.Fatal("Failed to open data store file")
	}

	a.dataStoreFile = dataStoreFile
}
