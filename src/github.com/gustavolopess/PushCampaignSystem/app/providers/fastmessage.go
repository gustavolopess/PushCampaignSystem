package providers

import "fmt"

type FastMessage struct{}

func (f *FastMessage) SendPushNotification(pushMessage, deviceId string) (err error) {
	// Write provider's delivery logic here
	fmt.Println("FastMessage provider!!")
	return
}
