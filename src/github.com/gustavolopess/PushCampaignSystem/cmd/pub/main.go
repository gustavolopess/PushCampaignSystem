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

	// Handle visits
	controller.ProcessVisitsFromLog(*visitLogPath, &natsConn, model.CampaignMongoCollection(), model.VisitMongoCollection())

	// Subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// Wait for SIGINT
	<- interruptChan
}
