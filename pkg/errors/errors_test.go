/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package errors

import (
	"net/http"
	"testing"
)

// TestNewAPIError tests the creation of a new API error
func TestNewAPIError(t *testing.T) {
	testCases := []struct {
		name           string
		code           int
		message        string
		detail         string
		expectedCode   int
		expectedMsg    string
		expectedDetail string
	}{
		{
			name:           "Standard Error",
			code:           http.StatusBadRequest,
			message:        "Invalid input",
			detail:         "The provided input is incorrect",
			expectedCode:   http.StatusBadRequest,
			expectedMsg:    "Invalid input",
			expectedDetail: "The provided input is incorrect",
		},
		{
			name:           "Error without Detail",
			code:           http.StatusInternalServerError,
			message:        "Server error",
			expectedCode:   http.StatusInternalServerError,
			expectedMsg:    "Server error",
			expectedDetail: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiErr *APIError
			if tc.detail != "" {
				apiErr = New(tc.code, tc.message, tc.detail)
			} else {
				apiErr = New(tc.code, tc.message)
			}

			if apiErr.Code != tc.expectedCode {
				t.Errorf("Expected code %d, got %d", tc.expectedCode, apiErr.Code)
			}

			if apiErr.Message != tc.expectedMsg {
				t.Errorf("Expected message %s, got %s", tc.expectedMsg, apiErr.Message)
			}

			if apiErr.Detail != tc.expectedDetail {
				t.Errorf("Expected detail %s, got %s", tc.expectedDetail, apiErr.Detail)
			}
		})
	}
}

// TestAPIErrorImplementsError checks if APIError implements the error interface
func TestAPIErrorImplementsError(t *testing.T) {
	apiErr := New(http.StatusBadRequest, "Test error")

	// Check if it can be used as an error
	var err error = apiErr
	if err == nil {
		t.Error("APIError should implement the error interface")
	}

	// Check Error() method
	if apiErr.Error() != apiErr.Message {
		t.Errorf("Error() should return the message, got %s", apiErr.Error())
	}
}
