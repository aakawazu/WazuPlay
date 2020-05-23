package router

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/aakawazu/WazuPlay/pkg/storage"
	"github.com/aakawazu/WazuPlay/pkg/token"
	"github.com/aakawazu/WazuPlay/pkg/upload"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/xfrr/goffmpeg/transcoder"
)

// UploadAudioResponse upload audio response
type UploadAudioResponse struct {
	URL string `json:"url"`
}

func convertAudio(file multipart.File) (*os.File, error) {
	tmpFileName, err := random.GenerateRandomString()
	if err != nil {
		return nil, err
	}

	tmpConvertedFileName, err := random.GenerateRandomString()
	if err != nil {
		return nil, err
	}

	tmpConvertedFileName += ".flac"

	f, err := os.Create(fmt.Sprintf("%s/tmp/%s", AudioFilesRoot, tmpFileName))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	io.Copy(f, file)

	trans := new(transcoder.Transcoder)

	err = trans.Initialize(
		fmt.Sprintf("%s/tmp/%s", AudioFilesRoot, tmpFileName),
		fmt.Sprintf("%s/tmp/%s", AudioFilesRoot, tmpConvertedFileName),
	)
	if err != nil {
		return nil, err
	}

	done := trans.Run(false)

	err = <-done
	if err != nil {
		return nil, err
	}

	convertedFile, err := os.Open(fmt.Sprintf("%s/tmp/%s", AudioFilesRoot, tmpConvertedFileName))
	if err != nil {
		return nil, err
	}

	return convertedFile, nil
}

// UploadHandler /upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	owners, err := leveldb.OpenFile(fmt.Sprintf("%s/file-owners/leveldb", AudioFilesRoot), nil)
	if checkerr.InternalServerError(&w, err) {
		return
	}
	defer owners.Close()

	r.ParseMultipartForm(32 << 20)

	accessToken := db.EscapeSinglequotation(r.FormValue("token"))

	userID, err := token.VerificationAccessToken(accessToken)

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

	err = upload.CheckIfTheAllowedFileType(formFile, []string{"audio/mpeg", "audio/flac", "audio/wav"})

	if checkerr.BadRequest(&w, err) {
		return
	}

	convertedAudio, err := convertAudio(formFile)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	file, audioID, err := storage.Create(AudioFilesRoot)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	if err := owners.Put([]byte(audioID), []byte(userID), nil); checkerr.InternalServerError(&w, err) {
		return
	}

	io.Copy(file, convertedAudio)

	scheme := os.Getenv("SCHEME")
	domain := os.Getenv("DOMAIN")

	res := UploadAudioResponse{
		URL: fmt.Sprintf("%s://audio.%s/%s", scheme, domain, audioID),
	}

	resjson, err := json.Marshal(res)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
