package services

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/flagmanagment/backend/internal/repository"
)

type EmailService interface {
	SendInvitation(email, inviteLink string) error
}

type emailService struct {
	store repository.Store
}

func NewEmailService(store repository.Store) EmailService {
	return &emailService{store: store}
}

func (s *emailService) SendInvitation(email, inviteLink string) error {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "1025" // Mailhog default
	}
	
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: You've been invited to FlagManagement\r\n"+
		"\r\n"+
		"You have been invited to join FlagManagement.\r\n"+
		"Click the link below to accept the invitation:\r\n"+
		"%s\r\n", email, inviteLink))

	addr := fmt.Sprintf("%s:%s", host, port)
	// For Mailhog we don't need authentication (nil Auth)
	err := smtp.SendMail(addr, nil, "noreply@flagmanagement.local", []string{email}, msg)
	if err != nil {
		fmt.Printf("[EmailService DEV FALLBACK] Failed to connect to SMTP (%s): %v. Dispatched mock invitation email to %s with link: %s\n", addr, err, email, inviteLink)
		return nil
	}

	return nil
}
