package email

import "log"

type MockEmailService struct{}

var _ EmailSender = (*MockEmailService)(nil)

func NewMockEmailService() EmailSender {
	return &MockEmailService{}
}

func (s *MockEmailService) SendWelcomeEmail(to, username string) error {
	log.Printf("Mock: Sending welcome email to %s (%s)", username, to)
	return nil
}

func (s *MockEmailService) SendVerificationEmail(to, username, token string) error {
	log.Printf("Mock: Sending verification email to %s (%s)", username, to)
	return nil
}

func (s *MockEmailService) SendPasswordResetEmail(email, resetLins string) error {
	log.Printf("Mock: Sending password reser email")
	return nil
}
