package campaign

import (
	"context"
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/config"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)


func init() {
	var mongoconfigtests string
	flag.StringVar(&mongoconfigtests, "mongoconfigtests", "",
		"File to MongoDB configurations which will be used on tests")
	flag.Parse()

	var mongoConn config.MongoConn
	mongoConn.LoadConfig(mongoconfigtests)

	mongoConn.Connect()

	_ = config.CampaignCollection().Drop(context.Background())
}

func findOneCampaign(filter bson.M) (campaign *Campaign) {
	err := config.CampaignCollection().FindOne(context.Background(), filter).Decode(&campaign)
	if err != nil {
		return nil
	}
	return
}

func findManyCampaigns(filter bson.M) (campaigns []*Campaign) {
	cur, err := config.CampaignCollection().Find(context.Background(), filter)
	if err != nil {
		return
	}

	for cur.Next(context.Background()) {
		var camp Campaign
		err = cur.Decode(&camp)
		if err != nil {
			return
		}

		campaigns = append(campaigns, &camp)
	}

	return
}

func TestCampaign_StoreOne(t *testing.T) {
	tests := []struct {
		name string
		campaign Campaign
		validate func() bool
	}{
		{
			"TestCampaign_StoreOne - Create event",
			Campaign{
				1,
				"localytics",
				"TestCampaign_StoreOnes",
				[]Place{
					{
						1,
						"Place1",
					},
					{
						2,
						"Place1",
					},
					{
						3,
						"Place1",
					},
				},
			},
			func() bool {
				campaign := findOneCampaign(bson.M{"_id": 1})
				return campaign != nil
			},
		},
		{
			"TestCampaign_StoreOne - Upsert event",
			Campaign{
				1,
				"localytics",
				"TestCampaign_StoreOne",
				[]Place{
					{
						1,
						"Place1",
					},
					{
						2,
						"Place1",
					},
					{
						3,
						"Place1",
					},
					{
						4,
						"Place1",
					},
				},
			},
			func() bool {
				campaign := findOneCampaign(bson.M{"_id": 1})
				if campaign == nil {
					return false
				}
				return len(campaign.Targeting) == 4
			},
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		err := tt.campaign.Store(config.CampaignCollection())
		if err != nil {
			t.Errorf("Campaign.Store() error = %v", err)
		} else if !tt.validate() {
			t.Errorf("Test validation failed")
		}
	}
}


func TestCampaign_StoreMultiple(t *testing.T) {
	campaign1 := Campaign{
		55,
		"localytics",
		"TestCampaign_StoreMany",
		[]Place{
			{
				1,
				"Place1",
			},
			{
				2,
				"Place1",
			},
			{
				3,
				"Place1",
			},
		},
	}

	campaign2 := Campaign{
		56,
		"localytics",
		"TestCampaign_StoreMany",
		[]Place{
			{
				1,
				"Place1",
			},
			{
				2,
				"Place1",
			},
			{
				3,
				"Place1",
			},
		},
	}

	campaign3 := Campaign{
		57,
		"localytics",
		"TestCampaign_StoreMany",
		[]Place{
			{
				1,
				"Place1",
			},
			{
				2,
				"Place1",
			},
			{
				3,
				"Place1",
			},
		},
	}

	tests := []struct {
		name string
		campaigns []Campaign
		validate func() bool
	}{
		{
			"TestCampaign_StoreMultiple - Store many events",
			[]Campaign{campaign1, campaign2, campaign3},
			func() bool {
				camps := findManyCampaigns(bson.M{"push_message": "TestCampaign_StoreMany"})
				return len(camps) == 3
			},
		},
		{
			"TestCampaign_StoreMultiple - Store many events with one upsertion",
			[]Campaign{campaign1, campaign2, {
				campaign3.ID,
				campaign3.Provider,
				campaign3.PushMessage,
				append(campaign3.Targeting, Place{5, "Place 5"}),

			}},
			func() bool {
				camp := findOneCampaign(bson.M{"_id": campaign3.ID})
				return len(camp.Targeting) == len(campaign3.Targeting) + 1
			},
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		err := StoreMultiple(tt.campaigns, config.CampaignCollection())
		if err != nil {
			t.Errorf("StoreMultiple() error = %v", err)
		} else if !tt.validate() {
			t.Errorf("Test validation failed")
		}
	}
}