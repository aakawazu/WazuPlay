package router

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	pgdb "github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/aakawazu/WazuPlay/pkg/token"
	"github.com/aakawazu/WazuPlay/pkg/upload"
	"github.com/nfnt/resize"
	"github.com/syndtr/goleveldb/leveldb"
)

// ImageFilesRoot image files root
var ImageFilesRoot string = "/wazuplay-files/images"

type uploadImageResponse struct {
	URL string `json:"url"`
}

func saveImage(fileName string, folderName string, formFile multipart.File, db *leveldb.DB) error {
	img, _, err := image.Decode(formFile)
	if err != nil {
		return err
	}

	resizedImg := resize.Thumbnail(1200, 1200, img, resize.Bicubic)

	out, err := os.Create(fmt.Sprintf("%s/files/%s/%s", ImageFilesRoot, folderName, fileName))
	if err != nil {
		return err
	}
	defer out.Close()

	if err := jpeg.Encode(out, resizedImg, nil); err != nil {
		return err
	}

	if err := db.Put([]byte(fileName), []byte(folderName), nil); err != nil {
		return err
	}
	return nil
}

// UploadHandler /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb/", ImageFilesRoot), nil)
	if checkerr.InternalServerError(&w, err) {
		return
	}
	defer db.Close()

	r.ParseMultipartForm(32 << 20)

	accessToken := pgdb.EscapeSinglequotation(r.FormValue("token"))

	if _, err := token.VerificationAccessToken(accessToken); err != nil {
		if err == token.ErrTokenNotfound {
			httpstates.BadRequest(&w)
			return
		}
		httpstates.InternalServerError(&w)
		return
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

	fileName, err := random.GenerateRandomString()
	if checkerr.InternalServerError(&w, err) {
		return
	}

	folderName, err := db.Get([]byte("latest_folder"), nil)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	if files, err := upload.NumberOfFiles(ImageFilesRoot, string(folderName)); checkerr.InternalServerError(&w, err) {
		return
	} else if files > 50000 {
		if folderName, err := upload.CreateNewFolder(ImageFilesRoot); checkerr.InternalServerError(&w, err) {
			return
		} else if err := db.Put([]byte("latest_folder"), []byte(folderName), nil); checkerr.InternalServerError(&w, err) {
			return
		}

		folderName, err = db.Get([]byte("latest_folder"), nil)
		if checkerr.InternalServerError(&w, err) {
			return
		}
	}

	formFile, _, err = r.FormFile("image")
	if checkerr.InternalServerError(&w, err) {
		return
	}

	if err := saveImage(fileName, string(folderName), formFile, db); checkerr.InternalServerError(&w, err) {
		return
	}

	scheme := os.Getenv("SCHEME")
	domain := os.Getenv("DOMAIN")

	res := uploadImageResponse{
		URL: fmt.Sprintf("%s://images.%s/%s", scheme, domain, fileName),
	}

	resjson, err := json.Marshal(res)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
