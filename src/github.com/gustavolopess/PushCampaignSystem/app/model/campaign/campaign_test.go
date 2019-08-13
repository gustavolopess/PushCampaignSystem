package campaign

import (
	"context"
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/config"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"os"
	"reflect"
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

func TestCampaign_Store(t *testing.T) {
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
			t.Errorf("Campaign.Store(): Test validation failed")
		}
	}
}


func TestStoreMultiple(t *testing.T) {
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
			t.Errorf("StoreMultiple(): error = %v", err)
		} else if !tt.validate() {
			t.Errorf("StoredMultiple(): Test validation failed")
		}
	}
}

func TestLoadCampaigns(t *testing.T) {
	strCampaigns := []byte(`
[
  {
    "id": 876,
    "provider": "localytics",
    "push_message": "Magna ad aute ad aliqua excepteur excepteur quis esse ad dolor enim incididunt.",
    "targeting": [
      {
        "place_id": 4155,
        "description": "Krag"
      },
      {
        "place_id": 7722,
        "description": "Genesynk"
      },
      {
        "place_id": 4940,
        "description": "Frosnex"
      }
    ]
  },
  {
    "id": 7248,
    "provider": "localytics",
    "push_message": "Pariatur qui esse est aliqua dolore elit cillum fugiat anim eiusmod non enim tempor amet.",
    "targeting": [
      {
        "place_id": 4173,
        "description": "Centice"
      },
      {
        "place_id": 7526,
        "description": "Rodemco"
      },
      {
        "place_id": 4222,
        "description": "Intergeek"
      }
    ]
  }]
`)

	expectedOutput := []Campaign{
		{
			876,
			"localytics",
			"Magna ad aute ad aliqua excepteur excepteur quis esse ad dolor enim incididunt.",
			[]Place{{4155, "Krag"}, {7722, "Genesynk"}, {4940, "Frosnex"}},
		},
		{
			7248,
			"localytics",
			"Pariatur qui esse est aliqua dolore elit cillum fugiat anim eiusmod non enim tempor amet.",
			[]Place{{4173, "Centice"}, {7526, "Rodemco"}, {4222, "Intergeek"}},
		},
	}

	validTempFile, err := ioutil.TempFile("", "testCampaigns_valid.json")
	if err != nil {
		t.Errorf("Could not create validTempFile with fictional campaigns: %v", err)
	}
	defer os.Remove(validTempFile.Name())


	invalidTempFile, err := ioutil.TempFile("","testCampaigns_invalid.json")
	if err != nil {
		t.Errorf("Could not create invalidTempFile with fictional campaigns: %v", err)
	}
	defer os.Remove(invalidTempFile.Name())


	if _, err = validTempFile.Write(strCampaigns); err != nil {
		t.Errorf("Could not write fictional campaigns to validTempFile %v", err)
	}

	if _, err = invalidTempFile.Write([]byte(`[{"id": 777}, {"push_message": "lhebs"}]`)); err != nil {
		t.Errorf("Could not write fictional campaigns to validTempFile %v", err)
	}

	tests := []struct {
		name string
		campaignFilePath string
		expetedCampaigns []Campaign
		wantError bool
	}{
		{
			"Decode a JSON with valid campaigns",
			validTempFile.Name(),
			expectedOutput,
			false,
		},
		{
			"Decode a JSON with invalid campaigns",
			invalidTempFile.Name(),
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		campaigns, err := LoadCampaigns(tt.campaignFilePath)
		if (err != nil) != tt.wantError {
			t.Errorf("LoadCampaigns(): undesired error behaviour (%v, wantError=%v, Decoded=%v)", err, tt.wantError, campaigns)
		} else if !tt.wantError && !reflect.DeepEqual(campaigns, tt.expetedCampaigns) {
			t.Errorf("LoadCampaigns(): decoded campaigns are different from expected (Decoded=%v, Expected=%v)", campaigns, tt.expetedCampaigns)
		}
	}

	if err = validTempFile.Close(); err != nil {
		t.Errorf("Could not close validTempFile %v", err)
	}
}