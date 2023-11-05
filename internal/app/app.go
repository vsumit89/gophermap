package app

import (
	"fmt"
	"net/http"

	httpx "gophermap/internal/http"
	"gophermap/internal/services"
	"gophermap/pkg/logger"

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
	Router            *chi.Mux
	server            *http.Server
	logger            logger.ILogger
	MapInstance       *services.Map
	TransactionLogger services.TransactionLogger
}

func NewApp(cfg *AppConfig, logger logger.ILogger) *App {
	app := &App{}

	app.logger = logger
	app.MapInstance = services.NewMap()
	app.TransactionLogger = services.NewTransactionLogger(cfg.PersistentType)

	app.initRouter()

	app.initServer(cfg.Port)

	err := app.TransactionLogger.Init()
	if err != nil {
		app.logger.Error("error while initializing transaction logger", "error", err.Error())
		app.logger.Warn("continuing without transaction logger")
	}

	events, errChan := app.TransactionLogger.ReadEvents()
	// read events from the event channel and update the map
	go func() {
		for {
			select {
			case event := <-events:
				if event.EventType == services.EventPut {
					app.MapInstance.Put(event.Key, event.Value)
				} else {
					app.MapInstance.Delete(event.Key)
				}
			case err := <-errChan:
				if err != nil {
					if err.Error() == "input parse error: EOF" {
						app.logger.Info("keys from file transaction logger successfully read")
					} else {
						app.logger.Error("error while reading events from transaction logger", "error", err.Error())
					}
				}
			}
		}
	}()

	return app
}

func (a *App) initRouter() {
	a.Router = chi.NewRouter()
	a.Router.Use(a.logger.GetHTTPMiddleWare())
	httpx.RegisterRoutes(a.MapInstance, a.Router, a.TransactionLogger)
}

func (a *App) initServer(port int) {
	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.Router,
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
