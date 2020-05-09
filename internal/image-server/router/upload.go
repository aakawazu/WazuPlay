package router

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
	"github.com/syndtr/goleveldb/leveldb"
)

// ImageFilesRoot image files root
var ImageFilesRoot string = "/wazuplay-files/images"

type uploadImageResponse struct {
	URL string `json:"url"`
}

func checkIfTheImage(file multipart.File) (bool, error) {
	allowedFiletype := []string{"image/png", "image/jpeg"}
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return false, err
	}

	if !mimetype.EqualsAny(mime.String(), allowedFiletype...) {
		return false, nil
	}

	return true, nil
}

func numberOfFiles(folderName string) (int, error) {
	fileList, err := ioutil.ReadDir(fmt.Sprintf("%s/files/%s", ImageFilesRoot, folderName))
	if err != nil {
		return 0, err
	}
	return len(fileList), err
}

// CreateNewFolder create new folder
func CreateNewFolder(db *leveldb.DB) error {
	newFolderName, err := random.GenerateRandomString()
	if err != nil {
		return err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/files/%s", ImageFilesRoot, newFolderName), 0777)
	if err != nil {
		return err
	}

	if err := db.Put([]byte("latest_folder"), []byte(newFolderName), nil); err != nil {
		return err
	}

	return nil
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

// UploadImage /upload
func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	db, err := leveldb.OpenFile(fmt.Sprintf("%s/leveldb/", ImageFilesRoot), nil)
	if checkerr.InternalServerError(err, &w) {
		return
	}
	defer db.Close()

	r.ParseMultipartForm(32 << 20)

	formFile, _, err := r.FormFile("image")
	if checkerr.InternalServerError(err, &w) {
		return
	}
	defer formFile.Close()

	if checkImg, err := checkIfTheImage(formFile); checkerr.InternalServerError(err, &w) {
		return
	} else if !checkImg {
		httpstates.BadRequest(&w)
		return
	}

	fileName, err := random.GenerateRandomString()
	if checkerr.InternalServerError(err, &w) {
		return
	}

	folderName, err := db.Get([]byte("latest_folder"), nil)
	if checkerr.InternalServerError(err, &w) {
		return
	}

	if files, err := numberOfFiles(string(folderName)); checkerr.InternalServerError(err, &w) {
		return
	} else if files > 5 {
		CreateNewFolder(db)
		folderName, err = db.Get([]byte("latest_folder"), nil)
		if checkerr.InternalServerError(err, &w) {
			return
		}
	}

	formFile, _, err = r.FormFile("image")
	if checkerr.InternalServerError(err, &w) {
		return
	}

	if err := saveImage(fileName, string(folderName), formFile, db); checkerr.InternalServerError(err, &w) {
		return
	}

	domain := os.Getenv("DOMAIN")

	res := uploadImageResponse{
		URL: fmt.Sprintf("https://images.%s/%s", domain, fileName),
	}

	resjson, err := json.Marshal(res)
	if checkerr.InternalServerError(err, &w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
