package controller

import (
	"github.com/gustavolopess/PushCampaignSystem/app/model"
	"log"
)

// Read campaigns in file and store them into MongoDB
func StoreCampaignsFromFile(filePath string)  {

	// Load campaigns from file
	campaigns, err := model.LoadCampaigns(filePath)
	if err != nil {
		log.Fatalf("Campaigns could not be loaded from file: %s", err.Error())
	}

	// Store all them into database
	err = model.StoreMultipleCampaigns(campaigns)
	if err != nil {
		log.Fatalf("Could not inserto campaigns into MongoDB: %s", err.Error())
	}
}

