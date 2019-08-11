package model

import (
	"encoding/json"
	"fmt"
	"github.com/gustavolopess/PushCampaignSystem/app/providers/factory"
	"log"
	"strings"
)

// Nats message model
type NatsMessage struct {
	VisitId			int64  `json:"visit_id"`
	Provider 		string	`json:"provider"`
	PushMessage 	string	`json:"push_message"`
	DeviceId		string 	`json:"device_id"`
}

var pushMessageTemplate = `
=> Push sent regarding visit %d
===> Device ID: "%s"
===> %s logging: { "message": "%s", device_id: "%s" }
`

func (n *NatsMessage) LoadMessage(data []byte) (err error) {
	err = json.Unmarshal(data, &n)
	return
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

func OnMessage(data []byte) {
	// Unmarshal message
	var natsMessage NatsMessage
	if err := json.Unmarshal(data, &natsMessage); err != nil {
		log.Fatalf("Invalid NATS message %s: %s", data, err.Error())
		return
	}


	// Instantiate provider from factory
	provider, err := factory.GetProvider(natsMessage.Provider)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// Crete push notification message and send it
	pushMessage := fmt.Sprintf(
		pushMessageTemplate,
		natsMessage.VisitId,
		natsMessage.DeviceId,
		strings.Title(natsMessage.Provider),
		natsMessage.PushMessage,
		natsMessage.DeviceId)
	provider.SendPushNotification(pushMessage)
}