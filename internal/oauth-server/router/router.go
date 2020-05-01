package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter oauth server router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", Index)
	router.HandleFunc("/token", TokenEndpoint)
	router.HandleFunc("/verificationcode/generate", GenerateVerificationCode)
	router.HandleFunc("/verificationcode/confirm", ConfirmVerificationCode)
	router.HandleFunc("/signup", SignUp)
	return router
}

//Index index
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, OAuth")
}
