package storage

import ()

type Entity interface {
	Type() string
	Pk() string
}

type StorageEngine interface {
	DoneChan() chan bool

	Create(e Entity) (Entity, error)
	Read(e Entity) (Entity, error)
	Update(e Entity) (Entity, error)
	Delete(e Entity) (Entity, error)

	// DumpAll - dump all data as json
	DumpAll() ([]byte, error)

	// Accounts - return all accounts
	Accounts() ([]*Account, error)

	// Telemetry
	// Telemetry(d *Device) ([]*Telemetry, error)

	// AddTelemetry
	AddTelemetry(d *Device, telemetry *Telemetry) error

	// Add
	// Add(t *TelemetryStorage, telemetry *Telemetry)

	// Read
	// Read(t *TelemetryStorage, from time.Time, to time.Time) []*Telemetry

	// All
	// All(t *TelemetryStorage) []*Telemetry
}

type NewFunc func(params map[string]string, shutdownChan chan bool) (StorageEngine, error)
