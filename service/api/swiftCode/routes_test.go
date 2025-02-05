package swiftCode_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/DroppedHard/SWIFT-service/service/api/swiftCode"
	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RoutesTestSuite struct {
	suite.Suite
	router *mux.Router
	store  *mockSwiftCodeStore
}

func (suite *RoutesTestSuite) SetupTest() {
	suite.store = new(mockSwiftCodeStore)
	handler := swiftCode.NewHandler(suite.store)

	suite.router = mux.NewRouter()
	handler.RegisterRoutes(suite.router)
}

func (suite *RoutesTestSuite) makeRequest(method, url string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, nil)
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)
	return rr
}

func (suite *RoutesTestSuite) makePostRequest(url string, body interface{}) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	suite.router.ServeHTTP(rr, req)
	return rr
}

func (suite *RoutesTestSuite) assertJSONResponse(rr *httptest.ResponseRecorder, expectedCode int, expectedData interface{}) {
	suite.Equal(expectedCode, rr.Code)
	if expectedData != nil {
		responseType := reflect.TypeOf(expectedData).Elem()
		response := reflect.New(responseType).Interface()
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		suite.NoError(err)
		suite.Equal(expectedData, response)
	}
}

func (suite *RoutesTestSuite) assertMessageResponse(rr *httptest.ResponseRecorder, expectedCode int, expectedMessage string) {
	suite.Equal(expectedCode, rr.Code)
	var errorResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	suite.NoError(err)
	suite.Contains(errorResponse[utils.ResponseMessageField], expectedMessage)
}

func (suite *RoutesTestSuite) resetMocks() {
	suite.store.ExpectedCalls = nil
	suite.store.Calls = nil
}

func TestRoutesSuite(t *testing.T) {
	suite.Run(t, &RoutesTestSuite{})
}

func (suite *RoutesTestSuite) TestGetBankDataBySwiftCode() {
	suite.Run("Positive Cases", func() {
		for _, testCase := range GetBankDataBySwiftCodePositiveTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(utils.GetFunctionName(
					types.BankDataStore.FindBankDetailsBySwiftCode),
					mock.Anything,
					testCase.SwiftCode,
				).Return(&testCase.ExpectedData.BankDataDetails, nil)

				if testCase.IsHeadquarter {
					suite.store.On(
						utils.GetFunctionName(types.BankDataStore.FindBranchesDataByHqSwiftCode),
						mock.Anything,
						testCase.SwiftCode,
					).Return(testCase.ExpectedData.Branches, nil)
				}
				defer suite.resetMocks()

				rr := suite.makeRequest("GET", "/swift-codes/"+testCase.SwiftCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				if testCase.IsHeadquarter {
					suite.assertJSONResponse(rr, testCase.ExpectedCode, &testCase.ExpectedData)
				} else {
					suite.assertJSONResponse(rr, testCase.ExpectedCode, &testCase.ExpectedData.BankDataDetails)
				}
				suite.store.AssertExpectations(suite.T())
			})
		}
	})
	suite.Run("Negative Cases", func() {
		for _, testCase := range GetBankDataBySwiftCodeNegativeTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(
					utils.GetFunctionName(types.BankDataStore.FindBankDetailsBySwiftCode),
					mock.Anything,
					testCase.SwiftCode,
				).Return(testCase.NegativeFindValue, testCase.NegativeFindError).Maybe()
				defer suite.resetMocks()

				rr := suite.makeRequest("GET", "/swift-codes/"+testCase.SwiftCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.ErrorIncludes)

				suite.store.AssertExpectations(suite.T())
			})
		}
	})

}

func (suite *RoutesTestSuite) TestGetBankDataByCountryCode() {
	suite.Run("Positive Cases", func() {
		for _, testCase := range GetBankDataByCountryCodePositiveTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(
					utils.GetFunctionName(types.BankDataStore.FindBanksDataByCountryCode),
					mock.Anything,
					testCase.CountryCode,
				).Return(testCase.ExpectedData.SwiftCodes, nil)
				defer suite.resetMocks()

				rr := suite.makeRequest("GET", "/swift-codes/country/"+testCase.CountryCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertJSONResponse(rr, testCase.ExpectedCode, &testCase.ExpectedData)
				suite.store.AssertExpectations(suite.T())
			})
		}
	})
	suite.Run("Negative Cases", func() {
		for _, testCase := range GetBankDataByCountryCodeNegativeTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(
					utils.GetFunctionName(types.BankDataStore.FindBanksDataByCountryCode),
					mock.Anything,
					testCase.CountryCode,
				).Return(testCase.NegativeFindValue, testCase.NegativeFindError).Maybe()
				defer suite.resetMocks()
				rr := suite.makeRequest("GET", "/swift-codes/country/"+testCase.CountryCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.ErrorIncludes)

				suite.store.AssertExpectations(suite.T())
			})
		}
	})
}

func (suite *RoutesTestSuite) TestPostBankData() {
	suite.Run("Positive Cases", func() {
		for _, testCase := range PostBankDataPositiveTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DoesSwiftCodeExist), mock.Anything, testCase.BankData.SwiftCode).Return(int64(0), nil)
				suite.store.On(utils.GetFunctionName(types.BankDataStore.SaveBankData), mock.Anything, testCase.BankData).Return(nil)
				defer suite.resetMocks()

				rr := suite.makePostRequest("/swift-codes/", testCase.BankData)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.MessageIncludes)
				suite.store.AssertExpectations(suite.T())
			})
		}
	})
	suite.Run("Negative Cases", func() {
		for _, testCase := range PostBankDataNegativeTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DoesSwiftCodeExist), mock.Anything, testCase.BankData.SwiftCode).Return(testCase.NegativeExistValue, testCase.NegativeExistError).Maybe()
				suite.store.On(utils.GetFunctionName(types.BankDataStore.SaveBankData), mock.Anything, testCase.BankData).Return(testCase.NegativeSaveError).Maybe()
				defer suite.resetMocks()

				rr := suite.makePostRequest("/swift-codes/", testCase.BankData)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.MessageIncludes)

				suite.store.AssertExpectations(suite.T())
			})
		}
	})
}

func (suite *RoutesTestSuite) TestDeleteBankData() {
	suite.Run("Positive Cases", func() {
		for _, testCase := range DeleteBankDataPositiveTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DoesSwiftCodeExist), mock.Anything, testCase.SwiftCode).Return(int64(1), nil)
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DeleteBankData), mock.Anything, testCase.SwiftCode).Return(nil)
				defer suite.resetMocks()

				rr := suite.makeRequest("DELETE", "/swift-codes/"+testCase.SwiftCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.MessageIncludes)
				suite.store.AssertExpectations(suite.T())
			})
		}
	})
	suite.Run("Negative Cases", func() {
		for _, testCase := range DeleteBankDataNegativeTestCases {
			suite.Run(testCase.Description, func() {
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DoesSwiftCodeExist), mock.Anything, testCase.SwiftCode).Return(testCase.NegativeExistValue, testCase.NegativeExistError).Maybe()
				suite.store.On(utils.GetFunctionName(types.BankDataStore.DeleteBankData), mock.Anything, testCase.SwiftCode).Return(testCase.NegativeDeleteError).Maybe()
				defer suite.resetMocks()

				rr := suite.makeRequest("DELETE", "/swift-codes/"+testCase.SwiftCode)

				suite.Equal(testCase.ExpectedCode, rr.Code)

				suite.assertMessageResponse(rr, testCase.ExpectedCode, testCase.MessageIncludes)

				suite.store.AssertExpectations(suite.T())
			})
		}
	})
}

type mockSwiftCodeStore struct {
	mock.Mock
}

func (m *mockSwiftCodeStore) FindBankDetailsBySwiftCode(ctx context.Context, swiftCode string) (*types.BankDataDetails, error) {
	args := m.Called(ctx, swiftCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.BankDataDetails), args.Error(1)
}
func (m *mockSwiftCodeStore) FindBranchesDataByHqSwiftCode(ctx context.Context, swiftCode string) ([]types.BankDataCore, error) {
	args := m.Called(ctx, swiftCode)
	return args.Get(0).([]types.BankDataCore), args.Error(1)
}
func (m *mockSwiftCodeStore) FindBanksDataByCountryCode(ctx context.Context, countryCode string) ([]types.BankDataCore, error) {
	args := m.Called(ctx, countryCode)
	return args.Get(0).([]types.BankDataCore), args.Error(1)
}
func (m *mockSwiftCodeStore) SaveBankData(ctx context.Context, data types.BankDataDetails) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}
func (m *mockSwiftCodeStore) DeleteBankData(ctx context.Context, swiftCode string) error {
	args := m.Called(ctx, swiftCode)
	return args.Error(0)
}
func (m *mockSwiftCodeStore) DoesSwiftCodeExist(ctx context.Context, swiftCode string) (int64, error) {
	args := m.Called(ctx, swiftCode)
	return args.Get(0).(int64), args.Error(1)
}
