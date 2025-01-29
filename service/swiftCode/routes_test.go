package swiftCode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DroppedHard/SWIFT-service/service/swiftCode"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

func TestSwiftCodeServiceHandlers(t *testing.T) {
	store := &mockSwiftCodeStore{}
	handler := swiftCode.NewHandler(store)
	router := mux.NewRouter()
	handler.RegisterRoutes(router)
	t.Run("should fail if the given payload is invalid", func(t *testing.T) {

		payload := types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				Address:       "",
				BankName:      "",
				CountryISO2:   "",
				IsHeadquarter: false,
				SwiftCode:     "",
			},
			CountryName: "",
		}
		marshalled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/swift-codes/", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockSwiftCodeStore struct {
	mock.Mock
}

func (m *mockSwiftCodeStore) GetBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*types.BankDataDetails, error) {
	return nil, nil
}
func (m *mockSwiftCodeStore) GetBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]types.BankDataCore, error) {
	return nil, nil
}
func (m *mockSwiftCodeStore) GetBanksDataByCountryCode(ctx context.Context, countryCode string) ([]types.BankDataCore, error) {
	return nil, nil
}
func (m *mockSwiftCodeStore) AddBankData(ctx context.Context, data types.BankDataDetails) error {
	return nil
}
func (m *mockSwiftCodeStore) DeleteBankData(ctx context.Context, swiftCode string) error {
	return nil
}
func (m *mockSwiftCodeStore) DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error) {
	return 0, nil
}
