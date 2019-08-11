package model

import (
	"encoding/json"
)

// Nats message model
type NatsMessage struct {
	Provider 	string	`json:"string"`
	Message 	string	`json:"string"`
	DeviceId	string 	`json:"string"`
}


func (n *NatsMessage) EnqueueIntoNats(natsConn *NatsConn) (message []byte, err error) {
	// Struct to JSON
	message, err = json.Marshal(n)
	if err != nil {
		return
	}

	// Enqueue JSON into Nats streaming
	err = natsConn.Publish(message)
	return
}