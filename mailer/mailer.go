package mailer

import (
	"gopkg.in/gomail.v2"
)

type GomailMailer struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
}

// SendEmail sends an email using Gomail
func (m *GomailMailer) SendEmail(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Username)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", body)

	d := gomail.NewDialer(m.SMTPHost, m.SMTPPort, m.Username, m.Password)

	return d.DialAndSend(msg)
}
