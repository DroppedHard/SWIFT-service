package swiftCode

import (
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.BankDataStore
}

func NewHandler(store types.BankDataStore) *Handler {
	return &Handler{store: store}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/swift-codes/{swift-code}", h.handleGetSwiftCodeData).Methods("GET")
	router.HandleFunc("/swift-codes/country/{countryISO2}", h.handleGetAllSwiftcodesForCountry).Methods("GET")
	router.HandleFunc("/swift-codes/", h.handleAddSwiftCode).Methods("POST")
	router.HandleFunc("/swift-codes/{swift-code}", h.handleDeleteSwiftCode).Methods("DELETE")
}

func (h *Handler) handleGetSwiftCodeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)["swift-code"]

	// validate the SWIFT code - maybe more details later on
	if err := utils.Validate.Var(swiftCode, "required,swiftCode"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
        return
	}

	bank, err := h.store.GetBankDetailsBySwiftCode(ctx, swiftCode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	utils.WriteJSON(w, http.StatusOK, bank)
}

func (h *Handler) handleGetAllSwiftcodesForCountry(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleAddSwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload types.BankDataDetails
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	exists, err := h.store.DoesSwiftCodeExist(ctx, payload.SwiftCode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check existence of key %s: %w", payload.SwiftCode, err))
		return
	}
	if exists > 0 {
		utils.WriteError(w, http.StatusConflict, fmt.Errorf("the SWIFT code %s already exists", payload.SwiftCode))
		return
	}

	if err := h.store.AddBankData(ctx, payload); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data to DB: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "bank data succesfully added"})
}

func (h *Handler) handleDeleteSwiftCode(w http.ResponseWriter, r *http.Request) {

}