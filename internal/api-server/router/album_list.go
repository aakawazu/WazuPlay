package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/token"
)

// AlbumListResponse album
type AlbumListResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Artist          string `json:"album_artist"`
	AlbumPictureURL string `json:"album_picture_url"`
}

// AlbumListResponses album list
type AlbumListResponses []AlbumListResponse

// AlbumListHandler return album list
func AlbumListHandler(w http.ResponseWriter, r *http.Request) {
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
	if checkerr.BadRequest(&w, err) {
		return
	}

	rows, err := db.RunSQL(fmt.Sprintf(
		"SELECT DISTINCT id, title, artist, album_picture_url FROM albums WHERE album_owner = '%s'",
		userID,
	))
	defer rows.Close()

	if checkerr.InternalServerError(&w, err) {
		return
	}

	var albumID string
	var title string
	var artist string
	var albumPictureURL string

	var res AlbumListResponses

	for rows.Next() {
		rows.Scan(&albumID, &title, &artist, &albumPictureURL)

		res = append(res, AlbumListResponse{
			ID:              albumID,
			Title:           title,
			Artist:          artist,
			AlbumPictureURL: albumPictureURL,
		})
	}

	resjson, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
