package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"strconv"
	"time"
)

// Visit model
type Visit struct {
	ID			int64	`json:"id"`
	PlaceId		int64	`json:"place_id"`
	DeviceId	string	`json:"device_id"`
}

// regex expression to capture Visit in log lines
var logLineRegex = regexp.MustCompile(`(?i)Visit:[\s\,]+?id=(?P<id>\d+)[\s\,]+?device_id=(?P<device_id>\w+)[\s\,]+?place_id=(?P<place_id>\d+)`)


// Receives a log line and returns a Visit
func ParseVisitFromLogLine(line string) (visit *Visit, err error) {
	match := logLineRegex.FindStringSubmatch(line)
	matchedID, err := strconv.Atoi(match[1])
	matchedPlaceID, err := strconv.Atoi(match[3])
	if err != nil {
		return
	}

	visit = &Visit{
		ID: int64(matchedID),
		PlaceId: int64(matchedPlaceID),
		DeviceId: match[2],
	}

	return
}

// Search campaigns by visit/targeting
func (v *Visit) ListCampaigns(mongoCollection *mongo.Collection) (results []*Campaign, err error) {

	// Search campaigns which contains the visit's place into its targeting
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := mongoCollection.Find(ctx, bson.M{"targeting.place_id": v.PlaceId})
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

// Store visit into MongoDB to indicate that this ID has been processed
func (v *Visit) Store(mongoCollection *mongo.Collection) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = mongoCollection.InsertOne(ctx, v)

	return
}

// Check if MongoDB has stored any visit with same ID
func (v *Visit) HasBeenProcessed(mongoCollection *mongo.Collection) bool {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var result Visit
	if err := mongoCollection.FindOne(ctx, bson.M{"_id": v.ID}).Decode(&result); err != nil {
		return true
	}
	return false
}