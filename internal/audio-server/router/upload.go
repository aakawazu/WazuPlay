package router

import (
	"mime/multipart"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/gabriel-vasile/mimetype"
)

type uploadAudioResponse struct {
	URL string `json:"url"`
}

func checkIfTheAudio(file multipart.File) (bool, error) {
	allowedFiletype := []string{"audio/mpeg", "audio/flac", "audio/wav"}
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return false, err
	}

	if !mimetype.EqualsAny(mime.String(), allowedFiletype...) {
		return false, nil
	}

	return true, nil
}

// UploadHandler /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

}
