package model

import (
	"context"
	"encoding/json"
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
