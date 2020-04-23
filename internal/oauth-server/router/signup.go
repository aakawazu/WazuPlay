package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	db "github.com/aakawazu/WazuPlay/pkg/db"
	httpStates "github.com/aakawazu/WazuPlay/pkg/http_states"
	mail "github.com/aakawazu/WazuPlay/pkg/mail"
)

// ConfirmMail confirm mail address
func ConfirmMail(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		mailAddress := r.FormValue("mail")
		if len(mailAddress) == 0 {
			httpStates.BadRequest(&w)
		} else {
			rand.Seed(time.Now().UnixNano())
			verificationCode := rand.Intn(999999)
			expiration := (time.Now().Add(15 * time.Minute)).Round(time.Second)
			expirationS := fmt.Sprintf("%s", expiration)
			expirationS = expirationS[0 : len(expirationS)-4]
			sqlStatement := fmt.Sprintf("DELETE FROM pending WHERE mail_address = '%s'", mailAddress)
			_, err := db.RunSQL(sqlStatement)
			if err != nil {
				httpStates.InternalServerError(&w)
				return
			}
			sqlStatement = fmt.Sprintf(
				"INSERT INTO pending (mail_address, verification_code, expiration) VALUES('%s', %d, '%s')",
				mailAddress, verificationCode, expirationS,
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
				fmt.Println(err)
				httpStates.InternalServerError(&w)
				return
			}
		}
	} else {
		httpStates.MethodNotAllowed(&w)
		return
	}
}
