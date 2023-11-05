package http

import (
	"gophermap/internal/http/handlers"
	"gophermap/internal/services"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(mapInstance *services.Map, router *chi.Mux) {
	handlers.RegisterMapHandler(mapInstance, router)
}
