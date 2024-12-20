// Package email provides email sending functionality for the application,
// including welcome and verification emails with HTML templates.
package email

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mail.v2"
)

// EmailSender defines the interface for sending various types of emails.
// This interface allows for easy mocking in tests and flexibility in
// implementation.
type EmailSender interface {
	// SendWelcomeEmail sends a welcome email to new users
	SendWelcomeEmail(to, username string) error

	// SendVerificationEmail sends an email verification link
	SendVerificationEmail(to, username, token string) error

	SendPasswordResetEmail(email, resetLins string) error
}

// EmailService implements the EmailSender interface and handles
// email sending operations using SMTP.
type EmailService struct {
	dialer    *mail.Dialer   // SMTP connection handler
	from      string         // Sender email address
	templates *EmailTemplate // Email template manager
	baseURL   string         // Base URL for email links
}

// NewEmailService creates a new email service instance with the provided configuration.
//
// Parameters:
//   - host: SMTP server hostname
//   - port: SMTP server port
//   - username: SMTP authentication username
//   - password: SMTP authentication password
//   - baseURL: Base URL for application links in emails
//
// Returns:
//   - *EmailService: Configured email service
//   - error: Any error during template initialization
func NewEmailService(host string, port int, username, from, fromName, password, baseURL string) (*EmailService, error) {
	templates, err := NewEmailTemplate()
	if err != nil {
		return nil, err
	}

	return &EmailService{
		dialer:    mail.NewDialer(host, port, username, password),
		from:      fmt.Sprintf("%s <%s>", fromName, from),
		templates: templates,
		baseURL:   baseURL,
	}, nil
}

// WelcomeEmailData contains the data needed for the welcome email template.
type WelcomeEmailData struct {
	Username string // User's display name
	LoginURL string // URL to the login page
	Year     int    // Current year for copyright notice
}

// SendWelcomeEmail sends a welcome email to a newly registered user.
//
// Parameters:
//   - to: Recipient email address
//   - username: Recipient's username
//
// Returns:
//   - error: Any error encountered during email sending
//
// Template Data:
//   - Username: User's display name
//   - LoginURL: URL to the login page
//   - Year: Current year for copyright
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	// Prepare template data
	data := WelcomeEmailData{
		Username: username,
		LoginURL: s.baseURL + "/login",
		Year:     time.Now().Year(),
	}

	// Execute email template
	body, err := s.templates.ExecuteTemplate("welcome.html", data)
	if err != nil {
		return err
	}

	// Create and send email
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Welcome to Task Manager!")
	m.SetBody("text/html", body)

	return s.dialer.DialAndSend(m)
}

// VerificationEmailData contains the data needed for the verification email template.
type VerificationEmailData struct {
	Username         string
	VerificationLink string
	Year             int
	CurrentTime      string
}

// SendVerificationEmail sends an email with a verification link to a user.
//
// Parameters:
//   - to: Recipient email address
//   - username: Recipient's username
//   - token: Verification token
//
// Returns:
//   - error: Any error encountered during email sending
//
// Template Data:
//   - Username: User's display name
//   - VerificationLink: Complete verification URL with token
//   - Year: Current year for copyright
func (s *EmailService) SendVerificationEmail(to, username, token string) error {
	// Create email message
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify Your Email Address")

	// Prepare template data
	data := VerificationEmailData{
		Username:         username,
		VerificationLink: fmt.Sprintf("%s/verify-email?token=%s", s.baseURL, token),
		Year:             time.Now().Year(),
		CurrentTime:      time.Now().Format(time.RFC3339),
	}

	// Execute template
	body, err := s.templates.ExecuteTemplate("verification.html", data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %v", err)
	}

	m.SetBody("text/html", body)
	return s.dialer.DialAndSend(m)
}

type PasswordResetEmailData struct {
	ResetLink string
	IPAddress string
	Year      int
	Timestamp string
}

func (s *EmailService) SendPasswordResetEmail(to, resetLink string) error {
	maskedTo := maskEmail(to)
	log.Printf("Sending password reset email to: %s", maskedTo)

	// Create email message
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	log.Printf("Setting header From: %s", s.from)

	m.SetHeader("To", to)
	log.Printf("Setting header To: %s", maskedTo)

	m.SetHeader("Subject", "Reset Your Password - ActionHub")
	log.Printf("Setting header Subject: Reset Your Password - ActionHub")

	// Get IP address from request context (you might need to pass this through)
	ipAddress := "Unknown"

	// Prepare template data
	data := PasswordResetEmailData{
		ResetLink: resetLink,
		IPAddress: ipAddress,
		Year:      time.Now().Year(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Execute template
	body, err := s.templates.ExecuteTemplate("password-reset.html", data)
	if err != nil {
		log.Printf("Failed to execute password reset email template: %v", err)
		return fmt.Errorf("failed to execute email template: %v", err)
	}
	log.Printf("Executed password reset email template")

	m.SetBody("text/html", body)
	// log.Printf("Set body text/html: %s", body)

	// Send email
	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send password reset email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}
	log.Printf("Successfully sent password reset email to: %s", maskedTo)
	return nil
}

// maskEmail masks part of the email for logging purposes
// Example: j***@example.com
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "invalid_email"
	}
	if len(parts[0]) <= 1 {
		return email
	}
	return parts[0][:1] + "***@" + parts[1]
}
