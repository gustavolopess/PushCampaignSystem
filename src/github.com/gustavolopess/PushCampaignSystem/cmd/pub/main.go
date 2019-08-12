package main

import (
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/app/controller"
	"github.com/gustavolopess/PushCampaignSystem/config"
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
	var natsConn config.NatsConn
	natsConn.LoadConfig(*natsConfigPath)
	natsConn.Connect(config.PubQueue)

	// Init MongoDB
	var mongoConn config.MongoConn
	mongoConn.LoadConfig(*mongoConfigPath)
	mongoConn.Connect()

	// Handle visits
	controller.ProcessVisitsFromLog(*visitLogPath, &natsConn, config.CampaignMongoCollection(), config.VisitMongoCollection())

	// Subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// Wait for SIGINT
	<- interruptChan
}
