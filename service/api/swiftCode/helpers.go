package swiftCode

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/DroppedHard/SWIFT-service/service/api"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
)

func (h *Handler) fetchBankDataBySwiftCode(w http.ResponseWriter, ctx context.Context, swiftCode string) *types.BankDataDetails {
	bank, err := h.store.FindBankDetailsBySwiftCode(ctx, swiftCode)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("fetching bank details failed: %v", err))
		return nil
	}
	if bank == nil {
		api.WriteError(w, http.StatusNotFound, fmt.Errorf("the SWIFT code %s was not found", swiftCode))
		return nil
	}
	return bank
}

func (h *Handler) writeBankHqData(w http.ResponseWriter, ctx context.Context, bank *types.BankDataDetails, swiftCode string) {
	branches, partialErr := h.store.FindBranchesDataByHqSwiftCode(ctx, swiftCode)
	bankHq := types.BankHeadquatersResponse{
		BankDataDetails: *bank,
		Branches:        branches,
	}
	if partialErr != nil {
		api.WriteJson(w, http.StatusPartialContent, bankHq)
		return
	}
	api.WriteJson(w, http.StatusOK, bankHq)
}

func (h *Handler) fetchBankDataByCountryCode(w http.ResponseWriter, ctx context.Context, countryCode string) *types.CountrySwiftCodesResponse {
	banks, partialErr := h.store.FindBanksDataByCountryCode(ctx, countryCode)
	response := types.CountrySwiftCodesResponse{
		CountryIso2: countryCode,
		CountryName: utils.GetCountryNameFromCountryCode(countryCode),
		SwiftCodes:  banks,
	}
	if partialErr != nil {
		api.WriteJson(w, http.StatusPartialContent, response)
		return nil
	}
	return &response
}

func (h *Handler) retrieveValidatedPayloadFromContext(w http.ResponseWriter, ctx context.Context) *types.BankDataDetails {
	payload, ok := ctx.Value(reflect.TypeOf(types.BankDataDetails{})).(types.BankDataDetails)
	if !ok {
		api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve validated payload"))
		return nil
	}
	return &payload
}

func (h *Handler) checkBankDataExistenceInStorage(w http.ResponseWriter, ctx context.Context, swiftCode string, shouldExist bool) bool {
	exists, err := h.store.DoesSwiftCodeExist(ctx, swiftCode)
	if err != nil {
		api.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to check existence of key %s: %w", swiftCode, err))
		return true
	}
	if !shouldExist && exists > 0 {
		api.WriteError(w, http.StatusConflict, fmt.Errorf("the SWIFT code %s already exists", swiftCode))
		return true
	}
	if shouldExist && exists == 0 {
		api.WriteError(w, http.StatusNotFound, fmt.Errorf("the SWIFT code %s does not exist", swiftCode))
		return true
	}
	return false
}
