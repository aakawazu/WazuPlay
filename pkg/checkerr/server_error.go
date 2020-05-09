package checkerr

import (
	"fmt"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/httpstates"
)

// InternalServerError error
func InternalServerError(w *http.ResponseWriter, err error) bool {
	if err != nil {
		fmt.Println(err)
		httpstates.InternalServerError(w)
		return true
	}
	return false
}

// BadRequest error
func BadRequest(w *http.ResponseWriter, err error) bool {
	if err != nil {
		httpstates.BadRequest(w)
		return true
	}
	return false
}
