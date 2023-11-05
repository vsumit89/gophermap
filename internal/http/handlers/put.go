package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *mapHttpHandler) PutKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value := chi.URLParam(r, "value")
	m.mapInstance.Put(key, value)

	w.WriteHeader(http.StatusCreated)
}
