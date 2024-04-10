package serverError

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ServerError struct {
	Message    string
	StatusCode int
}

func (e ServerError) Error() string {
	return e.Message
}

func ErrorResponse(w http.ResponseWriter, err error) {

	var serverError ServerError
	var httpStatusCode int
	var message string

	ok := errors.As(err, &serverError)
	if !ok {
		httpStatusCode = http.StatusInternalServerError
		message = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)

	resp := make(map[string]string)

	resp["error"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
