package campaign

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"time"
)

// place model
type Place struct {
	PlaceId     int64  `json:"place_id" bson:"place_id"`
	Description string `json:"description" bson:"description"`
}

// model campaign
type Campaign struct {
	ID          int64  	`json:"id" bson:"_id"`
	Provider    string 	`json:"provider" bson:"provider"`
	PushMessage string 	`json:"push_message" bson:"push_message"`
	Targeting	[]Place	`json:"targeting" bson:"targeting"`
}


// Returns an array of campaigns from input file
func LoadCampaigns(filePath string) (campaigns []Campaign, err error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(fileData, &campaigns)
	if err != nil {
		return
	}

	// Check if all campaigns are valid
	for _, campaign := range campaigns {
		if err = campaign.Valid(); err != nil {
			return
		}
	}

	return
}


// Check whether required fields are empty
func (c *Campaign) Valid() error {
	if c == nil {
		return fmt.Errorf("nil pointer campaign")
	}

	if c.ID ==  0 {
		return fmt.Errorf("campaign without ID")
	}

	if c.Provider == "" {
		return fmt.Errorf("campaign without provider")
	}

	if c.PushMessage == "" {
		return fmt.Errorf("campaign without push_message")
	}

	return nil
}

// Store campaign on database
func (c *Campaign) Store(mongoCollection *mongo.Collection) (err error) {
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
func StoreMultiple(campaigns []Campaign, mongoCollection *mongo.Collection) (err error) {
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
