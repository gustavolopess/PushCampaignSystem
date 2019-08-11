package model

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

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

// read config from JSON file at specified path
func (m *MongoConn) LoadConfig(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read MongoDB config file: %s", err.Error())
	}
	err = json.Unmarshal([]byte(file), m)
	if err != nil {
		log.Fatalf("Invalid Mongo config file: %s", err.Error())
	}
}

// Instance a new connection with MongoDB
func (m *MongoConn) Connect() {
	log.Println("Establishing new connection with MongoDB...")

	// Set client options
	clientOptions := options.Client().ApplyURI(m.Url)

	// Connect to MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err.Error())
	}

	// Check connection
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not discovery a MongoDB server: %s", err.Error())
	}

	// Instance ref to collection
	collection = client.Database(m.Database).Collection(m.Collection)

	log.Println("Connection with MongoDB established")
}

// Return object to MongoDB collection
func MongoCollection() *mongo.Collection {
	return collection
}