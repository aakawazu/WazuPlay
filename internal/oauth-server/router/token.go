package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// TokenRequest token request struct
type TokenRequest struct {
	Mailaddress string `validate:"required,email"`
	Password    string `validate:"required"`
}

// TokenResponse token response struct
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func validateTokenRequest(req *TokenRequest) bool {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return false
	}
	return true
}

// TokenHandler /token
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	r.ParseForm()
	mailaddress := db.EscapeSinglequotation(r.FormValue("mail_address"))
	password := db.EscapeSinglequotation(r.FormValue("password"))

	req := &TokenRequest{
		Mailaddress: mailaddress,
		Password:    password,
	}

	if !validateTokenRequest(req) {
		httpstates.BadRequest(&w)
		return
	}

	sqlStatement := fmt.Sprintf(
		"SELECT DISTINCT id, hashed_password FROM users WHERE mail_address = '%s'",
		mailaddress,
	)
	rows, err := db.RunSQL(sqlStatement)
	defer rows.Close()
	if checkerr.InternalServerError(&w, err) {
		return
	}

	var id string
	var hashedPassword string
	for rows.Next() {
		rows.Scan(&id, &hashedPassword)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); checkerr.BadRequest(&w, err) {
		return
	}

	accessToken, err := random.GenerateRandomString()
	if checkerr.InternalServerError(&w, err) {
		return
	}
	refreshToken, err := random.GenerateRandomString()
	if checkerr.InternalServerError(&w, err) {
		return
	}

	sqlStatement = fmt.Sprintf(
		"INSERT INTO access_token (token, expiration, user_id) VALUES('%s', '%s', '%s')",
		accessToken, db.TimeNow(60), id,
	)
	if _, err := db.RunSQL(sqlStatement); checkerr.InternalServerError(&w, err) {
		return
	}

	sqlStatement = fmt.Sprintf(
		"INSERT INTO refresh_token (token, expiration, user_id) VALUES('%s', '%s', '%s')",
		refreshToken, db.TimeNow(43200), id,
	)
	if _, err := db.RunSQL(sqlStatement); checkerr.InternalServerError(&w, err) {
		return
	}

	res := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}
	resjson, err := json.Marshal(res)
	if checkerr.InternalServerError(&w, err) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
