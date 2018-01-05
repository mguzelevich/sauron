package messages

import (
	"fmt"
	"strings"
)

type NmeaMessage interface {
	Marshal() ([]byte, error)
	Unmarshal([]string) error
}

type CreateMessageFunc func([]string) (NmeaMessage, error)

var register map[string]CreateMessageFunc

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

func validateChecksum(checksum string) error {
	if checksum[0] != '*' {
		return fmt.Errorf("checksum error (first symbol must be '*'")
	}
	return nil
}

func ParseMessage(data []byte) (NmeaMessage, error) {
	fields, err := split(data)
	if err != nil {
		return nil, err
	}

	msgLabel := fields[0]
	checksum := fields[len(fields)-1]

	f, ok := register[msgLabel]
	if !ok {
		return nil, fmt.Errorf("unknown message type [%s]", msgLabel)
	}
	if f == nil {
		return nil, fmt.Errorf("not supportable message type [%s]", msgLabel)
	}
	if err := validateChecksum(checksum); err != nil {
		return nil, fmt.Errorf("checksum [%s] is incorrect %s", checksum, err)
	}

	return f(fields)
}

/*
http://www.gpsinformation.org/dale/nmea.htm
*/
func init() {
	register = map[string]CreateMessageFunc{
		"$GPAAM": nil,       // Waypoint Arrival Alarm
		"$GPALM": nil,       // Almanac data
		"$GPAPA": nil,       // Auto Pilot A sentence
		"$GPAPB": nil,       // Auto Pilot B sentence
		"$GPBOD": nil,       // Bearing Origin to Destination
		"$GPBWC": nil,       // Bearing using Great Circle route
		"$GPDTM": nil,       // Datum being used.
		"$GPGGA": nil,       // Fix information
		"$GPGLL": nil,       // Lat/Lon data
		"$GPGRS": nil,       // GPS Range Residuals
		"$GPGSA": nil,       // Overall Satellite data
		"$GPGST": nil,       // GPS Pseudorange Noise Statistics
		"$GPGSV": nil,       // Detailed Satellite data
		"$GPMSK": nil,       // send control for a beacon receiver
		"$GPMSS": nil,       // Beacon receiver status information.
		"$GPRMA": nil,       // recommended Loran data
		"$GPRMB": nil,       // recommended navigation data for gps
		"$GPRMC": createRmc, // recommended minimum data for gps
		"$GPRTE": nil,       // route message
		"$GPTRF": nil,       // Transit Fix Data
		"$GPSTN": nil,       // Multiple Data ID
		"$GPVBW": nil,       // dual Ground / Water Spped
		"$GPVTG": nil,       // Vector track an Speed over the Ground
		"$GPWCV": nil,       // Waypoint closure velocity (Velocity Made Good)
		"$GPWPL": nil,       // Waypoint Location information
		"$GPXTC": nil,       // cross track error
		"$GPXTE": nil,       // measured cross track error
		"$GPZTG": nil,       // Zulu (UTC) time and time to go (to destination)
		"$GPZDA": nil,       // Date and Time
	}
}
