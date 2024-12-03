package email

import (
	"time"

	"gopkg.in/mail.v2"
)

type EmailSender interface {
	SendWelcomeEmail(to, username string) error
}

type EmailService struct {
	dialer    *mail.Dialer
	from      string
	templates *EmailTemplate
	baseURL   string
}

func NewEmailService(host string, port int, username, password, baseURL string) (*EmailService, error) {
	templates, err := NewEmailTemplate()
	if err != nil {
		return nil, err
	}

	return &EmailService{
		dialer:    mail.NewDialer(host, port, username, password),
		from:      username,
		templates: templates,
		baseURL:   baseURL,
	}, nil
}

func (s *EmailService) SendWelcomeEmail(to, username string) error {
	data := WelcomeEmailData{
		Username: username,
		LoginURL: s.baseURL + "/login",
		Year:     time.Now().Year(),
	}

	body, err := s.templates.ExecuteTemplate("welcome.html", data)
	if err != nil {
		return err
	}

	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Task Manager!")
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}
