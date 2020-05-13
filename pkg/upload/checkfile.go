package upload

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"

	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/gabriel-vasile/mimetype"
)

// CheckIfTheAllowedFileType check if the allowed file type
func CheckIfTheAllowedFileType(file multipart.File, allowedFiletype []string) (bool, error) {
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return false, err
	}

	if !mimetype.EqualsAny(mime.String(), allowedFiletype...) {
		return false, nil
	}

	return true, nil
}

// NumberOfFiles number of files
func NumberOfFiles(rootFolder string, folderName string) (int, error) {
	fileList, err := ioutil.ReadDir(fmt.Sprintf("%s/files/%s", rootFolder, folderName))
	if err != nil {
		return 0, err
	}
	return len(fileList), err
}

// CreateNewFolder create new folder
func CreateNewFolder(rootFolder string) (string, error) {
	newFolderName, err := random.GenerateRandomString()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(fmt.Sprintf("%s/files/%s", rootFolder, newFolderName), 0777)
	if err != nil {
		return "", err
	}

	return newFolderName, nil
}
