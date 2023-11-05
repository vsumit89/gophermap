package handlers

import (
	"gophermap/internal/services"

	"github.com/go-chi/chi/v5"
)

type mapHttpHandler struct {
	mapInstance *services.Map
}

func RegisterMapHandler(mapInstance *services.Map, router *chi.Mux) {
	mapHandlers := &mapHttpHandler{
		mapInstance: mapInstance,
	}

	router.Get("/map/{key}", mapHandlers.GetKey)
	router.Put("/map/{key}/{value}", mapHandlers.PutKey)

}
