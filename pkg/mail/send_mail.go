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

	// servername := "mail:25"
	// host, _, _ := net.SplitHostPort(servername)

	// msg := "From: " + m.From + "\r\n" +
	// 	"To: " + m.To + "\r\n" +
	// 	"Subject: " + m.Subject + "\r\n\r\n" +
	// 	m.Text + "\r\n"

	// password := os.Getenv("SMTP_PASSWORD")
	// auth := smtp.PlainAuth("", "wazuplay", password, host)

	// client, err := smtp.Dial(servername)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err = client.Auth(auth); err != nil {
	// 	log.Panic(err)
	// }
	// client.Mail(m.From)
	// client.Rcpt(m.To)

	// body, err := client.Data()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer body.Close()

	// buf := bytes.NewBufferString(msg)
	// if _, err = buf.WriteTo(body); err != nil {
	// 	log.Fatal(err)
	// }

	// client.Quit()

	return nil
}
