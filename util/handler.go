package util

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
)

//
// ErrorResponder responds to the passed in http.ResponsWriter
// using the given statusCode and error message
//
func ErrorResponder(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if statusCode >= 500 {
		sentry.CaptureException(err)
	}
	if err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
		panic(err)
	}
}

//
// JSONResponder responds to the passed in http.ResponsWriter with
// a json object of the passed in interface
//
func JSONResponder(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
