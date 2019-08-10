package model

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
)

// model campaign
type Campaign struct {
	ID          int64  	`json:"id"`
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
func (c *Campaign) Store() (err error) {
	_, err = MongoCollection().InsertOne(context.TODO(), *c)
	return
}


func SearchCampaignsByVisit(visit Visit) (results []*Campaign, err error) {

	cursor, err := MongoCollection().Find(context.TODO(), bson.M{"targeting.place_id": visit.PlaceId})
	if err != nil {
		return
	}


	for cursor.Next(context.TODO()) {
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
