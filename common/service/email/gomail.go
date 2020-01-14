package email

import (
	"crypto/tls"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"

	"github.com/golang-microservices/template/model"
)

// EmailClient manage all email action
type EmailClient struct {
	dialer  *gomail.Dialer
	conf    *model.Email
	session gomail.SendCloser
}

/*
	@sessionMappingGomail: Mapping between EmailConfig model (hashed) and EmailClient for singleton pattern
*/
var (
	sessionMappingGomail = make(map[string]*EmailClient)
)

// NewEmailClient function return gomain client based on singleton pattern
func NewEmailClient(config *model.Service) Mailstore {
	hash := config.Hash()
	currentSession := sessionMappingGomail[hash]
	if currentSession == nil {
		currentSession = &EmailClient{nil, nil, nil}

		dialer, err := GetDialer(&config.Email)
		if err != nil {
			log.Println("Error when try to make strconv port from config: ", err)
			panic(err)
		}

		session, err := dialer.Dial()
		if err != nil {
			log.Println("Error when try to dial to mail server: ", err)
			panic(err)
		}
		log.Println("Connected to Mail Server")

		currentSession.dialer = dialer
		currentSession.conf = &config.Email
		currentSession.session = session
		sessionMappingGomail[hash] = currentSession
	}

	return currentSession
}

// GetDialer function return a new Dialer
func GetDialer(conf *model.Email) (client *gomail.Dialer, err error) {
	port, err := strconv.Atoi(conf.Port)
	if err != nil {
		return nil, err
	}
	client = gomail.NewDialer(conf.Host, port, conf.Username, conf.Password)
	client.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: conf.Host}
	client.LocalName = conf.Host

	return client, nil
}

// Send function sent a message via email
func (e *EmailClient) Send(from, recipient, subject, msg string) (err error) {
	return e.mailSender(from, recipient, subject, msg)
}

// mailSender private function sent mail based on argument provide
func (e *EmailClient) mailSender(from, to, subject, message string) (err error) {
	msg := gomail.NewMessage()

	msg.SetHeader("From", from)
	msg.SetAddressHeader("To", to, "")
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)

	if e.session != nil {
		if err = e.session.Send(from, []string{to}, msg); err != nil {
			// close current session
			e.session.Close()

			// generate new client
			e.dialer, _ = GetDialer(e.conf)

			// reconnect smtp server
			if newSession, err := e.dialer.Dial(); err == nil {
				e.session = newSession
				// resend email
				return e.session.Send(from, []string{to}, msg)
			}
			return err
		}
		return nil
	}

	return e.dialer.DialAndSend(msg)
}
