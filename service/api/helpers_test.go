package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DroppedHard/SWIFT-service/service/api"
	"github.com/DroppedHard/SWIFT-service/utils"
	"github.com/stretchr/testify/assert"
)

type HelpersTestCase struct {
	Description  string
	RequestBody  string
	ExpectedErr  error
	StatusCode   int
	Message      string
	ExpectedBody map[string]string
}

func TestParseJson(t *testing.T) {
	tests := []HelpersTestCase{
		{
			Description: "Valid JSON",
			RequestBody: `{"field":"value"}`,
			ExpectedErr: nil,
		},
		{
			Description: "Missing Body",
			RequestBody: "",
			ExpectedErr: fmt.Errorf("EOF"),
		},
		{
			Description: "Invalid JSON",
			RequestBody: `{"field:"value"}`,
			ExpectedErr: fmt.Errorf("invalid character 'v' after object key"),
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBufferString(test.RequestBody))
			var payload map[string]string
			err := api.ParseJson(req, &payload)
			if test.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	tests := []HelpersTestCase{
		{
			Description:  "Write valid JSON response",
			StatusCode:   http.StatusOK,
			ExpectedBody: map[string]string{"message": "Success"},
		},
		{
			Description:  "Write JSON with status created",
			StatusCode:   http.StatusCreated,
			ExpectedBody: map[string]string{"message": "Resource created"},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := api.WriteJson(rr, test.StatusCode, test.ExpectedBody)
			assert.NoError(t, err)

			assert.Equal(t, test.StatusCode, rr.Code)

			var response map[string]string
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, response)
		})
	}
}

func TestWriteMessage(t *testing.T) {
	tests := []HelpersTestCase{
		{
			Description:  "Write message with OK status",
			StatusCode:   http.StatusOK,
			Message:      "Operation successful",
			ExpectedBody: map[string]string{utils.ResponseMessageField: "Operation successful"},
		},
		{
			Description:  "Write message with Bad Request status",
			StatusCode:   http.StatusBadRequest,
			Message:      "Invalid data",
			ExpectedBody: map[string]string{utils.ResponseMessageField: "Invalid data"},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			rr := httptest.NewRecorder()
			api.WriteMessage(rr, test.StatusCode, test.Message)

			assert.Equal(t, test.StatusCode, rr.Code)

			var response map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, response)
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []HelpersTestCase{
		{
			Description:  "Write error with OK status",
			StatusCode:   http.StatusBadRequest,
			Message:      "Invalid request",
			ExpectedBody: map[string]string{utils.ResponseMessageField: "Invalid request"},
		},
		{
			Description:  "Write error with Internal Server Error status",
			StatusCode:   http.StatusInternalServerError,
			Message:      "Server error",
			ExpectedBody: map[string]string{utils.ResponseMessageField: "Server error"},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			rr := httptest.NewRecorder()
			api.WriteError(rr, test.StatusCode, fmt.Errorf("%s", test.Message))

			assert.Equal(t, test.StatusCode, rr.Code)

			var response map[string]string
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, test.ExpectedBody, response)
		})
	}
}
