package natsmessage

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

// Check whether required fields are empty
func (n *NatsMessage) Validate() error {
	if n == nil {
		return fmt.Errorf("nil pointer message")
	}

	if n.VisitId ==  0 {
		return fmt.Errorf("message without visit_id")
	}

	if n.Provider == "" {
		return fmt.Errorf("message without provider")
	}

	if n.PushMessage == "" {
		return fmt.Errorf("message without push_message")
	}

	if n.DeviceId == "" {
		return fmt.Errorf("message without device_id")
	}

	return nil
}

// Decode JSON message into NatsMessage
func LoadMessage(data []byte) (natsMessage NatsMessage, err error) {
	err = json.Unmarshal(data, &natsMessage)
	if err != nil {
		return
	}

	return natsMessage, natsMessage.Validate()
}

func OnMessage(data []byte) {
	// Unmarshal message
	natsMessage, err := LoadMessage(data)
	if err != nil {
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
			log.Printf("Could not send push notification: %s", err.Error())
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