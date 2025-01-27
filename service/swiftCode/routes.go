package swiftCode

import (
	"context"
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
	router.HandleFunc("/swift-codes/country/{countryISO2}", h.handleGetAllSwiftCodesForCountry).Methods("GET")
	router.HandleFunc("/swift-codes/", h.handleAddSwiftCode).Methods("POST")
	router.HandleFunc("/swift-codes/{swift-code}", h.handleDeleteSwiftCode).Methods("DELETE")
}

func (h *Handler) handleGetSwiftCodeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)["swift-code"]

	if err := utils.Validate.Var(swiftCode, "required,swiftCode"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the given SWIFT code %s is not correct: %v", swiftCode, err))
		return
	}

	bank, err := h.store.GetBankDetailsBySwiftCode(ctx, swiftCode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("fetching bank details failed: %v", err))
		return
	}
	if bank == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the SWIFT code %s does not exist", swiftCode))
		return
	}
	if bank.IsHeadquarter {
		h.sendBankHqData(w, ctx, bank, swiftCode)
	} else {
		utils.WriteJson(w, http.StatusOK, bank)
	}
}

func (h *Handler) sendBankHqData(w http.ResponseWriter, ctx context.Context, bank *types.BankDataDetails, swiftCode string) {
	branches, partialErr := h.store.GetBranchesDataByHqSwiftCode(ctx, swiftCode)
	bankHq := types.BankHeadquatersResponse{
		BankDataDetails: *bank,
		Branches:        branches,
	}
	if partialErr != nil {
		utils.WriteJson(w, http.StatusPartialContent, map[string]interface{}{
			"data":    bankHq,
			"message": fmt.Sprintf("partial branch data fetched: %v", partialErr),
		})
		return
	}
	utils.WriteJson(w, http.StatusOK, bankHq)
}

func (h *Handler) handleGetAllSwiftCodesForCountry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryCode := mux.Vars(r)["countryISO2"]

	if err := utils.Validate.Var(countryCode, "required,countryISO2"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the given country ISO2 code %s is not correct: %v", countryCode, err))
		return
	}

	banks, partialErr := h.store.GetBanksDataByCountryCode(ctx, countryCode)
	response := types.CountrySwiftCodesResponse{
		CountryIso2: countryCode,
		CountryName: utils.GetCountryNameFromCountryCode(countryCode),
		SwiftCodes:  banks,
	}
	if partialErr != nil {
		utils.WriteJson(w, http.StatusPartialContent, map[string]interface{}{
			"data":    response,
			"message": fmt.Sprintf("failure while fetching banks: %v", partialErr),
		})
		return
	}
	utils.WriteJson(w, http.StatusOK, response)
}

func (h *Handler) handleAddSwiftCode(w http.ResponseWriter, r *http.Request) { // TODO verify data correctiness
	ctx := r.Context()
	var payload types.BankDataDetails
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
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
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to add data: %w", err))
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "bank data succesfully added"}) // TODO
}

func (h *Handler) handleDeleteSwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)["swift-code"]

	if err := utils.Validate.Var(swiftCode, "required,swiftCode"); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the SWIFT code %s is incorrect: %v", swiftCode, err))
		return
	}

	exists, err := h.store.DoesSwiftCodeExist(ctx, swiftCode)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check existence of key %s: %w", swiftCode, err))
		return
	}
	if exists == 0 {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("the SWIFT code %s does not exists", swiftCode))
		return
	}

	if err := h.store.DeleteBankData(ctx, swiftCode); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to delete data: %w", err))
		return
	}
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "bank data succesfully deleted"}) // TODO
}
