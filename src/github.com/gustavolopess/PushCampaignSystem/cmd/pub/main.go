package main

import (
	"flag"
	"os"
	"github.com/gustavolopess/PushCampaignSystem/config"
	"os/signal"
)

func main() {
	// command line flags
	natsConfigPath := flag.String("natsconfig", "etc/queue", "Path to file with NATS configuration")
	flag.Parse()

	// init NATS configuration instance
	var natsConfig config.NatsConfig
	natsConfig.LoadConfig(*natsConfigPath)


	// subscribe to SIGINT signals
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt)


	// wait for SIGINT
	<- interruptChan
}
