package userdata

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateRandomString generate random string
func GenerateRandomString() (string, error) {
	u0, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	u1, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	s := fmt.Sprintf("%s%s", u0, u1)
	return s, nil
}
