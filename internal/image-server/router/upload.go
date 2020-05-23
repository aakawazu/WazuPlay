package router

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	pgdb "github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/storage"
	"github.com/aakawazu/WazuPlay/pkg/token"
	"github.com/aakawazu/WazuPlay/pkg/upload"
	"github.com/nfnt/resize"
)

// ImageFilesRoot image files root
var ImageFilesRoot string = "/wazuplay-files/images"

// UploadImageResponse upload image response
type UploadImageResponse struct {
	URL string `json:"url"`
}

// UploadHandler /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	r.ParseMultipartForm(32 << 20)

	accessToken := pgdb.EscapeSinglequotation(r.FormValue("token"))

	_, err := token.VerificationAccessToken(accessToken)

	switch err {
	case nil:
	case token.ErrTokenNotfound:
		httpstates.BadRequest(&w)
	default:
		httpstates.InternalServerError(&w)
	}

	formFile, _, err := r.FormFile("file")
	if checkerr.InternalServerError(&w, err) {
		return
	}
	defer formFile.Close()

	err = upload.CheckIfTheAllowedFileType(formFile, []string{"image/png", "image/jpeg"})
	switch err {
	case nil:
	case upload.ErrFileTypeNotAllowed:
		httpstates.BadRequest(&w)
		return
	default:
		httpstates.InternalServerError(&w)
		return
	}

	file, imageID, err := storage.Create(ImageFilesRoot)
	defer file.Close()
	if checkerr.InternalServerError(&w, err) {
		return
	}

	formFile, _, err = r.FormFile("file")
	if checkerr.InternalServerError(&w, err) {
		return
	}

	img, _, err := image.Decode(formFile)
	if err != nil {
		return
	}

	resizedImg := resize.Thumbnail(1200, 1200, img, resize.Bicubic)

	if err := jpeg.Encode(file, resizedImg, nil); err != nil {
		return
	}

	scheme := os.Getenv("SCHEME")
	domain := os.Getenv("DOMAIN")

	res := UploadImageResponse{
		URL: fmt.Sprintf("%s://images.%s/%s", scheme, domain, imageID),
	}

	resjson, err := json.Marshal(res)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
