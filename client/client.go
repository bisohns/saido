package client

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	Socket          *websocket.Conn
	Send            chan *SendMessage
	Received        chan *ReceiveMessage
	StopHostPolling chan bool
}

// Write to websocket
func (client *Client) Write() {
	defer client.Socket.Close()
	var err error
	for msg := range client.Send {
		err = client.Socket.WriteJSON(msg)
		if err != nil {
			log.Error("Error inside client write ", err)
			// most likely socket connection has been closed so
			// just return
			return
		}
	}
}

// Read from websocket
func (client *Client) Read() {
	defer client.Socket.Close()
	for {
		var message *ReceiveMessage
		err := client.Socket.ReadJSON(&message)
		if err != nil {
			log.Errorf("While reading from client: %s", err)
			client.StopHostPolling <- true
			return
		} else {
			client.Received <- message
		}

	}
}
