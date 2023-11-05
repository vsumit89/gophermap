package handlers

import (
	"gophermap/pkg/helper"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *mapHttpHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	err := m.mapInstance.Delete(key)
	if err != nil {
		helper.SendJSONResponse(w, http.StatusNotFound, "error", err.Error())
		return
	}

	m.tsLogger.WriteDelete(key)
	
	helper.SendJSONResponse(w, http.StatusOK, "successfully deleted", nil)
}
