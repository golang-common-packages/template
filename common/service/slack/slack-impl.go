package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-microservices/template/model"
)

// SlackClient manage all slack action
type SlackClient struct{}

// SendSlackNotification will post to an 'Incoming Webook' url setup in Slack Apps. It accepts
// some text and the slack channel is saved within Slack.
func (s *SlackClient) SendSlackNotification(webhookURL, msg string) error {
	buf := new(bytes.Buffer)
	client := &http.Client{Timeout: 10 * time.Second}

	slackBody, _ := json.Marshal(model.Message{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}

	return nil
}
