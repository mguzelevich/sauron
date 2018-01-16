package mock

import (
	"github.com/mguzelevich/sauron/storage"
)

type StorageMock struct {
}

func (s StorageMock) Meta() (*storage.MetaInfo, error) {
	return &MetaInfo{}, nil
}

func (s StorageMock) Accounts() []string {
	return []string{}
}

func (s StorageMock) SaveTelemetry(telemetry *storage.Telemetry) {
	return
}

func (s StorageMock) GetTelemetry(key string) {
	return
}

func (s StorageMock) Dump() error {
	return nil
}
