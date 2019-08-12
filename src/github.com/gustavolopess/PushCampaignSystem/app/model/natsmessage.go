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
	VisitId			int64  	`json:"visit_id"`
	Provider 		string	`json:"provider,omitempty"`
	PushMessage 	string	`json:"push_message,omitempty"`
	DeviceId		string 	`json:"device_id,omitempty"`
	HasCampaign		bool	`json:"has_campaign"`
}

var withCampaignOutputTemplate =
`=> Push sent regarding visit %d
===> Device ID: "%s"
===> %s logging: { "message": "%s", device_id: "%s" }
`

var withoutCampaignOutputTemplate =
`=> Push sent regarding visit %d
===> No campaign with matching target
`

func (n *NatsMessage) LoadMessage(data []byte) (err error) {
	err = json.Unmarshal(data, &n)
	return
}


func OnMessage(data []byte) {
	// Unmarshal message
	var natsMessage NatsMessage
	if err := json.Unmarshal(data, &natsMessage); err != nil {
		log.Fatalf("Invalid NATS message %s: %s", data, err.Error())
		return
	}

	// Subscriber output
	var output string

	// Check if message a registered campaign
	if natsMessage.HasCampaign {
		// Instantiate provider from factory
		provider, err := factory.GetProvider(natsMessage.Provider)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}


		// Send campaign's push notification
		err = provider.SendPushNotification(natsMessage.PushMessage, natsMessage.DeviceId)
		if err != nil {
			log.Println("Could not send push notification: %s", err.Error())
			return
		}

		// Prints the "with campaign" output format
		output = fmt.Sprintf(withCampaignOutputTemplate, natsMessage.VisitId, natsMessage.DeviceId,
				strings.Title(natsMessage.Provider), natsMessage.PushMessage, natsMessage.DeviceId)

	} else {
		// Prints the "without campaign" output format
		output = fmt.Sprintf(withoutCampaignOutputTemplate, natsMessage.VisitId)
	}

	fmt.Println(output)
}