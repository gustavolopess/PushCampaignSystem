package controller

import (
	"github.com/gustavolopess/PushCampaignSystem/app/model/campaign"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

// Read campaigns in file and store them into MongoDB
func StoreCampaignsFromFile(filePath string, mongoCollection *mongo.Collection)  {

	// Load campaigns from file
	campaigns, err := campaign.LoadCampaigns(filePath)
	if err != nil {
		log.Fatalf("Campaigns could not be loaded from file: %s", err.Error())
	}

	// Store all them into database
	err = campaign.StoreMultiple(campaigns, mongoCollection)
	if err != nil {
		log.Fatalf("Could not inserto campaigns into MongoDB: %s", err.Error())
	}

	log.Println("Campaigns successfully stored into MongoDB")
}

