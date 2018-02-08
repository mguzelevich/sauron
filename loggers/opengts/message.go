package opengts

import (
	"strings"
	"time"

	"github.com/mguzelevich/go.nmea"

	"github.com/mguzelevich/sauron/storage"
)

type udpMessage struct {
	Timestamp *time.Time

	OriginalTimestamp *time.Time // time=%TIME

	Latitude  float64
	Longitude float64

	AndroidId string // androidId=%AID
	Serial    string // serial=%SER

	Speed     float64 // spd=%SPD
	Direction float64 // dir=%DIR
}

func (m *udpMessage) ParseRaw(p string) error {
	// data := "mgu/mi5s/$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E"
	packet := strings.Split(p, "/")
	//account := packet[0]
	//t.Device.DeviceId = packet[1]
	nmeaPacket := packet[2]
	// re := "(.+)/(.+)/($GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E)"

	msg, err := nmea.NewPacket([]byte(nmeaPacket))
	if err != nil {
		return err
	}
	rmc, err := msg.AsRmc()
	if err != nil {
		return err
	}
	// t.Provider = values.Get("prov")
	m.OriginalTimestamp = rmc.Timestamp

	lat, lon := rmc.Location.Float64()

	m.Latitude = lat
	m.Longitude = lon
	// Altitude:  rmc.Altitude,

	m.AndroidId = packet[0]
	m.Serial = packet[1]

	m.Speed = rmc.Speed
	m.Direction = rmc.Direction
	return nil
}

func (m udpMessage) Telemetry() *storage.Telemetry {
	t := &storage.Telemetry{
		Timestamp:         m.Timestamp,
		OriginalTimestamp: m.OriginalTimestamp,
		// Provider:  m.Provider,
		Location: storage.Location{
			Latitude:  m.Latitude,
			Longitude: m.Longitude,
			// Altitude:  m.Altitude,
		},

		// Annotation: m.Annotation,
		// Activity:   m.Activity,
	}

	// t.Device.Id:
	// t.Device.Battery = m.Battery
	// t.Gps.Satellites = m.Satellites
	// t.Gps.Accuracy = m.Accuracy
	t.Vehicle.Speed = m.Speed
	t.Vehicle.Direction = m.Direction

	return t
}

func (m udpMessage) Device() *storage.Device {
	return &storage.Device{
		AndroidId: m.AndroidId,
		Serial:    m.Serial,
	}
}
