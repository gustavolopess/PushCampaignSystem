# PushCampaignSystem 🛰🌎📲
Ad toy system which allows clients to impact a mobile device when it is at a specific place with a push notification

## Summary
* [Description](https://github.com/gustavolopess/PushCampaignSystem#-description)
* [Requirements](https://github.com/gustavolopess/PushCampaignSystem#-requirements)
* [Architecture](https://github.com/gustavolopess/PushCampaignSystem#-architecture)
	* [FAQ](https://github.com/gustavolopess/PushCampaignSystem#faq)
* [Providers](https://github.com/gustavolopess/PushCampaignSystem#-providers)
* [How to build](https://github.com/gustavolopess/PushCampaignSystem#-how-to-build)
* [How to run](https://github.com/gustavolopess/PushCampaignSystem#-how-to-run)
	* [Reader](https://github.com/gustavolopess/PushCampaignSystem#reader)
	* [Publisher](https://github.com/gustavolopess/PushCampaignSystem#publisher)
	* [Subscriber](https://github.com/gustavolopess/PushCampaignSystem#subscriber)
	* [MongoDB and NATS Streaming Server](https://github.com/gustavolopess/PushCampaignSystem#mongodb-and-nats-streaming-server)
* [Project Backlog/Roadmap](https://github.com/gustavolopess/PushCampaignSystem#-project-backlogroadmap)

## 📃 Description

Clients can hire ad campaigns based on geolocation. These campaigns are defined by an identifier,
a targeting, which is a list of places that must be visited by some device in order to get this device 
impacted by the campaign, a push message, which the message to be sent via push notification to devices visiting places listed
in targeting and a provider, which is the push notification provider.

The campaign struct is defined as below:
```go
type Campaign struct {
	ID              int64  	
	Provider        string 	
	PushMessage     string 	
	Targeting       []Place	
}
```

Where place contains the place's ID and description:
```go
type Place struct {
	PlaceId         int64  
	Description     string 
}
```

This system receives two inputs:
1. A file with __active campaigns__, similar to [this](https://github.com/gustavolopess/PushCampaignSystem/blob/develop/input/activeCampaigns.json).
2. A file with __log of visits__ from multiple devices, similar to [this](https://github.com/gustavolopess/PushCampaignSystem/blob/develop/input/visit.log). 

\* Each active campaign from first input is decoded into `Campaign struct` described above.

\* Each visit from second input is decoded into `Visit struct` defined below:
```go
type Visit struct {
    ID          int64
    PlaceId     int64
    DeviceId    string
}
```


This system provides an output following this format:
```log
=> Push sent regarding visit <VisitId>
===> Device ID: <DeviceID>
===> <Provider> logging: { "message": <PushMessage>, "device_id": <DeviceId> }
```

The overall idea is: if a place_id from some visit of visits log file is registered into some
campaign, this system must send a push notification to device in this visit (reported in `DeviceId` field).
The indicator that push notification were sent through campaign's provider is the output described above.

## ✅ Requirements

* It must be easy to add a new push notification provider, which could have its own
  delivery logic.
* A campaign has only one push message, but can have many places as targeting.
* A campaign should not be delivered to a device if its location does not match any
  campaign's targeting.
  
 ## 🏡 Architecture
 
![architecture](https://github.com/gustavolopess/PushCampaignSystem/blob/develop/assets/architecture.png "Architecture")

This system can be visualized as three loosely coupled services: __reader, publisher and subscriber__.

* The __reader__ is responsible for reading  the active campaigns and store them into MongoDB
* The __publisher__ reads each visit from visits log and verifies in MongoDB whether the visit's place_id is registered in 
any campaign. If it is registered, the publisher sends to [NATS streaming server](https://github.com/nats-io/nats-streaming-server)
a message with `visit_id`, `provider`, `push_message` and `device_id`. If it is not registered, the publisher
sends a message with `visit_id` and the flag `has_campaign` set to false to indicate that the `place_id` does not belong to
any campaign. After send visit to NATS, the publisher insert it on MongoDB to keep track of already processed visits. 
The Golang struct of the message that can be sent to NATS is defined below:
```go
type NatsMessage struct {
    VisitId             int64   `json:"visit_id"`
    Provider            string	`json:"provider"`
    PushMessage         string	`json:"push_message"`
    DeviceId            string 	`json:"device_id"`
    HasCampaign         bool    `json:"has_campaign"`
}
```
* The __subscriber__ consumes each message sent to NATS streaming server by publisher. If the flag `has_campaign` of the 
message is true, the subscriber forwards the message to its designated provider and prints an output
similar to the described in the output example above. If the flag `has_campaign` of the message is false,
the subscriber just prints an output similar to the example below:
```log
=> Push sent regarding visit <VisitId>
===> No campaign with matching target
```


#### FAQ
1. __What happens if a campaign arrives with same ID of some campaign already registered on MongoDB?__
    * The __reader__ only performs [*upserts*](https://docs.mongodb.com/manual/reference/method/db.collection.update/#mongodb30-upsert-id)
    operations to insert campaigns on database:
        * If a campaign's ID does not exist at database, a new one is created.
        * If a campaign's ID exists at database, the information of existing campaign is updated with 
        information of arrived campaign (like an update).

2. __What happens if a visit arrives with same ID of some previous visit?__
    * The __publisher__ uses MongoDB to keep track of processed `visit_id`, thus if a visit
    has been processed in the past (i.e. is repeated), the publisher just ignore it and 
    doesn't forward it to NATS. Note that this is a anomalous situation since each visit must
    have a unique ID.
    
3. __How to organize this architecture in terms of infrastructure?__
    * The services reader, publisher and subscriber can run in different hosts in different subnets. The only advice
    is that MongoDB should be placed at same subnet of publisher, because, since the publisher is the service with largest
    volume of write and read operations on MongoDB, put them at same subnet will prevent eventual slowness.
    
## 📲 Providers
This system has an easy way to write new providers to send the push notifications.
This is reached with help of factory design pattern.

The file [provider_factory.go](https://github.com/gustavolopess/PushCampaignSystem/blob/develop/src/github.com/gustavolopess/PushCampaignSystem/app/providers/factory/provider_factory.go)
contains the definition of PushNotificationProvider interface:
```go
type PushNotificationProvider interface {
    SendPushNotification(pushMessage, deviceId string) (err error)
}
``` 

Yet at [provider_factory.go](https://github.com/gustavolopess/PushCampaignSystem/blob/develop/src/github.com/gustavolopess/PushCampaignSystem/app/providers/factory/provider_factory.go)
there is a map of existent providers:
```
var existentProviders = map[string]PushNotificationProvider{
	"localytics": &providers.Localytics{},
	"mixpanel": &providers.MixPanel{},
}
```
When a request to send a push notification through provider `xyz` arrives, the system looks for this provider at 
`existentProviders` map. If it is exists, an object of type of desired provider is returned and 
the method `SendPushNotification` is called from it. If it does not exist, 
an error saying that required provider is invalid is returned.

The method `SendPushNotification` contains the provider's particular implementation to send a push notification to
the `deviceId` passed as argument.

Therefore, to add a new provider to system just follow these two steps:
1. Write a struct which implements the `PushNotificationProvider` interface.
2. Append to `existentProviders` map a new reference to an object of the new provider. 
Example: `"newprovider": &newproviderpackage.NewProvider{}`
    
## 🔨 How to build

This repository provides a Makefile to helps the process of building.
You just have to execute this command in your terminal:
```bash
$ make build
```
This will generate three binary files:
1. __./bin/reader__ - the reader service
2. __./bin/pub__ - the publisher service
3. __./bin/sub__ - the subscriber service
    
## 🏃🏽‍♀️ How to run

#### Reader
The reader binary requires two params to execute: 
* `campaignfile` - the path to JSON file with campaigns to be loaded into
database
* `mongoconfig` - the path to JSON file with MongoDB config

Example:
```bash
$ ./bin/reader -campaignfile=input/activeCampaigns.json -mongoconfig=etc/mongoConfig.json
```

#### Publisher
The publisher binary requires three params to execute:
* `natsconfig` - the path to JSON file with NATS config
* `mongoconfig` - the path to JSON file with MongoDB config
* `visitlogpath` - the path to file with log of visits

Example:
```bash
$ ./bin/publisher -natsconfig=etc/natsConfig.json -mongoconfig=etc/mongoConfig.json -visitlogpath=input/visit.log
```


#### Subscriber
The subscriber binary one param to execute:
* `natsconfig` - the path to JSON file with NATS config

Example:
```bash
$  ./bin/subscriber -natsconfig=etc/natsConfig.json
```

#### MongoDB and NATS Streaming Server
This repository provides a docker-compose file to run 
the required Mongo and NATS services. Just execute this command:
```bash
$ docker-compose up -d
```

__PS:__ 

The `mongoconfig` file must follow this format:
```json
{
  "url": "mongodb://127.0.0.1:27018",
  "database": "pushcampaignsystem",
  "campaign_collection": "campaigns",
  "visit_collection": "visits"
}
```

The `natsconfig` file must follow this format:
```json
{
  "host": "localhost",
  "port": 4223,
  "cluster_id": "test-cluster",
  "client_id": "publish-campaign-client",
  "subject": "publish-campaign-subject",
  "durable_name": "publish-campaign-durable"
}
```

The JSON values doesn't need to be equal to the values above, these are just examples.

\* To make sure that all messages sent to NATS by publisher will be received by subscriber, it's recommended that both services are running at same time or with few seconds of difference between message publishment time and start of execution
of subscriber service.

#### 🗺 Project Backlog/Roadmap
Check our [Trello board](https://trello.com/b/Wl6WaCkk/pushcampaignsystem) with project tasks.

