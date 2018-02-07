package storage

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"strings"
)

type Device struct {
	Id               string           `json:"id"`
	UserId           string           `json:"user_id"`
	AndroidId        string           `json:"android_id"`
	Serial           string           `json:"serial"`
	TelemetryStorage TelemetryStorage `json:"-"`
}

func (d Device) Type() string {
	return "Device"
}

func (d Device) Pk() string {
	if d.Id == "" {
		// d.Id = "00000000-0000-0000-0000-000000000000"
		d.Id = d.Hash()
	}
	return d.Id
}

func (d *Device) MarshalJSON() ([]byte, error) {
	d.Id = d.Pk()

	type alias Device
	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(d),
	})
}

func (d *Device) Telemetry() ([]*Telemetry, error) {
	return nil, nil
}

func (d *Device) AddTelemetry(telemetry *Telemetry) error {
	err := dataStorage.engine.AddTelemetry(d, telemetry)
	return err
}

func (d *Device) Hash() string {
	hash := []string{}
	for _, s := range []string{d.AndroidId, d.Serial} {
		if s != "" {
			hash = append(hash, s)
		}
	}
	checksum := crc32.Checksum([]byte(strings.Join(hash, "/")), crc32q)
	crc := fmt.Sprintf("%08x", checksum)
	return crc
}
