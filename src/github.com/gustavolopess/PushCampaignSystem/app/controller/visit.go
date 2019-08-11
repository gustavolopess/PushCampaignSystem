package controller

import (
	"github.com/gustavolopess/PushCampaignSystem/app/model"
	providerFactory "github.com/gustavolopess/PushCampaignSystem/app/providers/factory"
	"github.com/hpcloud/tail"
	"log"
)

// Return generator to lines in log file
func TailVisitLogFile(filePath string) chan *tail.Line {
	t, err := tail.TailFile(filePath, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Could not tail visit log file: %s", err.Error())
	}

	return t.Lines
}

// Search campaigns which contain same place_id as present in log line
func SearchCampaignsByLogLine(line string) (*model.Visit, []*model.Campaign) {
	// Parse line to visit
	visit, err := model.ParseVisitLogLine(line)
	if err != nil {
		log.Fatalf("Could not parse line '%s': %s", line, err.Error())
	}

	campaigns, err := model.SearchCampaignsByVisit(visit)
	if err != nil {
		log.Fatalf("Could not search campaign with visit %#v: %s", visit, err.Error())
	}

	return visit, campaigns
}

// Insert message into NATS pub queue
func EnqueueMessageIntoNats(natsMessage *model.NatsMessage, natsConn *model.NatsConn) {
	_, err := natsMessage.EnqueueIntoNats(natsConn)
	if err != nil {
		log.Fatalf("Could not enqueue message %#v into NATS streaming: %s", natsMessage, err.Error())
	}

	log.Printf("Message %#v successfuly enqueued into NATS streaming\n", natsMessage)
}

// Send push notification
func SendPushNotification(natsMessage *model.NatsMessage) {
	provider, err := providerFactory.GetProvider(natsMessage.Provider)

	if err != nil {
		log.Printf("Could not send push notification: %s", err.Error())
	} else {
		provider.SendPushNotification(natsMessage)
	}
}

