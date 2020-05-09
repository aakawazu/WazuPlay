package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter oauth server router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/token", TokenHandler)
	router.HandleFunc("/verificationcode/generate", GenerateVerificationCodeHandler)
	router.HandleFunc("/verificationcode/confirm", ConfirmVerificationCodeHandler)
	router.HandleFunc("/signup", SignUpHandler)
	return router
}

// IndexHandler /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, OAuth")
}
