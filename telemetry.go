package sauron

import (
	"net/url"
	"time"
)

// %LAT - Latitude
// %LON - Longitude
// %DESC - Annotation
// %SAT - Satellites
// %ALT - Altitude
// %SPD - Speed
// %ACC - Accuracy
// %DIR - Direction
// %PROV - Provider
// %TIMESTAMP - Timestamp (epoch)
// %TIME - Time UTC (2011-12-25T15:27:33Z)
// %STARTTIMESTAMP - Start timestamp (epoch)
// %BATT - Battery
// %AID - Android ID
// %SER - Serial
// %ACT - Activity
// %FILENAME - summary_current_filename

func parseTime(ts string) time.Time {
	// 2011-12-25T15:27:33Z
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	if t, err := time.Parse(time.RFC3339, ts); err == nil {
		return t
	}
	return time.Now().UTC()
}

type Location struct {
	Latitude  string `json:"lat"`           // lat=%LAT
	Longitude string `json:"lon"`           // lon=%LON
	Altitude  string `json:"alt,omitempty"` // alt=%ALT
}

type Vehicle struct {
	Speed     string `json:"spd,omitempty"` // spd=%SPD
	Direction string `json:"dir,omitempty"` // dir=%DIR
}

type Device struct {
	Battery   string `json:"battery,omitempty"`   // battery=%BATT
	AndroidId string `json:"androidId,omitempty"` // androidId=%AID
	Serial    string `json:"serial,omitempty"`    // serial=%SER
}

func (d *Device) Hash() string {
	return d.Serial
}

type Gps struct {
	Satellites string `json:"sat,omitempty"` // sat=%SAT
	Accuracy   string `json:"acc,omitempty"` // acc=%ACC
}

type Telemetry struct {
	Timestamp time.Time `json:"time,omitempty"` // time=%TIME
	Provider  string    `json:"prov,omitempty"` // prov=%PROV

	Location Location `json:"location"`
	Device   Device   `json:"device"`
	Gps      Gps      `json:"gps"`
	Vehicle  Vehicle  `json:"vehicle"`

	Annotation string `json:"desc,omitempty"` // desc=%DESC

	Activity string `json:"activity,omitempty"` // activity=%ACT
	// Epoch    string `json:"epoch,omitempty"`    // epoch=%TIMESTAMP"
}

func (t *Telemetry) ParseCustomUrl(p string) error {
	if values, err := url.ParseQuery(p); err != nil {
		return err
	} else {
		t.Provider = values.Get("prov")
		t.Timestamp = parseTime(values.Get("time"))

		t.Location = Location{
			Latitude:  values.Get("lat"),
			Longitude: values.Get("lon"),
			Altitude:  values.Get("alt"),
		}

		t.Device = Device{
			Battery:   values.Get("batt"),
			AndroidId: values.Get("aid"),
			Serial:    values.Get("ser"),
		}

		t.Gps = Gps{
			Satellites: values.Get("sat"),
			Accuracy:   values.Get("acc"),
		}

		t.Vehicle = Vehicle{
			Speed:     values.Get("spd"),
			Direction: values.Get("dir"),
		}

		t.Annotation = values.Get("desc")
		t.Activity = values.Get("act")
	}
	return nil
}

func (t *Telemetry) ParseUdpPacket(p string) error {
	// data := "mgu/mi5s/$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E"
	// re := "(.+)/(.+)/($GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E)"
	return nil
}
