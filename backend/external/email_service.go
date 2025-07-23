package external

import (
	"fmt"
	"log"
)

type EmailService struct {
	apiKey   string
	baseURL  string
	client   *HTTPClient
}

func NewEmailService(apiKey, baseURL string) *EmailService {
	return &EmailService{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  NewHTTPClient(baseURL, 30000),
	}
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	headers := map[string]string{
		"Authorization": "Bearer " + e.apiKey,
	}

	data := map[string]interface{}{
		"to":      to,
		"subject": subject,
		"body":    body,
	}

	_, err := e.client.Post("/send", data, headers)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}

func (e *EmailService) SendWelcomeEmail(userEmail, userName string) error {
	subject := "Welcome to Our Platform!"
	body := fmt.Sprintf("Hello %s,\n\nWelcome to our platform! We're excited to have you on board.\n\nBest regards,\nThe Team", userName)
	
	return e.SendEmail(userEmail, subject, body)
}