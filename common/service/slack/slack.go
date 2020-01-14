package slack

// Storage store function in slack package
type Storage interface {
	SendSlackNotification(webhookURL, msg string) error
}
