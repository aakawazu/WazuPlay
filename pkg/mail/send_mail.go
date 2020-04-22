package mail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

//Mail mail structure
type Mail struct {
	From    string
	To      string
	Subject string
	Text    string
}

//Send send mail
func Send(a Mail) error {
	m := gomail.NewMessage()
	m.SetHeader("From", a.From)
	m.SetHeader("To", a.To)
	m.SetHeader("Subject", a.Subject)
	m.SetBody("text/html", a.Text)

	d := gomail.NewDialer("mail", 25, "wazuplay", "test")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
