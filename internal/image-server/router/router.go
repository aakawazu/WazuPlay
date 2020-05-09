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
	router.HandleFunc("/upload", UploadHandler)
	router.HandleFunc("/{id}", GetHandler)
	return router
}

//IndexHandler /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Image")
}
