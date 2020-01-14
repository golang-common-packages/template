package email

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/golang-microservices/template/model"
)

// SendGridClient manage all email action
type SendGridClient struct {
	client *sendgrid.Client
	conf   *model.Email
}

/*
	@sessionMappingSendGrid: Mapping between model.EmailConfig and SendGridClient for singleton pattern
*/
var (
	sessionMappingSendGrid = make(map[string]*SendGridClient)
)

// NewSendGridClient function return sendgrid client based on singleton pattern
func NewSendGridClient(config *model.Service) Mailstore {
	hash := config.Hash()
	currentSession := sessionMappingSendGrid[hash]
	if currentSession == nil {
		currentSession = &SendGridClient{nil, nil}

		client := sendgrid.NewSendClient(config.Email.Key)
		log.Println("Connected to SendGrid Server")

		currentSession.client = client
		currentSession.conf = &config.Email
		sessionMappingSendGrid[hash] = currentSession
	}

	return currentSession
}

// Send function sent mail based on argument provide
func (s *SendGridClient) Send(from, recipient, subject, msg string) error {
	From := mail.NewEmail("Backend-golang Admin", from)
	To := mail.NewEmail("Backend-golang User", recipient)
	plainTextContent := msg
	htmlContent := "<strong>" + msg + "</strong>"
	message := mail.NewSingleEmail(From, subject, To, plainTextContent, htmlContent)
	_, err := s.client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
