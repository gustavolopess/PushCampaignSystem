package visit

import (
	"context"
	"flag"
	"github.com/gustavolopess/PushCampaignSystem/app/model/campaign"
	"github.com/gustavolopess/PushCampaignSystem/config"
	"go.mongodb.org/mongo-driver/bson"
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
}

func clearCollections() {
	_ = config.CampaignCollection().Drop(context.Background())
	_ = config.VisitCollection().Drop(context.Background())
}

func findOneVisit(filter bson.M) (visit *Visit) {
	err := config.VisitCollection().FindOne(context.Background(), filter).Decode(&visit)
	if err != nil {
		return nil
	}
	return
}

//func findManyVisits(filter bson.M) (visits []*Visit) {
//	cur, err := config.CampaignCollection().Find(context.Background(), filter)
//	if err != nil {
//		return
//	}
//
//	for cur.Next(context.Background()) {
//		var visit Visit
//		err = cur.Decode(&camp)
//		if err != nil {
//			return
//		}
//
//		campaigns = append(campaigns, &camp)
//	}
//
//	return
//}

func TestParseVisitFromLogLine(t *testing.T) {
	clearCollections()

	tests := []struct {
		name string
		logLine string
		expectedOutput *Visit
		wantError bool
	}{
		{
			"Decode a valid visit log (separated by comma and/or spaces)",
			"Visit: id=1234 device_id=id3123l, place_id=13",
			&Visit{
				1234,
				13,
				"id3123l",
			},
			false,
		},
		{
			"Decode a valid visit log (separated by comma and spaces)",
			"Visit: id=1234, device_id=id3123l, place_id=13",
			&Visit{
				1234,
				13,
				"id3123l",
			},
			false,
		},
		{
			"Decode a invalid visit log (without ID)",
			"Visit: id=, device_id=id3123l, place_id=13",
			nil,
			true,
		},
		{
			"Decode a invalid visit log (without device_id)",
			"Visit: id=1234, device_id=, place_id=13",
			nil,
			true,
		},
		{
			"Decode a invalid visit log (without place_id)",
			"Visit: id=1234, device_id=id3123l, place_id=",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		visit, err := ParseVisitFromLogLine(tt.logLine)
		if (err != nil) != tt.wantError {
			t.Errorf("ParseVisitFromLogLine(): undesired error behaviour (%v, wantError=%v)", err, tt.wantError)
		} else if !tt.wantError && !reflect.DeepEqual(*visit, *tt.expectedOutput) {
			t.Errorf("ParseVisitFromLogLine(): decoded visit is different from expected (Decoded=%v, Expected=%v)", visit, tt.expectedOutput)
		}
	}
}

func TestVisit_Store(t *testing.T) {
	clearCollections()

	tests := []struct {
		name string
		visit *Visit
		wantError bool
	}{
		{
			"Store visit",
			&Visit{
				12345,
				54321,
				"dev321123",
			},
			false,
		},
		{
			"Store visit",
			&Visit{
				12345,
				54321,
				"dev321123",
			},
			true,
		},
		{
			"Store visit",
			&Visit{
				312,
				3321,
				"dev32132",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		err := tt.visit.Store(config.VisitCollection())
		if (err != nil) != tt.wantError {
			t.Errorf("Visit.Store(): undesired error behaviour (%v, wantError=%v)", err, tt.wantError)
		} else {
			storedVisit := findOneVisit(bson.M{"_id": tt.visit.ID})
			if !reflect.DeepEqual(*storedVisit, *tt.visit) {
				t.Errorf("Visit.Store(): stored visit is different from original (stored=%x, original=%x)",
					*storedVisit, *tt.visit)
			}
		}
	}
}


func TestVisit_HasBeenProcessed(t *testing.T) {
	clearCollections()

	tests := []struct {
		name string
		visit *Visit
		toStore bool
	}{
		{
			"HasBeenProcessed with processed visit",
			&Visit{
				1234,
				4321,
				"dev32113",
			},
			true,
		},
		{
			"HasBeenProcessed with unprocessed visit",
			&Visit{
				12345,
				54321,
				"dev2113",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		if tt.toStore {
			err := tt.visit.Store(config.VisitCollection())
			if err != nil {
				t.Errorf("HasBeenProcessed(): Could not store visit on MongoDB: %v", err)
				continue
			}
		}

		hasBeenProcessed := tt.visit.HasBeenProcessed(config.VisitCollection())
		if hasBeenProcessed != tt.toStore {
			t.Errorf("HasBeenProcessed(): return different from toStore (%v, toStore=%v, ID=%d)",
				hasBeenProcessed, tt.toStore, tt.visit.ID)
		}
	}
}

func TestVisit_ListCampaigns(t *testing.T) {
	clearCollections()

	campaigns := []campaign.Campaign{
		{
			876,
			"localytics",
			"Magna ad aute ad aliqua excepteur excepteur quis esse ad dolor enim incididunt.",
			[]campaign.Place{{7722, "Genesynk"}},
		},
		{
			7248,
			"localytics",
			"Pariatur qui esse est aliqua dolore elit cillum fugiat anim eiusmod non enim tempor amet.",
			[]campaign.Place{{4173, "Centice"}, {7526, "Rodemco"}},
		},
		{
			772,
			"localytics",
			"Lorem ipsum dolor sit amet consectetur adipiscing elit.",
			[]campaign.Place{{4173, "Centice"}, {7722, "Genesynk"}},
		},
	}

	err := campaign.StoreMultiple(campaigns, config.CampaignCollection())
	if err != nil {
		t.Errorf("ListCampaigns(): could not store campaigns before start tests")
	}

	tests := []struct{
		name string
		visit Visit
		numberOfCampaigns int
	}{
		{
			"ListCampaigns of visit with one campaign",
			Visit{
				34688,
				7526,
				"dev7526",
			},
			1,
		},
		{
			"ListCampaigns of visit with two campaigns #1",
			Visit{
				3468,
				4173,
				"dev4173",
			},
			2,
		},
		{
			"ListCampaigns of visit with two campaigns #2",
			Visit{
				3768,
				7722,
				"dev7722",
			},
			2,
		},
		{
			"ListCampaigns of visit with zero campaigns",
			Visit{
				33768,
				12345,
				"dev12345",
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		campaigns, err := tt.visit.ListCampaigns(config.CampaignCollection())
		if err != nil {
			t.Errorf("ListCampaigns(): Could not list campaigns: %x", err)
			continue
		}

		if len(campaigns) != tt.numberOfCampaigns {
			t.Errorf("ListCampaigns(): numberOfCampaigns different from number of campaigns listed (%d != %d)",
				tt.numberOfCampaigns, len(campaigns))
		}
	}
}

