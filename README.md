# PushCampaignSystem 🛰🌎📲
Ad toy system which allows clients to impact a mobile device when it is at a specific place with a push notification

## 📃 Description

Clients can contract ad campaigns based on geolocation. These campaigns are defined by an identifier,
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
any campaign. The Golang struct of the message that can be sent to NATS is defined below:
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
