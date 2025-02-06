package api

import (
	"fmt"
	"net/http"

	_ "github.com/DroppedHard/SWIFT-service/docs"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type HealthCheckHandler struct {
	store types.BankDataStore
}

func NewHealthCheckHandler(store types.BankDataStore) *HealthCheckHandler {
	return &HealthCheckHandler{store: store}
}

func (h *HealthCheckHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", h.getHealthCheck).Methods("GET")
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}

// getHealthCheck godoc
// @Summary 		System health check
// @Description 	endpoint to verify whether system is healthy, or not
// @Tags			status
// @Produce  		json
// @Success	 		200		{object} types.ReturnMessage
// @Router 			/health [get]
func (h *HealthCheckHandler) getHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.store.Ping(ctx)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, fmt.Errorf("redis database unavailable: %v", err))
		return
	}

	WriteMessage(w, http.StatusOK, "OK")
}
