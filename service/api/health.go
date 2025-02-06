package api

import (
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/gorilla/mux"
)

type HealthCheckHandler struct {
	store types.BankDataStore
}

func NewHealthCheckHandler(store types.BankDataStore) *HealthCheckHandler {
	return &HealthCheckHandler{store: store}
}

func (h *HealthCheckHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.getHealthCheck).Methods("GET")
}

func (h *HealthCheckHandler) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.store.Ping(ctx)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("redis database unavailable: %v", err))
		return
	}

	WriteMessage(w, http.StatusOK, "OK")
}
