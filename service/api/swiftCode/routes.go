package swiftCode

import (
	"fmt"
	"net/http"
	"strings"

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

// getBankDataBySwiftCode godoc
// @Summary 		Swift code to bank data
// @Description 	Use it to fetch bank data by SWIFT code - if it is a HQ it branches will be retrieved too
// @Tags		bank
// @Produce  	json
// @Param 		swiftCode 	path 	string 	true 	"Bank swift code"
// @Success	 	200		{object}	types.BankHeadquatersResponse
// @Success	 	206		{object}	types.BankHeadquatersResponse
// @Failure	 	400		{object}	types.ReturnMessage
// @Failure	 	404		{object}	types.ReturnMessage
// @Failure	 	500		{object}	types.ReturnMessage
// @Router 		/swift-codes/{swiftCode} [get]
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

// getBankDataByCountryCode godoc
// @Summary 		Country code to bank data
// @Description 	Use it to fetch banks data by country ISO2 code
// @Tags		bank
// @Produce  	json
// @Param 		countryISO2 	path 	string 	true 	"country ISO2 code"
// @Success	 	200		{object}	types.CountrySwiftCodesResponse
// @Success	 	206		{object}	types.CountrySwiftCodesResponse
// @Failure	 	400		{object}	types.ReturnMessage
// @Failure	 	500		{object}	types.ReturnMessage
// @Router 		/swift-codes/country/{countryISO2} [get]
func (h *SwiftCodeHandler) getBankDataByCountryCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	countryCode := mux.Vars(r)[utils.PathParamCountryIso2]

	response := h.fetchBankDataByCountryCode(w, ctx, strings.ToUpper(countryCode))
	if response == nil {
		return
	}

	api.WriteJson(w, http.StatusOK, response)
}

// postBankData godoc
// @Summary 		Add bank data to the system
// @Description 	Use it to add new bank data - verify data correctiness
// @Tags		bank
// @Accept  	json
// @Produce  	json
// @Param 		bankData 	body 	types.BankDataDetails 	true 	"Bank data"
// @Success	 	201		{object}	types.ReturnMessage
// @Failure	 	400		{object}	types.ReturnMessage
// @Failure	 	409		{object}	types.ReturnMessage
// @Failure	 	500		{object}	types.ReturnMessage
// @Router 		/swift-codes/ [post]
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

// deleteBankData godoc
// @Summary 		Delete bank data from the system
// @Description 	Use it to delete bank data by SWIFT code
// @Tags		bank
// @Produce  	json
// @Param 		swiftCode 	path 	string 	true 	"Bank swift code"
// @Success	 	200		{object}	types.ReturnMessage
// @Failure	 	400		{object}	types.ReturnMessage
// @Failure	 	404		{object}	types.ReturnMessage
// @Failure	 	500		{object}	types.ReturnMessage
// @Router 		/swift-codes/{swiftCode} [delete]
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
