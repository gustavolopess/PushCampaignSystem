package controller

import (
	"github.com/gustavolopess/PushCampaignSystem/app/model/campaign"
	"github.com/gustavolopess/PushCampaignSystem/app/model/natsmessage"
	"github.com/gustavolopess/PushCampaignSystem/app/model/visit"
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

		// Parse line to visit
		v, err := visit.ParseVisitFromLogLine(line.Text)
		if err != nil {
			log.Fatalf("Could not parse line '%s': %s", line, err.Error())
		}

		if !v.HasBeenProcessed(visitMongoCollection) {
			// If this visit is new, handle it
			go func() {
				// Store visit in MongoDB to avoid future repetitions
				err := v.Store(visitMongoCollection)
				if err != nil {
					log.Printf("Could not store visit with ID %d: %s. Ignoring", v.ID, err.Error())
					return
				}

				// List campaigns with targeting containing visit's PlaceId
				campaigns, err := v.ListCampaigns(campaignMongoCollection)
				if err != nil {
					log.Printf("Could not search campaigns of visit with ID %d: %s. Ignoring", v.ID, err.Error())
					return
				}

				if len(campaigns) > 0 {
					// If it has campaigns, send messages to NATS with according with each campaign
					for _, c := range campaigns {
						EnqueueMessageIntoNats(v, c, natsConn)
					}
				} else {
					// If it hasn't campaigns, send messages to NATS with flag HasCampaign set to false
					EnqueueMessageIntoNats(v, nil, natsConn)
				}
			}()
		} else {
			// If this visit has been processed in the past, ignore it
			log.Printf("Visit with ID %d has been processed in the past. Ignoring.", v.ID)
		}

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
func EnqueueMessageIntoNats(v *visit.Visit, campaign *campaign.Campaign, natsConn *config.NatsConn) {
	var natsMessage *natsmessage.NatsMessage
	if campaign != nil {
		// If campaign exists, set flag HasCampaign to true and create message with data from campaign
		natsMessage = &natsmessage.NatsMessage{
			VisitId: v.ID,
			Provider: campaign.Provider,
			PushMessage:  campaign.PushMessage,
			DeviceId: v.DeviceId,
			HasCampaign: true,
		}
	} else {
		// If campaign doesn't exist, set flag HasCampaign to false
		natsMessage = &natsmessage.NatsMessage{
			VisitId: v.ID,
			HasCampaign: false,
		}
	}

	// Send created message to NATS
	err := natsConn.Publish(natsMessage)
	if err != nil {
		log.Fatalf("Could not enqueue message from visit with ID %d into NATS streaming: %s", v.ID, err.Error())
	}

	log.Printf("Message from visit with ID %d successfuly enqueued into NATS streaming\n", v.ID)
}


// Consume messages from NATS sub queue
func DequeueMessagesFromNats(natsConn *config.NatsConn) {
	natsConn.Subscribe(natsmessage.OnMessage)
}