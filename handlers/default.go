package handlers

import "net/http"

//
// MethodNotAllowedHandler is used when an error occurs
// outside of a handler
//
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(405)
	w.Write([]byte("Method not allowed"))
	return
}

//
// NotFoundHandler is used when an error occurs
// outside of a handler
//
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Not found"))
	return
}
