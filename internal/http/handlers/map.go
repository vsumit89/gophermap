package handlers

import (
	"gophermap/internal/services"

	"github.com/go-chi/chi/v5"
)

type mapHttpHandler struct {
	mapInstance *services.Map
	tsLogger    services.TransactionLogger // tsLogger is a transaction logger
}

func RegisterMapHandler(mapInstance *services.Map, router *chi.Mux, ts services.TransactionLogger) {
	mapHandlers := &mapHttpHandler{
		mapInstance: mapInstance,
		tsLogger:    ts,
	}

	router.Get("/map/{key}", mapHandlers.GetKey)
	router.Put("/map/{key}/{value}", mapHandlers.PutKey)

	router.Delete("/map/{key}", mapHandlers.DeleteKey)

}
