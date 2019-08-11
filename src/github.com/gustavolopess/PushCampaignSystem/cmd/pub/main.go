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

	// init MongoDB
	var mongoConn model.MongoConn
	mongoConn.LoadConfig(*mongoConfigPath)
	mongoConn.Connect()

	logFileTail := controller.TailVisitLogFile(*visitLogPath)
	for line := range logFileTail {
		go func() {
			visit, campaigns := controller.SearchCampaignsByLogLine(line.Text)
			for _, c := range campaigns {
				natsMessage := &model.NatsMessage{
					Provider: c.Provider,
					Message:  c.PushMessage,
					DeviceId: visit.DeviceId,
				}
				go controller.EnqueueMessageIntoNats(natsMessage, &natsConn)
			}
		}()
	}


	// Subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// Wait for SIGINT
	<- interruptChan
}
