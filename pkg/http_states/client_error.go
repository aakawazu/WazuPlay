package statescode

import "net/http"

// BadRequest return 400
func BadRequest(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(400)
	w.Write([]byte("400 - Bad Request"))
}

// MethodNotAllowed return 405
func MethodNotAllowed(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(405)
	w.Write([]byte("405 - Method Not Allowed"))
}
