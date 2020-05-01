package checkerr

import (
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/httpstates"
)

// InternalServerError error
func InternalServerError(err error, w *http.ResponseWriter) bool {
	if err != nil {
		httpstates.InternalServerError(w)
		return true
	}
	return false
}

// BadRequest error
func BadRequest(err error, w *http.ResponseWriter) bool {
	if err != nil {
		httpstates.BadRequest(w)
		return true
	}
	return false
}
