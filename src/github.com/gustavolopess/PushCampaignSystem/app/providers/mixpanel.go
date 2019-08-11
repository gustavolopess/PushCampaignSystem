package providers

import (
"fmt"
"github.com/gustavolopess/PushCampaignSystem/app/model"
)

type MixPanel struct {}

func (m *MixPanel) SendPushNotification(natsMessage *model.NatsMessage) {
	fmt.Printf(`
		=> Push sent regarding visit %d\n
		===> Device ID: "%s"\n
		===> Mixpanel logging: { "message": "%s", device_id: "%s" }\n
		\n
	`, natsMessage.VisitId, natsMessage.DeviceId, natsMessage.PushMessage, natsMessage.DeviceId)
}
