package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/token"
	"github.com/gorilla/mux"
)

// MusicListResponse music
type MusicListResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"album_artist"`
	Instrument string `json:"instrument"`
}

// MusicListResponses music list
type MusicListResponses []MusicListResponse

// AlbumInfoHandler return music list
func AlbumInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	accessToken := db.EscapeSinglequotation(r.Header.Get("token"))
	if accessToken == "" {
		httpstates.BadRequest(&w)
		return
	}

	userID, err := token.VerificationAccessToken(accessToken)
	if checkerr.BadRequest(&w, err) {
		return
	}

	albumID := mux.Vars(r)["id"]

	rows, err := db.RunSQL(fmt.Sprintf(
		"SELECT DISTINCT id, title, artist, instrument FROM audio_files WHERE file_owner = '%s' AND album = '%s'",
		userID, albumID,
	))
	defer rows.Close()

	if checkerr.InternalServerError(&w, err) {
		return
	}

	var musicID string
	var title string
	var artist string
	var instrument string

	var res MusicListResponses

	for rows.Next() {
		rows.Scan(&musicID, &title, &artist, &instrument)

		res = append(res, MusicListResponse{
			ID:         musicID,
			Title:      title,
			Artist:     artist,
			Instrument: instrument,
		})
	}

	resjson, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
