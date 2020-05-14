package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter audio server router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/upload", UploadHandler)
	return router
}

//IndexHandler /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Audio")
}
