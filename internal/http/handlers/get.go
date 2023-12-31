package handlers

import (
	"gophermap/pkg/helper"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *mapHttpHandler) GetKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	value, err := m.mapInstance.Get(key)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusNotFound, "error", err.Error())
		return
	}

	helper.SendJSONResponse(w, http.StatusOK, "success", value)
}
