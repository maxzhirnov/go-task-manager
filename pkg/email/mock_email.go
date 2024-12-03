package email

import "log"

type MockEmailService struct{}

// Ensure MockEmailService implements EmailSender
var _ EmailSender = (*MockEmailService)(nil)

func NewMockEmailService() EmailSender {
	return &MockEmailService{}
}

func (s *MockEmailService) SendWelcomeEmail(to, username string) error {
	log.Printf("Mock: Sending welcome email to %s (%s)", username, to)
	return nil
}
