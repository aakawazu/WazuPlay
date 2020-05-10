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

func generateTokenResponse(id string) ([]byte, error) {
	accessToken, err := random.GenerateRandomString()
	if err != nil {
		return nil, err
	}
	refreshToken, err := random.GenerateRandomString()
	if err != nil {
		return nil, err
	}

	sqlStatement := fmt.Sprintf(
		"INSERT INTO access_token (token, expiration, user_id) VALUES('%s', '%s', '%s')",
		accessToken, db.TimeNow(60), id,
	)
	if _, err := db.RunSQL(sqlStatement); err != nil {
		return nil, err
	}

	sqlStatement = fmt.Sprintf(
		"INSERT INTO refresh_token (token, expiration, user_id) VALUES('%s', '%s', '%s')",
		refreshToken, db.TimeNow(43200), id,
	)
	if _, err := db.RunSQL(sqlStatement); err != nil {
		return nil, err
	}

	res := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
	}

	resjson, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return resjson, nil
}

// TokenHandler /token
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}

	r.ParseForm()
	grantType := r.FormValue("grant_type")

	var resjson []byte

	if grantType == "password" {
		mailaddress := db.EscapeSinglequotation(r.FormValue("mail_address"))
		password := db.EscapeSinglequotation(r.FormValue("password"))

		req := &TokenRequest{
			Mailaddress: mailaddress,
			Password:    password,
		}

		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			httpstates.BadRequest(&w)
			return
		}

		sqlStatement := fmt.Sprintf(
			"SELECT DISTINCT user_id, hashed_password FROM users WHERE mail_address = '%s'",
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

		resjson, err = generateTokenResponse(id)
		if checkerr.InternalServerError(&w, err) {
			return
		}

	} else if grantType == "refresh_token" {
		refreshToken := db.EscapeSinglequotation(r.FormValue("refresh_token"))
		rows, err := db.RunSQL(fmt.Sprintf(
			"SELECT DISTINCT user_id FROM refresh_token WHERE token = '%s'",
			refreshToken,
		))
		defer rows.Close()
		if checkerr.InternalServerError(&w, err) {
			return
		}

		var id string
		for rows.Next() {
			rows.Scan(&id)
		}

		resjson, err = generateTokenResponse(id)
		if checkerr.InternalServerError(&w, err) {
			return
		}

	} else {
		httpstates.BadRequest(&w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resjson))
}
