package main

import (
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/app/controller"
	"github.com/gustavolopess/PushCampaignSystem/app/model"
	"os"
	"os/signal"
)

func main() {
	// Command line flags
	natsConfigPath := flag.String(
		"natsconfig",
		"etc/queue.json",
		"Path to file with NATS configuration")

	mongoConfigPath := flag.String(
		"mongoconfig",
		"etc/mongoConfig.json",
		"Path to file with MongoDB configuration")

	visitLogPath := flag.String(
		"visitlogpath",
		"input/visit.log",
		"Path to log with visits which will be tailed")

	flag.Parse()

	// Init NATS configuration instance
	var natsConn model.NatsConn
	natsConn.LoadConfig(*natsConfigPath)
	natsConn.Connect(model.PubQueue)

	// Init MongoDB
	var mongoConn model.MongoConn
	mongoConn.LoadConfig(*mongoConfigPath)
	mongoConn.Connect()

	// Tail log file with visits
	logFileTail := controller.TailVisitLogFile(*visitLogPath)
	for line := range logFileTail {
		go func() {
			visit, campaigns := controller.SearchCampaignsByLogLine(line.Text)
			if len(campaigns) > 0 {
				for _, c := range campaigns {
					natsMessage := &model.NatsMessage{
						VisitId: visit.ID,
						Provider: c.Provider,
						PushMessage:  c.PushMessage,
						DeviceId: visit.DeviceId,
						HasCampaign: true,
					}
					go controller.EnqueueMessageIntoNats(&natsConn, natsMessage)
				}
			} else {
				natsMessage := &model.NatsMessage{
					VisitId: visit.ID,
					HasCampaign: false,
				}
				go controller.EnqueueMessageIntoNats(&natsConn, natsMessage)
			}
		}()
	}


	// Subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// Wait for SIGINT
	<- interruptChan
}
