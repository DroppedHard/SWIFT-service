package swiftCode

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

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
	router.HandleFunc("/swift-codes/{swift-code}", middleware.CustomValidationMiddleware(validateSwiftCode)(h.handleGetSwiftCodeData)).Methods("GET")
	router.HandleFunc("/swift-codes/country/{countryISO2}", middleware.CustomValidationMiddleware(validateCountryCode)(h.handleGetAllSwiftCodesForCountry)).Methods("GET")
	router.HandleFunc("/swift-codes/", middleware.BodyValidationMiddleware(validateAddSwiftCode)(h.handleAddSwiftCode)).Methods("POST")
	router.HandleFunc("/swift-codes/{swift-code}", middleware.CustomValidationMiddleware(validateSwiftCode)(h.handleDeleteSwiftCode)).Methods("DELETE")
}

func (h *Handler) handleGetSwiftCodeData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	swiftCode := mux.Vars(r)["swift-code"]

	bank, err := h.store.GetBankDetailsBySwiftCode(ctx, swiftCode)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("fetching bank details failed: %v", err))
		return
	}
	if bank == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("the SWIFT code %s was not found", swiftCode))
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

func (h *Handler) handleAddSwiftCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payload, ok := ctx.Value(reflect.TypeOf(types.BankDataDetails{})).(types.BankDataDetails)
	if !ok {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve validated payload"))
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
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "bank data succesfully added"})
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
	utils.WriteJson(w, http.StatusCreated, map[string]string{"message": "bank data succesfully deleted"})
}
