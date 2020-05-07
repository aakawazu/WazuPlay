package router

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"github.com/syndtr/goleveldb/leveldb"
)

// UploadImage /upload
func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	r.ParseMultipartForm(32 << 20)
	formFile, _, err := r.FormFile("image")
	if checkerr.InternalServerError(err, &w) {
		return
	}
	defer formFile.Close()

	if checkerr.InternalServerError(err, &w) {
		return
	}

	allowedFiletype := []string{"image/png", "image/jpeg"}
	mime, err := mimetype.DetectReader(formFile)
	if checkerr.InternalServerError(err, &w) {
		return
	}

	if !mimetype.EqualsAny(mime.String(), allowedFiletype...) {
		httpstates.BadRequest(&w)
		return
	}

	formFile, _, err = r.FormFile("image")
	if checkerr.InternalServerError(err, &w) {
		return
	}

	img, _, err := image.Decode(formFile)
	if checkerr.InternalServerError(err, &w) {
		return
	}

	resizedImg := resize.Thumbnail(1200, 1200, img, resize.Bicubic)
	db, err := leveldb.OpenFile("/wazuplay-files/images/leveldb/", nil)
	if checkerr.InternalServerError(err, &w) {
		return
	}
	defer db.Close()

	filesinLatestFolderS, err := db.Get([]byte("filesin_latest_folder"), nil)
	if checkerr.InternalServerError(err, &w) {
		return
	}
	filesinLatestFolder, err := strconv.Atoi(string(filesinLatestFolderS))
	if checkerr.InternalServerError(err, &w) {
		return
	}
	if filesinLatestFolder > 50000 {
		newFolderName, err := random.GenerateRandomString()
		if checkerr.InternalServerError(err, &w) {
			return
		}

		if err := db.Put([]byte("latest_folder"), []byte(newFolderName), nil); checkerr.InternalServerError(err, &w) {
			return
		}

		if err := db.Put([]byte("filesin_latest_folder"), []byte("0"), nil); checkerr.InternalServerError(err, &w) {
			return
		}

		err = os.MkdirAll(fmt.Sprintf("/wazuplay-files/images/files/%s", newFolderName), 0777)
		if checkerr.InternalServerError(err, &w) {
			return
		}
	}
	folderName, err := db.Get([]byte("latest_folder"), nil)
	if checkerr.InternalServerError(err, &w) {
		return
	}

	fileName, err := random.GenerateRandomString()
	if checkerr.InternalServerError(err, &w) {
		return
	}

	out, err := os.Create(fmt.Sprintf("/wazuplay-files/images/files/%s/%s", folderName, fileName))
	if checkerr.InternalServerError(err, &w) {
		return
	}
	defer out.Close()

	if err := jpeg.Encode(out, resizedImg, nil); checkerr.InternalServerError(err, &w) {
		return
	}

	if err := db.Put([]byte(fileName), []byte(folderName), nil); checkerr.InternalServerError(err, &w) {
		return
	}

}
