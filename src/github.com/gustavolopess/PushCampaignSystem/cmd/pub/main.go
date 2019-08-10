package main

import (
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/app/model"
	"os"
	"os/signal"
)

func main() {
	// command line flags
	natsConfigPath := flag.String("natsconfig", "etc/queue.json", "Path to file with NATS configuration")
	flag.Parse()

	// init NATS configuration instance
	var natsConn model.NatsConn
	natsConn.LoadConfig(*natsConfigPath)


	// subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// wait for SIGINT
	<- interruptChan
}
