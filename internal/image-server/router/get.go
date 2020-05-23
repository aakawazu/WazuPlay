package router

import (
	"io"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/storage"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
)

// GetHandler /*
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	fileName := mux.Vars(r)["id"]

	img, err := storage.Open(ImageFilesRoot, fileName)
	if err != nil {
		if err == leveldb.ErrNotFound {
			httpstates.NotFound(&w)
			return
		}
		httpstates.InternalServerError(&w)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}
