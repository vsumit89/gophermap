package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *mapHttpHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	value, err := m.mapInstance.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Write([]byte(value))

}
