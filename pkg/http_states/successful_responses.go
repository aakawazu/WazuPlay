package statescode

import "net/http"

// OK return 200
func OK(wa *http.ResponseWriter) {
	w := *wa
	w.WriteHeader(200)
	w.Write([]byte("200 - OK"))
}
