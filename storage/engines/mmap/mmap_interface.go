package mmap

import (
	"encoding/json"
	"fmt"

	"github.com/mguzelevich/go.log"

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
	s.db[key][pk] = json.RawMessage(buff)

	// s.DumpCollection(e)
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
	s.db[key][pk] = json.RawMessage(buff)
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
	buff, err := json.Marshal(s.db[key])
	log.Debug.Printf("Dump %v: %v (%v)", e.Type(), string(buff), err)
	return string(buff), err
}

func (s StorageMemory) DumpAll() ([]byte, error) {
	return s.dump()
}

func (s StorageMemory) Accounts() ([]*storage.Account, error) {
	key, err := mapByType(storage.Account{})
	if err != nil {
		return nil, err
	}

	accounts := []*storage.Account{}
	for k := range s.db[key] {
		if entity, err := s.Read(&storage.Account{Id: k}); err != nil {
			log.Error.Printf("read account %v error %v", k, err)
		} else {
			accounts = append(accounts, entity.(*storage.Account))
		}
	}
	return accounts, nil
}

func (s StorageMemory) AddTelemetry(d *storage.Device, telemetry *storage.Telemetry) error {
	key := fmt.Sprintf("telemetry.%v", d.Pk())
	log.Trace.Printf("AddTelemetry %v: %v: %v", d, key, telemetry)
	telemetry.Device.Id = d.Pk()
	if _, ok := s.db[key]; !ok {
		s.db[key] = make(map[string]json.RawMessage)
	}
	if buff, err := json.Marshal(telemetry); err != nil {
		log.Error.Printf("save device %v telemetry %v error %v", d.Pk(), telemetry, err)
	} else {
		s.db[key][telemetry.Pk()] = json.RawMessage(buff)
		// log.Error.Printf("save device %v telemetry [%v]", d.Pk(), string(buff))
	}
	return nil
}
