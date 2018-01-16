package dummy

import (
	"github.com/mguzelevich/sauron/storage"
)

type StorageDummy struct {
}

func (s StorageDummy) Meta() (*storage.MetaInfo, error) {
	panic("not implemented")
	return nil, nil
}

func (s StorageDummy) Accounts() []string {
	panic("not implemented")
	return nil
}

func (s StorageDummy) SaveTelemetry(telemetry *storage.Telemetry) {
	panic("not implemented")
	return
}

func (s StorageDummy) GetTelemetry(key string) {
	panic("not implemented")
	return
}

func (s StorageDummy) Dump() error {
	panic("not implemented")
	return nil
}
