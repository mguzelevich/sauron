package mmap

import (
	"encoding/json"

	"github.com/mguzelevich/go-log"

	"github.com/mguzelevich/sauron/storage"
)

func mapByType(e storage.Entity) (string, error) {
	m := map[string]string{
		"Account": "accounts",
		"Device":  "devices",
	}
	mapKey, ok := m[e.Type()]
	if !ok {
		return "", storage.ErrCollectionNotFound
	}
	return mapKey, nil
}

func (s StorageMemory) Create(e storage.Entity) (storage.Entity, error) {
	log.Trace.Printf("Create %v: %v", e.Type(), e)
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	pk := e.Pk()
	buff, _ := json.Marshal(e)
	s.db[key][pk] = string(buff)

	s.DumpCollection(e)
	return e, nil
}

func (s StorageMemory) Read(e storage.Entity) (storage.Entity, error) {
	log.Trace.Printf("Read %v: %v", e.Type(), e)
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	pk := e.Pk()
	data, ok := s.db[key][pk]
	if !ok {
		return e, storage.ErrEntityNotFound
	}
	err = json.Unmarshal([]byte(data), e)
	return e, nil
}

func (s StorageMemory) Update(e storage.Entity) (storage.Entity, error) {
	log.Debug.Printf("Update %v: %v", e.Type(), e)
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	pk := e.Pk()
	buff, _ := json.Marshal(e)

	e, err = s.Read(e)
	if err != nil {
		return e, err
	}

	err = json.Unmarshal(buff, &e)
	if err != nil {
		return e, err
	}

	buff, _ = json.Marshal(e)
	s.db[key][pk] = string(buff)
	return e, nil
}

func (s StorageMemory) Delete(e storage.Entity) (storage.Entity, error) {
	log.Trace.Printf("Delete %v: %v", e.Type(), e)
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	pk := e.Pk()
	if _, ok := s.db[key][pk]; !ok {
		return e, storage.ErrEntityNotFound
	} else {
		delete(s.db[key], pk)
	}
	return nil, nil
}

func (s StorageMemory) DumpCollection(e storage.Entity) (string, error) {
	key, err := mapByType(e)
	if err != nil {
		return "{}", err
	}
	buff, _ := json.Marshal(s.db[key])
	log.Debug.Printf("Dump %v: %v", e.Type(), string(buff))
	return string(buff), nil
}

func (s StorageMemory) Accounts() ([]*storage.Account, error) {
	accounts := []*storage.Account{}
	log.Debug.Printf("accounts")
	for k, v := range s.db["accounts"] {
		log.Debug.Printf("account %s = %s", k, v)
		accounts = append(accounts, &storage.Account{})
	}
	return accounts, nil
}

func (s StorageMemory) GetDevice(device *storage.Device) (*storage.Device, error) {
	data, ok := s.db["devices"][device.Id]
	if !ok {
		return device, storage.ErrEntityNotFound
	}
	err := json.Unmarshal([]byte(data), device)
	return device, err
}

func (s StorageMemory) ReadDevice(account *storage.Account, device *storage.Device) (*storage.Device, error) {
	e, err := s.Read(account)
	account = e.(*storage.Account)
	if err != nil {
		return device, err
	}
	device, err = s.GetDevice(device)
	return device, err
}

func (s StorageMemory) CreateDevice(account *storage.Account, device *storage.Device) (*storage.Device, error) {
	buff, _ := json.Marshal(device)
	s.db["devices"][device.Id] = string(buff)
	return device, nil
}

func (s StorageMemory) GetDevices(a *storage.Account) []*storage.Device {
	devices := []*storage.Device{}
	log.Debug.Printf("devices")
	for k, v := range s.db["devices"] {
		log.Debug.Printf("account %s = %s", k, v)
		devices = append(devices, &storage.Device{Id: k})
	}
	return devices
}
