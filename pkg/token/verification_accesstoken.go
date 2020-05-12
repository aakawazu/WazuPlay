package token

import (
	"errors"
	"fmt"

	"github.com/aakawazu/WazuPlay/pkg/db"
)

var (
	// ErrTokenNotfound token not found
	ErrTokenNotfound = errors.New("token not found")
)

// VerificationAccessToken verification access token. return userID
func VerificationAccessToken(token string) (string, error) {
	rows, err := db.RunSQL(fmt.Sprintf(
		"SELECT DISTINCT user_id FROM access_token WHERE token = '%s' and expiration > '%s'",
		token, db.TimeNow(0),
	))
	if err != nil {
		return "", err
	}

	var id string
	for rows.Next() {
		rows.Scan(&id)
	}

	if id == "" {
		return "", ErrTokenNotfound
	}

	return id, nil
}
