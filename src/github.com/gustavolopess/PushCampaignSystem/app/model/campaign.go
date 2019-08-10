package model

import (
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


func LoadCampaigns(filePath string) (campaigns []Campaign, err error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	err = json.Unmarshal(fileData, &campaigns)
	return
}


func (c *Campaign) Store() {

}