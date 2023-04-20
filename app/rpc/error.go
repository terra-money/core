package rpc

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
)

// ErrorResponse defines the attributes of a JSON error response.
type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err1 := w.Write(legacy.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
	if err1 != nil {
		panic(err1)
	}
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}

// CheckError takes care of writing an error response if err is not nil.
// Returns false when err is nil; it returns true otherwise.
func CheckError(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		WriteErrorResponse(w, status, err.Error())
		return true
	}

	return false
}

// CheckInternalServerError attaches an error message to an HTTP 500 INTERNAL SERVER ERROR response.
// Returns false when err is nil; it returns true otherwise.
func CheckInternalServerError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusInternalServerError, err)
}

// PostProcessResponseBare post processes a body similar to PostProcessResponse
// except it does not wrap the body and inject the height.
func PostProcessResponseBare(w http.ResponseWriter, ctx client.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		resp, err = ctx.LegacyAmino.MarshalJSON(body)
		if CheckInternalServerError(w, err) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		panic(err)
	}
}
