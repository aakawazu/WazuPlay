package router

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/storage"
	"github.com/aakawazu/WazuPlay/pkg/token"
	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
)

// GetHandler /*
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	accessToken := db.EscapeSinglequotation(r.URL.Query().Get("token"))
	if accessToken == "" {
		httpstates.BadRequest(&w)
		return
	}

	userID, err := token.VerificationAccessToken(accessToken)

	fileName := mux.Vars(r)["id"]

	rows, err := db.RunSQL(fmt.Sprintf("SELECT FROM audio_files WHERE id = '%s' AND file_owner = '%s'", fileName, userID))

	if checkerr.InternalServerError(&w, err) {
		return
	}

	if !rows.Next() {
		httpstates.NotFound(&w)
		return
	}

	audioFile, err := storage.Open(AudioFilesRoot, fileName)

	switch err {
	case nil:
	case leveldb.ErrNotFound:
		httpstates.NotFound(&w)
	default:
		httpstates.InternalServerError(&w)
	}

	defer audioFile.Close()

	w.Header().Set("Content-Type", "audio/x-flac")
	io.Copy(w, audioFile)
}
