package opengts

import (
	"strings"
	"time"

	"github.com/mguzelevich/go.nmea"

	"github.com/mguzelevich/sauron/storage"
)

type udpMessage struct {
	Timestamp *time.Time // time=%TIME

	Latitude  float64
	Longitude float64

	device struct {
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
	m.Timestamp = rmc.Timestamp

	lat, lon := rmc.Location.Float64()

	m.Latitude = lat
	m.Longitude = lon
	// Altitude:  rmc.Altitude,

	m.device.AndroidId = packet[0]
	m.device.Serial = packet[1]

	// m.Gps = Gps{
	// 	Satellites: values.Get("sat"),
	// 	Accuracy:   values.Get("acc"),
	// }

	m.Vehicle.Speed = rmc.Speed
	m.Vehicle.Direction = rmc.Direction
	return nil
}

func (m udpMessage) Telemetry() *storage.Telemetry {
	return &storage.Telemetry{
		Timestamp: m.Timestamp,
	}
}

func (m udpMessage) Device() *storage.Device {
	return &storage.Device{
		AndroidId: m.device.AndroidId,
		Serial:    m.device.Serial,
	}
}
