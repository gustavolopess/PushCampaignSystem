package natsmessage

import (
	"reflect"
	"testing"
)

func TestLoadMessage(t *testing.T) {
	tests := []struct{
		name string
		message []byte
		natsMessage *NatsMessage
		wantError bool
	}{
		{
			"LoadMesssage with valid message #1",
			[]byte(`{"visit_id": 1234, "provider": "localytics", "push_message": "lorem ipsum dolor", "device_id": "dev1234", "has_campaign": true}`),
			&NatsMessage{
				1234,
				"localytics",
				"lorem ipsum dolor",
				"dev1234",
				true,
			},
			false,
		},
		{
			"LoadMesssage with valid message #2",
			[]byte(`{"visit_id": 12345, "provider": "localytics", "push_message": "lorem ipsum dolor", "device_id": "dev12345", "has_campaign": false}`),
			&NatsMessage{
				12345,
				"localytics",
				"lorem ipsum dolor",
				"dev12345",
				false,
			},
			false,
		},
		{
			"LoadMesssage with invalid message (missing provider)",
			[]byte(`{"visit_id": 12345, "push_message": "lorem ipsum dolor", "device_id": "dev12345", "has_campaign": false}`),
			nil,
			true,
		},
		{
			"LoadMesssage with invalid message (missing push_message)",
			[]byte(`{"visit_id": 12345, "provider": "localytics", "device_id": "dev12345", "has_campaign": false}`),
			nil,
			true,
		},
		{
			"LoadMesssage with invalid message (missing visit_id)",
			[]byte(`{"provider": "localytics", "push_message": "lorem ipsum dolor", "device_id": "dev12345", "has_campaign": false}`),
			nil,
			true,
		},
		{
			"LoadMesssage with invalid message (missing device_id)",
			[]byte(`{"visit_id": 12345, "provider": "localytics", "push_message": "lorem ipsum dolor", "has_campaign": false}`),
			nil,
			true,
		},
		{
			"LoadMesssage with valid message #3",
			[]byte(`{"visit_id": 1235, "provider": "localytics", "push_message": "lorem ipsum dolor", "device_id": "dev1235"}`),
			&NatsMessage{
				1235,
				"localytics",
				"lorem ipsum dolor",
				"dev1235",
				false,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		tNatsMessage, err := LoadMessage(tt.message)

		if (err != nil) != tt.wantError {
			t.Errorf("LoadMessage(): Error behaviour different form expected (%v, wantErr=%v)", err, tt.wantError)
			continue
		}

		if tt.natsMessage != nil && !reflect.DeepEqual(*tt.natsMessage, tNatsMessage) {
			t.Errorf("LoadMessage(): returned message different from expected (returned=%v, expected=%v)", tNatsMessage, tt.natsMessage)
		}
	}
}