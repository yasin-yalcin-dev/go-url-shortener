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
	"encoding/json"
	"net/http"
)

// APIError is a custom error structure for API errors
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// New creates a new APIError
func New(code int, message string, detail ...string) *APIError {
	apiErr := &APIError{
		Code:    code,
		Message: message,
	}

	if len(detail) > 0 {
		apiErr.Detail = detail[0]
	}

	return apiErr
}

// Predefined common errors
var (
	ErrBadRequest = &APIError{Code: http.StatusBadRequest, Message: "Invalid request"}
	ErrForbidden  = &APIError{Code: http.StatusForbidden, Message: "Operation not allowed"}
	ErrInternal   = &APIError{Code: http.StatusInternalServerError, Message: "Server error"}
)

// WriteResponse writes the error as an HTTP response
func (e *APIError) WriteResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
}
