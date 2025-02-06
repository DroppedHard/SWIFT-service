package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"testing"
	"time"

	"github.com/DroppedHard/SWIFT-service/config"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/stretchr/testify/suite"
)

const (
	startupWaitTime          = 30 * time.Second
	testSwiftCode            = "ALBPPLPWXXX"
	testCountryCode          = "PL"
	healthCheckEndpoint      = "/health"
	swiftCodeEndpoint        = "/swift-codes/"
	swiftCodeCountryEndpoint = "/swift-codes/country/"
)

var (
	baseURL = fmt.Sprintf("http://%s", config.Envs.PublicHost+config.Envs.Port+utils.ApiPrefix)
	wasOn   = false
)

type IntegrationTestSuite struct {
	suite.Suite
	client http.Client
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	if s.waitForService(baseURL + healthCheckEndpoint) {
		fmt.Println("API service detected. Starting tests...")
		wasOn = true
		return
	}
	s.client = http.Client{}

	fmt.Println("Starting Docker Compose...")

	cmd := exec.Command("make", "compose-up")
	cmd.Dir = "../../"
	output, err := cmd.CombinedOutput() // Capture combined output
	if err != nil {
		fmt.Printf("Error running Docker Compose: %s\n", output)
	}
	s.Require().NoError(err, "Failed to start Docker Compose")

	fmt.Printf("Compose started")
	s.waitForService(baseURL + healthCheckEndpoint)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if !wasOn {
		fmt.Println("Shutting Docker Compose down...")
		cmd := exec.Command("make", "compose-down")
		err := cmd.Run()
		s.Require().NoError(err, "Failed to stop Docker Compose")
	} else {
		fmt.Println("Test suite complete - leaving Docker Compose up")
	}
}

func (s *IntegrationTestSuite) waitForService(url string) bool {
	fmt.Printf("Checking for API service at %s...\n", url)
	deadline := time.Now().Add(startupWaitTime)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("Service is up!")
			return true
		}
		fmt.Printf("Service failed to response. Trying again after 2 seconds...\n")
		time.Sleep(2 * time.Second)
	}
	s.FailNow("Service did not start in time")
	return false
}

func (s *IntegrationTestSuite) TestHealthCheck() {
	resp, err := http.Get(baseURL + healthCheckEndpoint)
	s.NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	s.Contains(string(body), `"message":"OK"`)
}

func (s *IntegrationTestSuite) TestFindBankDetailsBySwiftCode() {
	resp, err := http.Get(baseURL + swiftCodeEndpoint + testSwiftCode)
	s.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var bankData map[string]interface{}
		err = json.Unmarshal(body, &bankData)
		s.NoError(err)

		s.Equal(testSwiftCode, bankData["swiftCode"])
		s.NotEmpty(bankData["bankName"])
		s.NotEmpty(bankData["address"])
	} else {
		s.Equal(http.StatusNotFound, resp.StatusCode)
		s.Contains(string(body), "not found")
	}
}

func (s *IntegrationTestSuite) TestFindBanksByCountryCode() {
	resp, err := http.Get(baseURL + swiftCodeCountryEndpoint + testCountryCode)
	s.NoError(err)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	s.Equal(http.StatusOK, resp.StatusCode)

	var banks map[string]interface{}
	err = json.Unmarshal(body, &banks)

	s.NoError(err)
	s.NotEmpty(banks["countryISO2"])
	s.NotEmpty(banks["countryName"])
}

func (s *IntegrationTestSuite) TestDeleteBankData() {

	resp, err := http.Get(baseURL + swiftCodeEndpoint + testSwiftCode)
	s.NoError(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		req, _ := http.NewRequest("DELETE", baseURL+swiftCodeEndpoint+testSwiftCode, nil)
		delResp, err := s.client.Do(req)
		s.NoError(err)
		defer delResp.Body.Close()

		s.Equal(http.StatusOK, delResp.StatusCode)

		afterDelResp, err := http.Get(baseURL + swiftCodeEndpoint + testSwiftCode)
		s.NoError(err)
		defer afterDelResp.Body.Close()

		s.Equal(http.StatusNotFound, afterDelResp.StatusCode)
	} else {
		req, _ := http.NewRequest("DELETE", baseURL+swiftCodeEndpoint+testSwiftCode, nil)
		delResp, err := s.client.Do(req)
		s.NoError(err)
		defer delResp.Body.Close()

		s.Equal(http.StatusNotFound, delResp.StatusCode)
	}
}

func (s *IntegrationTestSuite) TestAddBankData() {
	newBank := map[string]interface{}{
		"SwiftCode":     testSwiftCode,
		"BankName":      "Test Bank",
		"Address":       "123 Test St",
		"isHeadquarter": true,
		"CountryIso2":   testCountryCode,
		"CountryName":   utils.GetCountryNameFromCountryCode(testCountryCode),
	}

	bankJSON, _ := json.Marshal(newBank)
	resp, err := http.Post(baseURL+swiftCodeEndpoint, "application/json", bytes.NewBuffer(bankJSON))
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusCreated, resp.StatusCode)
}
