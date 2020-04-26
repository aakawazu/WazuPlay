package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	db "github.com/aakawazu/WazuPlay/pkg/db"
	httpStates "github.com/aakawazu/WazuPlay/pkg/http_states"
	mail "github.com/aakawazu/WazuPlay/pkg/mail"
	"github.com/go-playground/validator/v10"
)

// GenerateVerificationCode generate verification
func GenerateVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		mailAddress := r.FormValue("mail_address")
		verificationCode := rand.Intn(999999)
		expiration := db.TimeNow(15)
		type Request struct {
			MailAddress string `validate:"required,email"`
		}
		req := &Request{
			MailAddress: mailAddress,
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			httpStates.BadRequest(&w)
			return
		}
		mailAddress = db.EscapeSinglequotation(mailAddress)
		rand.Seed(time.Now().UnixNano())
		sqlStatement := fmt.Sprintf(
			"DELETE FROM pending WHERE mail_address = '%s'",
			mailAddress,
		)
		_, err := db.RunSQL(sqlStatement)
		if err != nil {
			httpStates.InternalServerError(&w)
			return
		}
		sqlStatement = fmt.Sprintf(
			"INSERT INTO pending (mail_address, verification_code, expiration) VALUES('%s', %d, '%s')",
			mailAddress, verificationCode, expiration,
		)
		if _, err := db.RunSQL(sqlStatement); err != nil {
			httpStates.InternalServerError(&w)
			return
		}
		msg := fmt.Sprintf(
			"アカウントを作成するには確認コードを入力してください \r\n <h1>%d</h1>",
			verificationCode,
		)
		subject := fmt.Sprintf("確認コード: %d", verificationCode)
		m := mail.Mail{
			From:    "noreply@wazuplay.online",
			To:      mailAddress,
			Subject: subject,
			Text:    msg,
		}
		if err := mail.Send(m); err != nil {
			httpStates.InternalServerError(&w)
			return
		}
		httpStates.OK(&w)
	} else {
		httpStates.MethodNotAllowed(&w)
	}
}

// ConfirmVerificationCode Check the verification code
func ConfirmVerificationCode(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		mailAddress := r.FormValue("mail_address")
		verificationCode, err := strconv.Atoi(r.FormValue("verification_code"))
		if err != nil {
			httpStates.InternalServerError(&w)
			return
		}
		type Request struct {
			MailAddress      string `validate:"required,email"`
			VerificationCode int    `validate:"gte=0,lt=999999"`
		}
		req := &Request{
			MailAddress:      mailAddress,
			VerificationCode: verificationCode,
		}
		validate := validator.New()
		if err := validate.Struct(req); err != nil {
			httpStates.BadRequest(&w)
			return
		}
		mailAddress = db.EscapeSinglequotation(mailAddress)
		sqlStatement := fmt.Sprintf(
			"SELECT * FROM pending WHERE mail_address = '%s' and verification_code = %d and expiration > '%s'",
			mailAddress, verificationCode, db.TimeNow(0),
		)
		rows, err := db.RunSQL(sqlStatement)
		defer rows.Close()
		if err != nil {
			httpStates.InternalServerError(&w)
			return
		}
		if !rows.Next() {
			httpStates.BadRequest(&w)
			return
		}
		httpStates.OK(&w)
	} else {
		httpStates.MethodNotAllowed(&w)
	}
}
