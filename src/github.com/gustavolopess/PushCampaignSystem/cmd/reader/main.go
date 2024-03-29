package main

import (
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/app/controller"
	"github.com/gustavolopess/PushCampaignSystem/config"
)

func main() {
	// Command line flags
	filePath := flag.String(
		"campaignfile",
		"input/activeCampaigns.json",
		"Path to active campaigns file")

	mongoConfigPath := flag.String(
		"mongoconfig",
		"etc/mongoConfig.json",
		"Path to file with MongoDB configuration")

	flag.Parse()

	// init MongoDB
	var mongoConn config.MongoConn
	mongoConn.LoadConfig(*mongoConfigPath)
	mongoConn.Connect()

	// Store campaigns from file
	controller.StoreCampaignsFromFile(*filePath, config.CampaignCollection())
}
