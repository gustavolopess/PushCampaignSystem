package providers

import (
	"fmt"
	"github.com/gustavolopess/PushCampaignSystem/app/model"
)

type Localytics struct {}

func (l *Localytics) SendPushNotification(natsMessage *model.NatsMessage) {
	fmt.Printf(`
		=> Push sent regarding visit %d\n
		===> Device ID: "%s"\n
		===> Localytics logging: { "message": "%s", device_id: "%s" }\n
		\n
	`, natsMessage.VisitId, natsMessage.DeviceId, natsMessage.PushMessage, natsMessage.DeviceId)
}
