package email

import (
	"github.com/golang-common-packages/template/model"
)

// Mailstore store function in email package
type Mailstore interface {
	Send(from, recipient, subject, msg string) error
}

const (
	GOMAIL = iota
	SENDGRID
)

// NewMailClient function for Factory Pattern
func NewMailClient(mailClientType int, config *model.Service) Mailstore {

	switch mailClientType {
	case GOMAIL:
		return NewEmailClient(config)
	case SENDGRID:
		return NewSendGridClient(config)
	}

	return nil
}
