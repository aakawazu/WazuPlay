package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
)

// GetHandler /*
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb/", ImageFilesRoot), nil)
	if checkerr.InternalServerError(&w, err) {
		return
	}
	defer db.Close()

	fileName := mux.Vars(r)["id"]
	folderName, err := db.Get([]byte(fileName), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			httpstates.NotFound(&w)
			return
		}
		httpstates.InternalServerError(&w)
		return
	}
	img, err := os.Open(fmt.Sprintf("%s/files/%s/%s", ImageFilesRoot, folderName, fileName))

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, img)
}
