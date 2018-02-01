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
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	buff, _ := json.Marshal(e)
	s.db[key][e.Pk()] = string(buff)
	return e, nil
}

func (s StorageMemory) Read(e storage.Entity) (storage.Entity, error) {
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}
	data, ok := s.db[key][e.Pk()]
	if !ok {
		return e, storage.ErrEntityNotFound
	}
	err = json.Unmarshal([]byte(data), e)
	return e, nil
}

func (s StorageMemory) Update(e storage.Entity) (storage.Entity, error) {
	key, err := mapByType(e)
	if err != nil {
		return e, err
	}

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
	s.db[key][e.Pk()] = string(buff)
	return e, nil
}

func (s StorageMemory) Delete(e storage.Entity) error {
	return nil
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

func (s StorageMemory) CreateAccount(account *storage.Account) (*storage.Account, error) {
	buff, _ := json.Marshal(account)
	s.db["accounts"][account.Id] = string(buff)
	return account, nil
}

func (s StorageMemory) ReadAccount(account *storage.Account) (*storage.Account, error) {
	data, ok := s.db["accounts"][account.Id]
	if !ok {
		return account, storage.ErrEntityNotFound
	}
	err := json.Unmarshal([]byte(data), account)
	return account, err
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
	account, err := s.ReadAccount(account)
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
