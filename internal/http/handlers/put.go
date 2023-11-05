package handlers

import (
	"gophermap/pkg/helper"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *mapHttpHandler) PutKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value := chi.URLParam(r, "value")
	m.mapInstance.Put(key, value)

	m.tsLogger.WritePut(key, value)

	helper.SendJSONResponse(w, http.StatusOK, "success", nil)
}
