package swiftCode

import (
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/service/api"
	"github.com/DroppedHard/SWIFT-service/service/middleware"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/gorilla/mux"
)

type SwiftCodeHandler struct {
	store types.BankDataStore
}

func NewSwiftCodeHandler(store types.BankDataStore) *SwiftCodeHandler {
	return &SwiftCodeHandler{store: store}
}

func (h *SwiftCodeHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/swift-codes/{"+utils.PathParamSwiftCode+"}", middleware.CustomPathParameterValidationMiddleware(api.ValidateSwiftCode)(h.getBankDataBySwiftCode)).Methods("GET")
	router.HandleFunc("/swift-codes/country/{"+utils.PathParamCountryIso2+"}", middleware.CustomPathParameterValidationMiddleware(api.ValidateCountryCode)(h.getBankDataByCountryCode)).Methods("GET")
	router.HandleFunc("/swift-codes/", middleware.BodyValidationMiddleware(api.ValidatePostSwiftCodePayload)(h.postBankData)).Methods("POST")
	router.HandleFunc("/swift-codes/{"+utils.PathParamSwiftCode+"}", middleware.CustomPathParameterValidationMiddleware(api.ValidateSwiftCode)(h.deleteBankData)).Methods("DELETE")
}

func (h *SwiftCodeHandler) getBankDataBySwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)[utils.PathParamSwiftCode]

	bank := h.fetchBankDataBySwiftCode(w, ctx, swiftCode)
	if bank == nil {
		return
	}
	if bank.IsHeadquarter {
		h.writeBankHqData(w, ctx, bank, swiftCode)
	} else {
		api.WriteJson(w, http.StatusOK, bank)
	}
}

func (h *SwiftCodeHandler) getBankDataByCountryCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryCode := mux.Vars(r)[utils.PathParamCountryIso2]

	response := h.fetchBankDataByCountryCode(w, ctx, countryCode)
	if response == nil {
		return
	}

	api.WriteJson(w, http.StatusOK, response)
}

func (h *SwiftCodeHandler) postBankData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := h.retrieveValidatedPayloadFromContext(w, ctx)
	if payload == nil {
		return
	}
	isResponseSent := h.checkBankDataExistenceInStorage(w, ctx, payload.SwiftCode, false)
	if isResponseSent {
		return
	}
	if err := h.store.SaveBankData(ctx, *payload); err != nil {
		api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data: %w", err))
		return
	}
	api.WriteMessage(w, http.StatusCreated, "bank data succesfully added")
}

func (h *SwiftCodeHandler) deleteBankData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)[utils.PathParamSwiftCode]

	isResponseSent := h.checkBankDataExistenceInStorage(w, ctx, swiftCode, true)
	if isResponseSent {
		return
	}

	if err := h.store.DeleteBankData(ctx, swiftCode); err != nil {
		api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete data: %w", err))
		return
	}
	api.WriteMessage(w, http.StatusOK, "bank data succesfully deleted")
}
