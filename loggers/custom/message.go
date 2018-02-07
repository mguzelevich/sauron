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

	Latitude  float64 // lat=%LAT
	Longitude float64 // lon=%LON
	Altitude  float64 // alt=%ALT

	Battery   string // battery=%BATT
	AndroidId string // androidId=%AID
	Serial    string // serial=%SER

	Satellites string // sat=%SAT
	Accuracy   string // acc=%ACC

	Speed     float64 // spd=%SPD
	Direction float64 // dir=%DIR

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

		m.Latitude = lat
		m.Longitude = lon
		m.Altitude = alt

		m.Battery = values.Get("batt")
		m.AndroidId = values.Get("aid")
		m.Serial = values.Get("ser")

		m.Satellites = values.Get("sat")
		m.Accuracy = values.Get("acc")

		speed, _ := strconv.ParseFloat(values.Get("spd"), 64)
		direction, _ := strconv.ParseFloat(values.Get("dir"), 64)

		m.Speed = speed
		m.Direction = direction

		m.Annotation = values.Get("desc")
		m.Activity = values.Get("act")
	}
	return nil
}

func (m message) Telemetry() *storage.Telemetry {
	t := &storage.Telemetry{
		Timestamp: m.Timestamp,
		Provider:  m.Provider,
		Location: storage.Location{
			Latitude:  m.Latitude,
			Longitude: m.Longitude,
			Altitude:  m.Altitude,
		},

		Annotation: m.Annotation,
		Activity:   m.Activity,
	}

	// t.Device.Id:
	t.Device.Battery = m.Battery
	t.Gps.Satellites = m.Satellites
	t.Gps.Accuracy = m.Accuracy
	t.Vehicle.Speed = m.Speed
	t.Vehicle.Direction = m.Direction

	return t
}

func (m message) Device() *storage.Device {
	return &storage.Device{
		AndroidId: m.AndroidId,
		Serial:    m.Serial,
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
