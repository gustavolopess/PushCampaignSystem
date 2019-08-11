package model

import (
	"encoding/json"
)

// Nats message model
type NatsMessage struct {
	VisitId			int64  `json:"visit_id"`
	Provider 		string	`json:"provider"`
	PushMessage 	string	`json:"push_message"`
	DeviceId		string 	`json:"device_id"`
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