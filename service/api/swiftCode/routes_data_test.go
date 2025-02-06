package swiftCode_test

import (
	"fmt"
	"net/http"

	"github.com/DroppedHard/SWIFT-service/types"
	"github.com/DroppedHard/SWIFT-service/utils"
)

type GetSwiftCodeTestCase struct {
	Description       string
	SwiftCode         string
	ExpectedData      *types.BankHeadquatersResponse
	ExpectedCode      int
	IsHeadquarter     bool
	ErrorIncludes     string
	NegativeFindError error
	NegativeFindValue *types.BankDataDetails
}

var GetBankDataBySwiftCodePositiveTestCases = []GetSwiftCodeTestCase{
	{
		Description:   "Valid HQ SWIFT Code (1 branch)",
		SwiftCode:     "ALBPPLPWXXX",
		IsHeadquarter: true,
		ExpectedCode:  http.StatusOK,
		ExpectedData: &types.BankHeadquatersResponse{
			BankDataDetails: types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					BankName:      "Headquarters Bank",
					CountryIso2:   "PL",
					IsHeadquarter: true,
					Address:       "HQ Street 1",
				},
				CountryName: utils.GetCountryNameFromCountryCode("PL"),
			},
			Branches: []types.BankDataCore{
				{
					SwiftCode:     "ALBPPLPWCUS",
					BankName:      "Branch 1",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 1",
				},
			},
		},
	},
	{
		Description:   "Valid Branch SWIFT Code",
		SwiftCode:     "ALBPPLPWCUS",
		IsHeadquarter: false,
		ExpectedCode:  http.StatusOK,
		ExpectedData: &types.BankHeadquatersResponse{
			BankDataDetails: types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWCUS",
					BankName:      "Branch 1",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 1",
				},
				CountryName: utils.GetCountryNameFromCountryCode("PL"),
			},
		},
	},
	{
		Description:   "Valid HQ SWIFT Code (3 branches)",
		SwiftCode:     "ALBPPLPWXXX",
		IsHeadquarter: true,
		ExpectedCode:  http.StatusOK,
		ExpectedData: &types.BankHeadquatersResponse{
			BankDataDetails: types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					BankName:      "Headquarters Bank",
					CountryIso2:   "PL",
					IsHeadquarter: true,
					Address:       "HQ Street 1",
				},
				CountryName: utils.GetCountryNameFromCountryCode("PL"),
			},
			Branches: []types.BankDataCore{
				{
					SwiftCode:     "ALBPPLPW001",
					BankName:      "Branch 1",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 1",
				},
				{
					SwiftCode:     "ALBPPLPW002",
					BankName:      "Branch 2",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 2",
				},
				{
					SwiftCode:     "ALBPPLPW003",
					BankName:      "Branch 3",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 3",
				},
			},
		},
	},
	{
		Description:   "Valid HQ SWIFT Code (No Branches)",
		SwiftCode:     "ALBPPLPWXXX",
		IsHeadquarter: true,
		ExpectedCode:  http.StatusOK,
		ExpectedData: &types.BankHeadquatersResponse{
			BankDataDetails: types.BankDataDetails{
				BankDataCore: types.BankDataCore{
					SwiftCode:     "ALBPPLPWXXX",
					BankName:      "Headquarters Bank",
					CountryIso2:   "PL",
					IsHeadquarter: true,
					Address:       "HQ Street 1",
				},
				CountryName: "Poland",
			},
			Branches: nil,
		},
	},
}
var GetBankDataBySwiftCodeNegativeTestCases = []GetSwiftCodeTestCase{
	{
		Description:   "Invalid Swift Code (Nonexistent)",
		SwiftCode:     "ALBPPLXAXXX",
		IsHeadquarter: true,
		ExpectedCode:  http.StatusNotFound,
		ExpectedData:  nil,
		ErrorIncludes: "the SWIFT code",
	},
	{
		Description:   "Invalid Swift Code (Invalid Characters)",
		SwiftCode:     "ALBPPL__123",
		IsHeadquarter: false,
		ExpectedCode:  http.StatusBadRequest,
		ExpectedData:  nil,
		ErrorIncludes: "validation failed on 'swiftCode' tag",
	},
	{
		Description:   "Invalid Swift Code (Too Short)",
		SwiftCode:     "ALBPPL",
		IsHeadquarter: false,
		ExpectedCode:  http.StatusBadRequest,
		ExpectedData:  nil,
		ErrorIncludes: "validation failed on 'swiftCode' tag",
	},
	{
		Description:       "Internal server error",
		SwiftCode:         "ALBPPLXAXXX",
		IsHeadquarter:     false,
		ExpectedCode:      http.StatusInternalServerError,
		ExpectedData:      nil,
		NegativeFindError: fmt.Errorf("internal server error message"),
		ErrorIncludes:     "internal server error message",
	},
}

type CountryCodeTestCase struct {
	Description       string
	CountryCode       string
	ExpectedData      *types.CountrySwiftCodesResponse
	ExpectedCode      int
	ErrorIncludes     string
	NegativeFindError error
	NegativeFindValue []types.BankDataCore
}

var GetBankDataByCountryCodePositiveTestCases = []CountryCodeTestCase{
	{
		Description:  "Valid Country Code (1 result)",
		CountryCode:  "PL",
		ExpectedCode: http.StatusOK,
		ExpectedData: &types.CountrySwiftCodesResponse{
			CountryIso2: "PL",
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
			SwiftCodes: []types.BankDataCore{
				{
					SwiftCode:     "ALBPPLPWCUS",
					BankName:      "Branch 1",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 1",
				},
			},
		},
	},
	{
		Description:  "Valid Country Code (3 results)",
		CountryCode:  "PL",
		ExpectedCode: http.StatusOK,
		ExpectedData: &types.CountrySwiftCodesResponse{
			CountryIso2: "PL",
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
			SwiftCodes: []types.BankDataCore{
				{
					SwiftCode:     "ALBPPLPW001",
					BankName:      "Branch 1",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 1",
				},
				{
					SwiftCode:     "ALBPPLPW002",
					BankName:      "Branch 2",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 2",
				},
				{
					SwiftCode:     "ALBPPLPW003",
					BankName:      "Branch 3",
					CountryIso2:   "PL",
					IsHeadquarter: false,
					Address:       "Branch Street 3",
				},
			},
		},
	},
}

var GetBankDataByCountryCodeNegativeTestCases = []CountryCodeTestCase{
	{
		Description:   "Invalid Country Code (too long)",
		CountryCode:   "AAAAAAA",
		ExpectedCode:  http.StatusBadRequest,
		ExpectedData:  nil,
		ErrorIncludes: "validation failed on 'countryISO2' tag",
	},
	{
		Description:   "Invalid Country Code (too short)",
		CountryCode:   "X",
		ExpectedCode:  http.StatusBadRequest,
		ExpectedData:  nil,
		ErrorIncludes: "validation failed on 'countryISO2' tag",
	},
	{
		Description:   "Invalid Country Code (incorrect example)",
		CountryCode:   "__",
		ExpectedCode:  http.StatusBadRequest,
		ExpectedData:  nil,
		ErrorIncludes: "validation failed on 'countryISO2' tag",
	},
	{
		Description:       "Internal server error (partial return)",
		CountryCode:       "PL",
		ExpectedCode:      http.StatusPartialContent,
		ExpectedData:      nil,
		NegativeFindValue: nil,
		NegativeFindError: fmt.Errorf("internal server error message"),
		ErrorIncludes:     "",
	},
}

type PostBankTestCase struct {
	Description        string
	BankData           types.BankDataDetails
	MessageIncludes    string
	ExpectedCode       int
	NegativeExistValue int64
	NegativeExistError error
	NegativeSaveError  error
}

var PostBankDataPositiveTestCases = []PostBankTestCase{
	{
		Description: "Valid bank HQ data",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWXXX",
				BankName:      "Headquarters Bank",
				CountryIso2:   "PL",
				IsHeadquarter: true,
				Address:       "HQ Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusCreated,
		MessageIncludes: "bank data succesfully added",
	},
	{
		Description: "Valid bank branch data",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWCUS",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusCreated,
		MessageIncludes: "bank data succesfully added",
	},
}

var PostBankDataNegativeTestCases = []PostBankTestCase{
	{
		Description: "Invalid bank data (swift code 1)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWaaa",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "failed on the 'swiftCode' tag",
	},
	{
		Description: "Invalid bank data (swift code 2)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "A__P__PWXXX",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "failed on the 'swiftCode' tag",
	},
	{
		Description: "Invalid bank data (country code)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPW",
				BankName:      "Branch Bank",
				CountryIso2:   "XXX",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "failed on the 'countryISO2' tag",
	},
	{
		Description: "Invalid bank data (isHeadquarter)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWXXX",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "isHeadquarter value 'false' does not match the swiftCode value 'ALBPPLPWXXX'",
	},
	{
		Description: "Internal server error (exist check)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWCUS",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		NegativeExistError: fmt.Errorf("error message"),
		ExpectedCode:       http.StatusInternalServerError,
		MessageIncludes:    "error message",
	},
	{
		Description: "Swift code already exists",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWCUS",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		NegativeExistValue: 1,
		ExpectedCode:       http.StatusConflict,
		MessageIncludes:    "the SWIFT code ALBPPLPWCUS already exists",
	},
	{
		Description: "Internal server error (save)",
		BankData: types.BankDataDetails{
			BankDataCore: types.BankDataCore{
				SwiftCode:     "ALBPPLPWCUS",
				BankName:      "Branch Bank",
				CountryIso2:   "PL",
				IsHeadquarter: false,
				Address:       "Branch Street 1",
			},
			CountryName: utils.GetCountryNameFromCountryCode("PL"),
		},
		NegativeSaveError: fmt.Errorf("error message"),
		ExpectedCode:      http.StatusInternalServerError,
		MessageIncludes:   "error message",
	},
}

type DeleteSwiftCodeTestCase struct {
	Description         string
	SwiftCode           string
	ExpectedCode        int
	MessageIncludes     string
	NegativeExistValue  int64
	NegativeExistError  error
	NegativeDeleteError error
}

var DeleteBankDataPositiveTestCases = []DeleteSwiftCodeTestCase{
	{
		Description:     "Valid Swift Code HQ",
		SwiftCode:       "ALBPPLPWXXX",
		ExpectedCode:    http.StatusOK,
		MessageIncludes: "bank data succesfully deleted",
	},
	{
		Description:     "Valid Swift Code branch",
		SwiftCode:       "ALBPPLPWCUS",
		ExpectedCode:    http.StatusOK,
		MessageIncludes: "bank data succesfully deleted",
	},
}
var DeleteBankDataNegativeTestCases = []DeleteSwiftCodeTestCase{
	{
		Description:     "Invalid Swift Code",
		SwiftCode:       "ALBPPL__XXX",
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "validation failed on 'swiftCode' tag",
	},
	{
		Description:     "Invalid Swift Code length",
		SwiftCode:       "ALBPPL__",
		ExpectedCode:    http.StatusBadRequest,
		MessageIncludes: "validation failed on 'swiftCode' tag",
	},
	{
		Description:        "Internal server error (exist check)",
		SwiftCode:          "ALBPPLPWXXX",
		ExpectedCode:       http.StatusInternalServerError,
		NegativeExistError: fmt.Errorf("internal server error message"),
		MessageIncludes:    "internal server error message",
	},
	{
		Description:        "Swift code does not exist",
		SwiftCode:          "ALBPPLPWXXX",
		ExpectedCode:       http.StatusNotFound,
		NegativeExistValue: 0,
		MessageIncludes:    "the SWIFT code ALBPPLPWXXX does not exist",
	},
	{
		Description:         "Internal server error (delete)",
		SwiftCode:           "ALBPPLPWXXX",
		ExpectedCode:        http.StatusInternalServerError,
		NegativeExistValue:  1,
		NegativeDeleteError: fmt.Errorf("internal server error message"),
		MessageIncludes:     "internal server error message",
	},
}
