package model

// place model
type Place struct {
	PlaceId     int64  `json:"place_id" bson:"place_id"`
	Description string `json:"description" bson:"description"`
}