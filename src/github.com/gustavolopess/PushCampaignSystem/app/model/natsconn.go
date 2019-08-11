package model

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats-streaming"
	"io/ioutil"
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

// publisher/subscriber connection
var (
	queue stan.Conn
)

// constants to indicate which kind of connection must be initiated
const (
	PubQueue = "pub-queue"
	SubQeueue = "sub-queue"
)

// read config from JSON file at specified path
func (c *NatsConn) LoadConfig(configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read NATS config file: %s", err.Error())
	}
	err = json.Unmarshal([]byte(file), &c)
	if err != nil {
		log.Fatalf("Invalid NATS config file: %s", err.Error())
	}
}


func (c *NatsConn) Connect(queueType string) {
	log.Println("Establishing new connection with NATS streaming server...")

	url := fmt.Sprintf("nats://%s:%d", c.Host, c.Port)

	// Avoid pub and sub use same clientID
	clientId := fmt.Sprintf("%s-%s", c.ClientID, queueType)

	var err error
	queue, err = stan.Connect(c.ClusterID, clientId, stan.NatsURL(url))

	if err != nil {
		log.Fatalf("Could not connect to NATS streaming server: %s", err.Error())
	}

	log.Println("Connection with NATS established")
}


func (c *NatsConn) Subscribe(onMessage func(data []byte)) {
	_, err := queue.Subscribe(c.Subject, func(m *stan.Msg){
		onMessage(m.Data)
	}, stan.DurableName(c.DurableName))

	if err != nil {
		log.Fatalf("Could not subscribe push notification messages from NATS: %s", err.Error())
	}
}

// Publish push notification message to stream queue
func (c *NatsConn) Publish(natsMessage *NatsMessage) (err error) {
	// Marshal NATS message
	msg, err := json.Marshal(natsMessage)
	if err != nil {
		return
	}

	// Send message to NATS
	err = queue.Publish(c.Subject, msg)

	return
}


// returns connection object
func Connection() stan.Conn {
	return queue
}
