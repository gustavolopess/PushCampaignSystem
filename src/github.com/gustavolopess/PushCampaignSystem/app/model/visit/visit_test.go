package visit

import (
	"reflect"
	"testing"
)

func TestParseVisitFromLogLine(t *testing.T) {
	tests := []struct {
		name string
		logLine string
		expectedOutput *Visit
		wantError bool
	}{
		{
			"Decode a valid visit log (separated by comma and/or spaces)",
			"Visit: id=1234 device_id=id3123l, place_id=13",
			&Visit{
				1234,
				13,
				"id3123l",
			},
			false,
		},
		{
			"Decode a valid visit log (separated by comma and spaces)",
			"Visit: id=1234, device_id=id3123l, place_id=13",
			&Visit{
				1234,
				13,
				"id3123l",
			},
			false,
		},
		{
			"Decode a invalid visit log (without ID)",
			"Visit: id=, device_id=id3123l, place_id=13",
			nil,
			true,
		},
		{
			"Decode a invalid visit log (without device_id)",
			"Visit: id=1234, device_id=, place_id=13",
			nil,
			true,
		},
		{
			"Decode a invalid visit log (without place_id)",
			"Visit: id=1234, device_id=id3123l, place_id=",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Logf("Running %s", tt.name)
		visit, err := ParseVisitFromLogLine(tt.logLine)
		if (err != nil) != tt.wantError {
			t.Errorf("ParseVisitFromLogLine(): undesired error behaviour (%v, wantError=%v)", err, tt.wantError)
		} else if !tt.wantError && !reflect.DeepEqual(*visit, *tt.expectedOutput) {
			t.Errorf("ParseVisitFromLogLine(): decoded visit is different from expected (Decoded=%v, Expected=%v)", visit, tt.expectedOutput)
		}

	}
}
