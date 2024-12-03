package email

import (
	"fmt"

	"gopkg.in/mail.v2"
)

type EmailSender interface {
	SendWelcomeEmail(to, username string) error
}

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService(host string, port int, username, password string) EmailSender {
	dialer := mail.NewDialer(host, port, username, password)
	return &EmailService{
		dialer: dialer,
		from:   username,
	}
}

func (s *EmailService) SendWelcomeEmail(to, username string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Task Manager!")

	body := fmt.Sprintf(`
        <h2>Welcome to Task Manager, %s!</h2>
        <p>Thank you for registering with us. We're excited to have you on board!</p>
        <p>You can now start creating and managing your tasks.</p>
        <p>Best regards,<br>Task Manager Team</p>
    `, username)

	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
