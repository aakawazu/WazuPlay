package router

import (
	"net/http"
	// db "github.com/aakawazu/WazuPlay/pkg/db"
)

// TokenResponse token response structure
type TokenResponse struct {
	accessToken  string
	tokenType    string
	expiresIn    int
	refreshToken string
}

// TokenEndpoint token endpoint
func TokenEndpoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if len(username) == 0 || len(password) == 0 {
			w.WriteHeader(400)
			w.Write([]byte("400 - username and password required"))
		} else {

		}
	}
}
