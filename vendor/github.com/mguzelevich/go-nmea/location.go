package nmea

import (
	"fmt"
	"strconv"
	// "strings"
	// "time"
)

const (
	multiplier = 10000000000
)

type Location struct {
	latitude  int64
	longitude int64
}

func (l *Location) Equal(otherLocation *Location) bool {
	return otherLocation != nil && l.latitude == otherLocation.latitude && l.longitude == otherLocation.longitude
}

func (l *Location) Float64() (float64, float64) {
	return float64(l.latitude) / multiplier, float64(l.longitude) / multiplier
}

func strToLatLon(str string, sign string) (float64, error) {
	var err error

	f, _ := strconv.ParseFloat(str, 64)
	deg := int(f*multiplier) / multiplier / 100
	value := float64(deg) + (f-float64(deg*100))/60

	switch sign {
	case "N":
		err = nil
	case "E":
		err = nil
	case "S":
		value = -value
	case "W":
		value = -value
	default:
		err = fmt.Errorf("unknown sign %s", sign)
	}
	return value, err
}

func strToUint(str string, sign string) (int64, error) {
	fv, err := strToLatLon(str, sign)
	if err != nil {
		return 0, err
	}
	return int64(fv * multiplier), nil
}

func nmeaToLocation(lat string, latSign string, lon string, lonSign string) (*Location, error) {
	if latSign != "N" && latSign != "S" {
		return nil, fmt.Errorf("incorrect lat sign")
	}
	if lonSign != "E" && lonSign != "W" {
		return nil, fmt.Errorf("incorrect lon sign")
	}

	ll := &Location{}
	ll.latitude, _ = strToUint(lat, latSign)
	ll.longitude, _ = strToUint(lon, lonSign)
	return ll, nil
}
