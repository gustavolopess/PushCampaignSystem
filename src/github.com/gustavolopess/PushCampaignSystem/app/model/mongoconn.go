package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

type MongoConn struct {
	Url string	`json:"url"`
	Database string `json:"database"`
	Collection string `json:"collection"`
}

// Instance a new connection with MongoDB
func (m *MongoConn) Connect() {
	// Set client options
	clientOptions := options.Client().ApplyURI(m.Url)

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err.Error())
	}

	// Check connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Could not maintain a conncetion with MongoDB: %s", err.Error())
	}

	// Instance ref to collection
	collection = client.Database(m.Database).Collection(m.Collection)
}

// Return object to MongoDB collection
func MongoCollection() *mongo.Collection {
	return collection
}