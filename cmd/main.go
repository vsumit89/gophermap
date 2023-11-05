package main

import (
	"flag"
	appPkg "gophermap/internal/app"
	"gophermap/pkg/logger"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	// reading the debug flag from the environment
	// if the debug flag is set, the logger will print debug logs
	// if the debug flag is not set, the logger will print info logs
	var debug bool
	flag.BoolVar(&debug, "debug", false, "sets the logger to debug mode")
	flag.Parse()

	log := logger.GetInstance(debug)

	log.Info("Starting Application")
	log.Info("Reading the config file")

	fileBuffer, err := os.ReadFile("../config.yaml")
	if err != nil {
		log.Error("error while reading config file", "error", err.Error())
		return
	}

	var cfg *appPkg.AppConfig
	err = yaml.Unmarshal(fileBuffer, &cfg)
	if err != nil {
		log.Error("error while unmarshalling config file", "error", err.Error())
		return
	}

	// initializes the app and the services it uses
	app := appPkg.NewApp(cfg, log)

	err = app.StartHTTPServer()
	if err != nil {
		log.Error("error while starting HTTP server", "error", err.Error())
		return
	}

}
