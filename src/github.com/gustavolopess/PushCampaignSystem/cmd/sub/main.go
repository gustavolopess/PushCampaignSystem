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
	flag.Parse()

	// Init NATS configuration instance and connect to streaming server
	var natsConn model.NatsConn
	natsConn.LoadConfig(*natsConfigPath)
	natsConn.Connect(model.SubQeueue)

	// Subscribe to NATS messages
	controller.DequeueMessagesFromNats(&natsConn)

	// Subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// Wait for SIGINT
	<- interruptChan
}