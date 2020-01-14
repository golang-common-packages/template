package websocket

// Storage store function in websocket package
type Storage interface {
	Reader() (message string, err error)
	Sender(message string) (err error)
}
