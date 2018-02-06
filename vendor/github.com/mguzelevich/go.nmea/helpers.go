package nmea

import (
	"strconv"
	"strings"
	"time"
)

func split(data []byte) ([]string, error) {
	raw := strings.TrimSpace(string(data))
	checksum := ""
	if idx := strings.LastIndex(raw, "*"); idx != -1 {
		// checksum finded
		checksum = raw[idx:]
		raw = raw[:idx]
	}

	result := strings.Split(raw, ",")
	if checksum != "" {
		result = append(result, checksum)
	}
	return result, nil
}

func toTime(date string, timeStr string) (*time.Time, error) {
	year, _ := strconv.Atoi(date[4:6])
	month, _ := strconv.Atoi(date[2:4])
	day, _ := strconv.Atoi(date[0:2])
	hour, _ := strconv.Atoi(timeStr[0:2])
	min, _ := strconv.Atoi(timeStr[2:4])
	sec, _ := strconv.Atoi(timeStr[4:6])

	if year > 60 {
		year += 1900
	} else {
		year += 2000
	}
	ts := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	return &ts, nil
}
