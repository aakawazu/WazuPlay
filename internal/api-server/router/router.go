package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter api server router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	// router.HandleFunc("/album/register", RegisterAlbumHandler)
	router.HandleFunc("/album/list", AlbumListHandler)
	router.HandleFunc("/album/{id}", AlbumInfoHandler)
	// router.HandleFunc("/music/{id}", MusicInfoHandler)
	return router
}

// IndexHandler /
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, WazuPlay")
}
