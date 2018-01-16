package storage

import (
	"time"
)

type TelemetryStorage struct {
}

func (t *TelemetryStorage) Add(telemetry *Telemetry) {

}

func (t *TelemetryStorage) Read(from time.Time, to time.Time) []*Telemetry {
	return nil
}

func (t *TelemetryStorage) All() []*Telemetry {
	return nil
}
