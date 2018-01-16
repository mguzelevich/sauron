package storage

import (
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

func ChangeAccount() {

}

func (d *Device) Telemetry() ([]*Telemetry, error) {
	return nil, nil
}

func (d *Device) AddTelemetry(telemetry *Telemetry) error {
	return nil
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
