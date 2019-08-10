package model

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"time"
)

// model campaign
type Campaign struct {
	ID          int64  	`json:"id" bson:"_id"`
	Provider    string 	`json:"provider"`
	PushMessage string 	`json:"push_message"`
	Targeting	[]Place	`json:"targeting"`
}


// Returns an array of campaigns from input file
func LoadCampaigns(filePath string) (campaigns []Campaign, err error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(fileData, &campaigns)
	return
}

// Store campaign on database
func (c *Campaign) StoreOne(mongoCollection *mongo.Collection) (err error) {
	// Uses upsert to update if it exists and create if doesn't and avoid existent key error
	updateOptions := &options.UpdateOptions{}
	updateOptions = updateOptions.SetUpsert(true)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = mongoCollection.UpdateOne(
		ctx,
		bson.M{"_id": c.ID},
		bson.D{
			{"$set", bson.D{
				{"provider", c.Provider},
				{"push_message", c.PushMessage},
				{"targeting", c.Targeting},
			}},
	}, updateOptions)
	return
}

// Store multiple campaigns
func StoreMultipleCampaigns(campaigns []Campaign, mongoCollection *mongo.Collection) (err error) {
	var operations []mongo.WriteModel

	// Uses upsert to update if it exists and create if doesn't and avoid existent key error
	for _, c := range campaigns {
		operation := mongo.NewUpdateOneModel()
		operation.SetFilter(bson.M{"_id": c.ID})
		operation.SetUpdate(bson.D{
			{"$set", bson.D{
				{"provider", c.Provider},
				{"push_message", c.PushMessage},
				{"targeting", c.Targeting},
			}},
		})
		operation.SetUpsert(true)

		operations = append(operations, operation)
	}

	// Bulk write to mitigate lacks of performance
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = mongoCollection.BulkWrite(ctx, operations)

	return
}


// Search campaigns by visit/targeting
func SearchCampaignsByVisit(visit Visit) (results []*Campaign, err error) {

	// Search campaigns which contains the visit's place into its targeting
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := MongoCollection().Find(ctx, bson.M{"targeting.place_id": visit.PlaceId})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	// Decode results into array of campaigns
	for cursor.Next(ctx) {
		// object to receive decoded document
		var campaign Campaign
		err = cursor.Decode(&campaign)
		if err != nil {
			return
		}

		results = append(results, &campaign)
	}
	err = cursor.Err()

	return
}
