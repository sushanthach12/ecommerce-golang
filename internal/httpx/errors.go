package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/constants"
)

/*
Status Code		Custom Code
400				invalid_id
404				not_found
500				internal_error
400 			malformed_json
422 			validation_failed
401 			unauthenticated
403 			forbidden
409 			conflict
429 			rate_limited
*/

// ErrorCode is a sealed enum: the underlying string field is unexported,
// so code outside this package cannot construct an ErrorCode directly
// (httpx.ErrorCode{value: "whatever"} won't compile from another package).
// The only valid values are the exported Code* constants below.
type ErrorCode struct {
	value string
}

// String lets ErrorCode be printed/logged normally.
func (c ErrorCode) String() string {
	return c.value
}

// MarshalJSON controls how it's serialized in the response body.
func (c ErrorCode) MarshalJSON() ([]byte, error) { // required: otherwise the value is not serialized inside of json encoder
	return json.Marshal(c.value)
}

var (
	CodeInvalidId        = ErrorCode{"invalid_id"}
	CodeNotFound         = ErrorCode{"not_found"}
	CodeInternalError    = ErrorCode{"internal_error"}
	CodeMalformedJson    = ErrorCode{"malformed_json"}
	CodeValidationFailed = ErrorCode{"validation_failed"}
	CodeUnAuthenticated  = ErrorCode{"unauthenticated"}
	CodeForbidden        = ErrorCode{"forbidden"}
	CodeConflict         = ErrorCode{"conflict"}
	CodeRateLimited      = ErrorCode{"rate_limited"}
	CodeUnauthorized     = ErrorCode{"unauthorized"}
)

type errorPayload struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Field   string    `json:"field,omitempty"` // when empty the field is not shown
}

type errorResponse struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status_code int, message string, code ErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorPayload{
			Code:    code,
			Message: message,
		},
	})
}

func ValidationError(w http.ResponseWriter, status_code int, message string, code ErrorCode, field string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: errorPayload{
			Code:    code,
			Message: message,
			Field:   field,
		},
	})
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ConflictError struct {
	Field   string
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

// HandleError inspects err's concrete type and writes the appropriate HTTP
// response. Falls back to a generic 500 if the error doesn't match any
// known type.
func HandleError(w http.ResponseWriter, err error) {
	var notFoundErr *NotFoundError
	var unauthorizedErr *UnauthorizedError
	var conflictErr *ConflictError
	var validationErr *constants.ValidationError

	switch {
	case errors.As(err, &notFoundErr):
		Error(w, http.StatusNotFound, notFoundErr.Message, CodeNotFound)

	case errors.As(err, &unauthorizedErr):
		Error(w, http.StatusUnauthorized, unauthorizedErr.Message, CodeUnauthorized)

	case errors.As(err, &conflictErr):
		ValidationError(w, http.StatusConflict, conflictErr.Message, CodeConflict, conflictErr.Field)

	case errors.As(err, &validationErr):
		ValidationError(w, http.StatusUnprocessableEntity, validationErr.Message, CodeValidationFailed, validationErr.Field)

	default:
		Error(w, http.StatusInternalServerError, "Something went wrong!", CodeInternalError)
	}
}
