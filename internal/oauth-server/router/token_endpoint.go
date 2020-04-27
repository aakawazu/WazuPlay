package router

import (
	"net/http"

	db "github.com/aakawazu/WazuPlay/pkg/db"
	httpStates "github.com/aakawazu/WazuPlay/pkg/http_states"
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
		username := db.EscapeSinglequotation(r.FormValue("username"))
		password := db.EscapeSinglequotation(r.FormValue("password"))
		if len(username) == 0 || len(password) == 0 {
			httpStates.BadRequest(&w)
		} else {

		}
	} else {
		httpStates.MethodNotAllowed(&w)
	}
}
