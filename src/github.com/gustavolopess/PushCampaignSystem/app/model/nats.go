package model

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/nats-io/go-nats-streaming"
	"log"
)

// struct with NATS configurations
// some fields are used by both publisher and subscriber (e.g.: host, port, cluster_id and client_id)
// other fields are used only by subscriber (e.g.: durable_name and subject)
type NatsConn struct {
	Host		string	`json:"host"`
	Port		int		`json:"port"`
	ClusterID	string	`json:"cluster_id"`
	ClientID	string	`json:"client_id"`
	Subject		string	`json:"subject"`
	DurableName string	`json:"durable_name"` // only used by subscriber
}

var (
	pubQueue stan.Conn
	subQueue stan.Conn

	PUB_QUEUE = "PUB_QUEUE"
	SUB_QUEUE = "SUB_QUEUE"
)

// read config from JSON file at specified path
func (c *NatsConn) LoadConfig(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read NATS config file: %s", err.Error())
	}
	err = json.Unmarshal([]byte(file), c)
	if err != nil {
		log.Fatalf("Invalid NATS config file: %s", err.Error())
	}
}


func (c *NatsConn) Connect(queueType string) {
	url := fmt.Sprintf("nats://%s:%d", c.Host, c.Port)
	queue, err := stan.Connect(c.ClusterID, c.ClientID, stan.NatsURL(url))

	if err != nil {
		log.Fatalf("Could not connect to NATS streaming server: %s", err.Error())
	}

	switch queueType {
	case PUB_QUEUE:
		pubQueue = queue
		break
	case SUB_QUEUE:
		subQueue = queue
		break
	default:
		log.Fatalf("Invalid queue type")
	}

}

// Publish push notification message to stream queue
func (c *NatsConn) Publish(msg []byte) error {
	return pubQueue.Publish(c.Subject, msg)
}


// returns publisher connection object
func PublisherConnection() stan.Conn {
	return pubQueue
}


// returns subscriber connection object
func SubscriberConnection() stan.Conn {
	return subQueue
}
