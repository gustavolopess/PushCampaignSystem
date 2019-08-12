package controller

import (
	"github.com/gustavolopess/PushCampaignSystem/app/model/natsmessage"
	visit2 "github.com/gustavolopess/PushCampaignSystem/app/model/visit"
	"github.com/gustavolopess/PushCampaignSystem/config"
	"github.com/hpcloud/tail"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Orchestrate visits computation
func ProcessVisitsFromLog(visitsLogPath string, natsConn *config.NatsConn, campaignMongoCollection, visitMongoCollection *mongo.Collection) {

	// Tail visits log file
	logFileTail := TailVisitLogFile(visitsLogPath)
	for line := range logFileTail {

		go func() {
			// Parse line to visit
			visit, err := visit2.ParseVisitFromLogLine(line.Text)
			if err != nil {
				log.Fatalf("Could not parse line '%s': %s", line, err.Error())
			}

			// If this visit has been processed in the past, ignore it
			if visit.HasBeenProcessed(visitMongoCollection) {
				log.Printf("Visit with ID %d has been processed in the past. Ignoring.", visit.ID)
				return
			}

			// List campaigns with targeting containing visit's PlaceId
			campaigns, err := visit.ListCampaigns(campaignMongoCollection)
			if err != nil {
				log.Printf("Could not search campaign with visit %#v: %s. Ignoring", visit, err.Error())
				return
			}

			if len(campaigns) > 0 {
				// If it has campaigns, send messages to NATS with flag HasCampaign set to true
				for _, c := range campaigns {
					natsMessage := &natsmessage.NatsMessage{
						VisitId: visit.ID,
						Provider: c.Provider,
						PushMessage:  c.PushMessage,
						DeviceId: visit.DeviceId,
						HasCampaign: true,
					}
					go EnqueueMessageIntoNats(natsConn, natsMessage)
				}
			} else {
				// If it hasn't campaigns, send messages to NATS with flag HasCampaign set to false
				natsMessage := &natsmessage.NatsMessage{
					VisitId: visit.ID,
					HasCampaign: false,
				}
				go EnqueueMessageIntoNats(natsConn, natsMessage)
			}
		}()

	}
}

// Return generator to lines in log file
func TailVisitLogFile(filePath string) chan *tail.Line {
	t, err := tail.TailFile(filePath, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Could not tail visit log file: %s", err.Error())
	}

	return t.Lines
}

// Insert message into NATS pub queue
func EnqueueMessageIntoNats(natsConn *config.NatsConn, natsMessage *natsmessage.NatsMessage) {
	err := natsConn.Publish(natsMessage)
	if err != nil {
		log.Fatalf("Could not enqueue message %#v into NATS streaming: %s", natsMessage, err.Error())
	}

	log.Printf("Message %#v successfuly enqueued into NATS streaming\n", natsMessage)
}


// Consume messages from NATS sub queue
func DequeueMessagesFromNats(natsConn *config.NatsConn) {
	natsConn.Subscribe(natsmessage.OnMessage)
}