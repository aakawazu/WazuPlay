package httpstates

import "net/http"

// InternalServerError return 500
func InternalServerError(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(500)
	w.Write([]byte("500 - Internal Server Error"))
}
