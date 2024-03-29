package factory

import (
	"fmt"
	"github.com/gustavolopess/PushCampaignSystem/app/providers"
)

// Interface which defines a provider
type PushNotificationProvider interface {
	SendPushNotification(pushMessage, deviceId string) (err error)
}

// Map with existent providers, new providers must be appended
var existentProviders = map[string]PushNotificationProvider{
	"localytics": &providers.Localytics{},
	"mixpanel": &providers.MixPanel{},
	"fastmessage": &providers.FastMessage{},
}

// Return provider based on its name
func GetProvider(provider string) (PushNotificationProvider, error) {
	// Check if required provider exists
	if _, ok := existentProviders[provider]; !ok {
		return nil, fmt.Errorf("%s is not a valid provider")
	}

	return existentProviders[provider], nil
}

