package httpstates

import "net/http"

// BadRequest return 400
func BadRequest(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(400)
	w.Write([]byte("400 bad request"))
}

// MethodNotAllowed return 405
func MethodNotAllowed(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(405)
	w.Write([]byte("405 method not allowed"))
}
