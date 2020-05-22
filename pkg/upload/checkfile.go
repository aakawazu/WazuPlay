package upload

import (
	"errors"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
)

var (
	// ErrFileTypeNotAllowed file type not allowed
	ErrFileTypeNotAllowed = errors.New("file type not allowed")
)

// CheckIfTheAllowedFileType check if the allowed file type
func CheckIfTheAllowedFileType(file multipart.File, allowedFiletype []string) error {
	mime, err := mimetype.DetectReader(file)
	if err != nil {
		return err
	}

	if !mimetype.EqualsAny(mime.String(), allowedFiletype...) {
		return ErrFileTypeNotAllowed
	}

	return nil
}
