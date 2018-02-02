package custom

import (
	"net/url"
	"strconv"
	"time"

	"github.com/mguzelevich/sauron/storage"
)

/*
  %LAT - Latitude
  %LON - Longitude
  %DESC - Annotation
  %SAT - Satellites
  %ALT - Altitude
  %SPD - Speed
  %ACC - Accuracy
  %DIR - Direction
  %PROV - Provider
  %TIMESTAMP - Timestamp (epoch)
  %TIME - Time UTC (2011-12-25T15:27:33Z)
  %STARTTIMESTAMP - Start timestamp (epoch)
  %BATT - Battery
  %AID - Android ID
  %SER - Serial
  %ACT - Activity
  %FILENAME - summary_current_filename
*/
type message struct {
	Timestamp *time.Time // time=%TIME
	Provider  string     // prov=%PROV

	Location struct {
		Latitude  float64 // lat=%LAT
		Longitude float64 // lon=%LON
		Altitude  float64 // alt=%ALT
	}
	device struct {
		Battery   string // battery=%BATT
		AndroidId string // androidId=%AID
		Serial    string // serial=%SER
	}
	Gps struct {
		Satellites string // sat=%SAT
		Accuracy   string // acc=%ACC
	}
	Vehicle struct {
		Speed     float64 // spd=%SPD
		Direction float64 // dir=%DIR
	}

	Annotation string // desc=%DESC

	Activity string // activity=%ACT
	// Epoch    string `json:"epoch,omitempty"`    // epoch=%TIMESTAMP"
}

func (m *message) ParseRaw(p string) error {
	if values, err := url.ParseQuery(p); err != nil {
		return err
	} else {
		m.Provider = values.Get("prov")
		m.Timestamp = parseTime(values.Get("time"))

		lat, _ := strconv.ParseFloat(values.Get("lat"), 64)
		lon, _ := strconv.ParseFloat(values.Get("lon"), 64)
		alt, _ := strconv.ParseFloat(values.Get("alt"), 64)

		m.Location.Latitude = lat
		m.Location.Longitude = lon
		m.Location.Altitude = alt

		m.device.Battery = values.Get("batt")
		m.device.AndroidId = values.Get("aid")
		m.device.Serial = values.Get("ser")

		m.Gps.Satellites = values.Get("sat")
		m.Gps.Accuracy = values.Get("acc")

		speed, _ := strconv.ParseFloat(values.Get("spd"), 64)
		direction, _ := strconv.ParseFloat(values.Get("dir"), 64)

		m.Vehicle.Speed = speed
		m.Vehicle.Direction = direction

		m.Annotation = values.Get("desc")
		m.Activity = values.Get("act")
	}
	return nil
}

func (m message) Telemetry() *storage.Telemetry {
	return &storage.Telemetry{
		Timestamp: m.Timestamp,
	}
}

func (m message) Device() *storage.Device {
	return &storage.Device{
		AndroidId: m.device.AndroidId,
		Serial:    m.device.Serial,
	}
}

func parseTime(ts string) *time.Time {
	// 2011-12-25T15:27:33Z
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	result := time.Now().UTC()
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		result = t
	}
	return &result
}
