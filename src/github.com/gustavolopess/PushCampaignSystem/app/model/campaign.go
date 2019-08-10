package model

type Campaign struct {
	ID          int64  	`json:"id"`
	Provider    string 	`json:"provider"`
	PushMessage string 	`json:"push_message"`
	Targeting	[]Place	`json:"targeting"`
}
