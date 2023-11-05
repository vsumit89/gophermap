package app

import (
	"fmt"
	"gophermap/internal/db"
	"gophermap/pkg/logger"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type AppConfig struct {
	AppName        string   `yaml:"app_name"`
	Version        string   `yaml:"version"`
	PersistentType string   `yaml:"persistent_type"`
	Database       DBConfig `yaml:"database"`
	Port           int      `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSLMode  string `yaml:"ssl_mode"`
}

// App contains the top level components of the application
// includes the router, database, and other services
type App struct {
	router        *chi.Mux
	server        *http.Server
	db            db.IDatabase
	dataStoreFile *os.File
	logger        logger.ILogger
}

func NewApp(cfg *AppConfig, logger logger.ILogger) *App {
	app := &App{}

	app.initRouter()
	app.initLogger(logger)
	logger.Error("test", "config", cfg)
	app.initServer(cfg.Port)

	if cfg.PersistentType == "logfile" {
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

func (a *App) initLogger(logger logger.ILogger) {
	a.logger = logger
}

func (a *App) initDataStoreFile() {
	dataStoreFile, err := os.OpenFile("ds.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		a.logger.Fatal("Failed to open data store file")
	}

	a.dataStoreFile = dataStoreFile
}

func (a *App) initServer(port int) {
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.router,
	}
}

func (a *App) StartHTTPServer() error {
	a.logger.Info("Starting HTTP Server on port " + a.server.Addr)
	err := a.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
