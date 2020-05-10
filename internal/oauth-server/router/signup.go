package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aakawazu/WazuPlay/pkg/checkerr"
	"github.com/aakawazu/WazuPlay/pkg/db"
	"github.com/aakawazu/WazuPlay/pkg/httpstates"
	"github.com/aakawazu/WazuPlay/pkg/mail"
	"github.com/aakawazu/WazuPlay/pkg/random"
	"github.com/go-playground/validator/v10"
)

// GenerateVerificationCodeRequest generate verificationcode request
type GenerateVerificationCodeRequest struct {
	MailAddress string `validate:"required,email"`
}

// ConfirmVerificationCodeRequest confirm verification code request
type ConfirmVerificationCodeRequest struct {
	MailAddress      string `validate:"required,email"`
	VerificationCode int    `validate:"gte=0,lt=999999"`
}

// SignUpRequest sign up request
type SignUpRequest struct {
	MailAddress      string `validate:"required,email"`
	VerificationCode int    `validate:"min=1,max=999999"`
	Username         string `validate:"required"`
	Password         string `validate:"required,min=5,max=50"`
}

func findMailAddressDuplicate(w *http.ResponseWriter, mailAddress string) bool {
	rows, err := db.RunSQL(fmt.Sprintf(
		"SELECT * FROM users WHERE mail_address = '%s'",
		mailAddress,
	))
	defer rows.Close()
	if checkerr.InternalServerError(w, err) {
		return true
	}
	if rows.Next() {
		httpstates.BadRequest(w)
		return true
	}
	return false
}

func findVerificationCode(w *http.ResponseWriter, mailAddress string, verificationCode int) bool {
	rows, err := db.RunSQL(fmt.Sprintf(
		"SELECT * FROM pending WHERE mail_address = '%s' and verification_code = %d and expiration > '%s'",
		mailAddress, verificationCode, db.TimeNow(0),
	))
	defer rows.Close()
	if checkerr.InternalServerError(w, err) {
		return false
	}
	if !rows.Next() {
		httpstates.BadRequest(w)
		return false
	}
	return true
}

// GenerateVerificationCodeHandler /verificationcode/generate
func GenerateVerificationCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}
	rand.Seed(time.Now().UnixNano())
	r.ParseForm()
	mailAddress := db.EscapeSinglequotation(r.FormValue("mail_address"))
	verificationCode := rand.Intn(999999)
	expiration := db.TimeNow(15)

	req := &GenerateVerificationCodeRequest{
		MailAddress: mailAddress,
	}
	validate := validator.New()
	if err := validate.Struct(req); checkerr.BadRequest(&w, err) {
		return
	}

	if findMailAddressDuplicate(&w, mailAddress) {
		return
	}

	if _, err := db.RunSQL(fmt.Sprintf(
		"DELETE FROM pending WHERE mail_address = '%s'",
		mailAddress,
	)); checkerr.InternalServerError(&w, err) {
		return
	}

	if _, err := db.RunSQL(fmt.Sprintf(
		"INSERT INTO pending (mail_address, verification_code, expiration) VALUES('%s', %d, '%s')",
		mailAddress, verificationCode, expiration,
	)); checkerr.InternalServerError(&w, err) {
		return
	}

	msg := fmt.Sprintf(
		"アカウントを作成するには15分以内に確認コードを入力してください \r\n <h1>%d</h1> \r\n",
		verificationCode,
	)
	m := mail.Mail{
		From:    "noreply@wazuplay.online",
		To:      mailAddress,
		Subject: fmt.Sprintf("確認コード: %d", verificationCode),
		Text:    msg,
	}
	if err := mail.Send(m); checkerr.InternalServerError(&w, err) {
		return
	}
	httpstates.OK(&w)
}

// ConfirmVerificationCodeHandler /verificationcode/confirm
func ConfirmVerificationCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}
	r.ParseForm()
	mailAddress := db.EscapeSinglequotation(r.FormValue("mail_address"))
	verificationCode, err := strconv.Atoi(r.FormValue("verification_code"))
	if checkerr.InternalServerError(&w, err) {
		return
	}

	req := &ConfirmVerificationCodeRequest{
		MailAddress:      mailAddress,
		VerificationCode: verificationCode,
	}
	validate := validator.New()
	if err := validate.Struct(req); checkerr.BadRequest(&w, err) {
		return
	}

	if !findVerificationCode(&w, mailAddress, verificationCode) {
		return
	}

	httpstates.OK(&w)
}

// SignUpHandler /signup
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpstates.MethodNotAllowed(&w)
		return
	}
	r.ParseForm()
	mailAddress := db.EscapeSinglequotation(r.FormValue("mail_address"))
	verificationCode, err := strconv.Atoi(r.FormValue("verification_code"))
	if checkerr.InternalServerError(&w, err) {
		return
	}
	username := db.EscapeSinglequotation(r.FormValue("username"))
	password := db.EscapeSinglequotation(r.FormValue("password"))

	req := &SignUpRequest{
		MailAddress:      mailAddress,
		VerificationCode: verificationCode,
		Username:         username,
		Password:         password,
	}
	validate := validator.New()
	if err := validate.Struct(req); checkerr.BadRequest(&w, err) {
		return
	}

	if !findVerificationCode(&w, mailAddress, verificationCode) {
		return
	}

	if findMailAddressDuplicate(&w, mailAddress) {
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if checkerr.InternalServerError(&w, err) {
		return
	}
	id, err := random.GenerateRandomString()
	if checkerr.InternalServerError(&w, err) {
		return
	}

	if _, err := db.RunSQL(fmt.Sprintf(
		"INSERT INTO users (user_id, username, mail_address, hashed_password) VALUES('%s', '%s', '%s', '%s')",
		id, username, mailAddress, hashedPassword,
	)); err != nil {
		httpstates.InternalServerError(&w)
		return
	}

	httpstates.OK(&w)
}
