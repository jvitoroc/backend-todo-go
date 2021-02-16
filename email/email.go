package email

import (
	"fmt"
	"net/smtp"

	"github.com/jvitoroc/todo-go/config"
)

type EmailService struct {
	SmtpUser string
	SmtpPass string
	SmtpHost string
	SmtpAddr string
}

func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		SmtpUser: cfg.Email.SmtpUser,
		SmtpPass: cfg.Email.SmtpPass,
		SmtpHost: cfg.Email.SmtpHost,
		SmtpAddr: cfg.Email.SmtpAddr,
	}
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.SmtpUser, es.SmtpPass, es.SmtpHost)
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	return smtp.SendMail(es.SmtpAddr, auth, es.SmtpUser, []string{to}, msg)
}
