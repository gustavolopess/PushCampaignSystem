package visit

import (
	"regexp"
	"strconv"
)

// Visit model
type Visit struct {
	ID			int64	`json:"id"`
	PlaceId		int64	`json:"place_id"`
	DeviceId	string	`json:"device_id"`
}

// regex expression to capture Visit in log lines
var logLineRegex = regexp.MustCompile(`(?i)Visit:[\s\,]+?id=(?P<id>\d+)[\s\,]+?device_id=(?P<device_id>\w+)[\s\,]+?place_id=(?P<place_id>\d+)`)


// Receives a log line and returns a Visit
func ParseVisitLogLine(line string) (visit *Visit, err error) {
	match := logLineRegex.FindStringSubmatch(line)
	matchedID, err := strconv.Atoi(match[1])
	matchedPlaceID, err := strconv.Atoi(match[3])
	if err != nil {
		return
	}

	visit = &Visit{
		ID: int64(matchedID),
		PlaceId: int64(matchedPlaceID),
		DeviceId: match[2],
	}

	return
}