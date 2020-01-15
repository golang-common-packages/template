package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/golang-common-packages/template/model"
)

// SharepointClient manage all sharepoint action
type SharepointClient struct {
	websocket *websocket.Conn
	Config    *model.Websocket
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Connect function establish websocket endpoint
func Connect(config model.Websocket) *SharepointClient {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	if err = ws.WriteMessage(1, []byte("Welcome to backend-golang websocket!!!")); err != nil {
		log.Println(err)
	}

	log.Println("Websocket endpoint is ready")

	return &SharepointClient{ws, &config}
}

// Reader function return message from websocket service
func (s *SharepointClient) Reader() (message string, err error) {
	messageType, p, err := s.websocket.ReadMessage()
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(p), nil
}

// Sender function sent message to websocket service
func (s *SharepointClient) Sender(message string) (err error) {
	if err = s.websocket.WriteMessage(1, []byte("Hi Client!")); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
