package swiftCode

import (
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/service/middleware"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.BankDataStore
}

func NewHandler(store types.BankDataStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/swift-codes/{"+utils.PathParamSwiftCode+"}", middleware.CustomPathParameterValidationMiddleware(validateSwiftCode)(h.getBankDataBySwiftCode)).Methods("GET")
	router.HandleFunc("/swift-codes/country/{"+utils.PathParamCountryIso2+"}", middleware.CustomPathParameterValidationMiddleware(validateCountryCode)(h.getBankDataByCountryCode)).Methods("GET")
	router.HandleFunc("/swift-codes/", middleware.BodyValidationMiddleware(validateAddSwiftCode)(h.postBankData)).Methods("POST")
	router.HandleFunc("/swift-codes/{"+utils.PathParamSwiftCode+"}", middleware.CustomPathParameterValidationMiddleware(validateSwiftCode)(h.deleteBankData)).Methods("DELETE")
}

func (h *Handler) getBankDataBySwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)[utils.PathParamSwiftCode]

	bank := h.fetchBankDataBySwiftCode(w, ctx, swiftCode)
	if bank == nil {
		return
	}
	if bank.IsHeadquarter {
		h.writeBankHqData(w, ctx, bank, swiftCode)
	} else {
		utils.WriteJson(w, http.StatusOK, bank)
	}
}

func (h *Handler) getBankDataByCountryCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryCode := mux.Vars(r)[utils.PathParamCountryIso2]

	response := h.fetchBankDataByCountryCode(w, ctx, countryCode)
	if response == nil {
		return
	}

	utils.WriteJson(w, http.StatusOK, response)
}

func (h *Handler) postBankData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload := h.retrieveValidatedPayloadFromContext(w, ctx)
	if payload == nil {
		return
	}
	isResponseSent := h.checkBankDataExistenceInStorage(w, ctx, payload.SwiftCode, false)
	if isResponseSent {
		return
	}
	if err := h.store.AddBankData(ctx, *payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data: %w", err))
		return
	}
	utils.WriteMessage(w, http.StatusCreated, "bank data succesfully added")
}

func (h *Handler) deleteBankData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)[utils.PathParamSwiftCode]

	isResponseSent := h.checkBankDataExistenceInStorage(w, ctx, swiftCode, true)
	if isResponseSent {
		return
	}
	
	if err := h.store.DeleteBankData(ctx, swiftCode); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete data: %w", err))
		return
	}
	utils.WriteMessage(w, http.StatusCreated, "bank data succesfully deleted")
}
