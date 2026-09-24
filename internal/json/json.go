package json

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/sushanthach12/ecom-go/internal/constants"
	"github.com/sushanthach12/ecom-go/internal/httpx"
)

func Decode[T any, PT interface {
	*T
}](w http.ResponseWriter, r *http.Request, logger *slog.Logger) (T, bool) {
	var payload T

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("Malformed Payload:", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid Payload", httpx.CodeMalformedJson)
		return payload, false
	}

	return payload, true
}

func DecodeAndValidate[T any, PT interface {
	*T
	Validate() error
}](w http.ResponseWriter, r *http.Request, logger *slog.Logger) (T, bool) {
	var payload T

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logger.Error("Malformed Payload:", "error", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid Payload", httpx.CodeMalformedJson)
		return payload, false
	}

	if validationErr := PT(&payload).Validate(); validationErr != nil {
		vErr, ok := errors.AsType[*constants.ValidationError](validationErr)
		if ok {
			logger.Error("Validation failed:", "error", vErr.Error())
			httpx.ValidationError(w, http.StatusUnprocessableEntity, vErr.Error(), httpx.CodeValidationFailed, vErr.Field)
			return payload, false
		}

		logger.Error("Validation failed:", "error", validationErr.Error())
		httpx.Error(w, http.StatusBadRequest, "Invalid Payload", httpx.CodeValidationFailed)
		return payload, false
	}

	return payload, true
}
