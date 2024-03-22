package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type errResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"some error message"`
} //@name Error

var (
	ErrBadRequest         = errors.New("bad request")
	ErrNoAuth             = errors.New("received no auth info")
	ErrEmptyRequest       = errors.New("received empty request")
	ErrParse              = errors.New("failed to parse parameter")
	ErrAuthFailed         = errors.New("failed to authenticate")
	ErrInvalidDateFormat  = errors.New("received invalid date")
	ErrInvalidInterval    = errors.New("end date is beforehand the start date or matches it")
	ErrExpiredDate        = errors.New("date is expired")
	ErrIncompleteInterval = errors.New("received no start date or no end date")
	ErrNoUserID           = errors.New("received no user id")
	ErrIncompleteRequest  = errors.New("in case telegram account is changed, both id and nickname should be set")

	ValidateErr = new(validator.ValidationErrors)
)

func WriteWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := errResponse{
		Status:  statusCode,
		Message: errMsg,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Printf("Failed to encode error response into JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// WriteWithStatus sets the response header to application/json, write the header
// with a status code, and encodes and writes the data json.NewEncoder()
func WriteWithStatus(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("failed to encode API response into JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func WriteValidationError(w http.ResponseWriter, errs validator.ValidationErrors) {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	WriteWithError(w, http.StatusBadRequest, strings.Join(errMsgs, ", "))
}
