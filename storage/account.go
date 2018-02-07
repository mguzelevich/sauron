package storage

import (
	"encoding/json"

	"github.com/mguzelevich/go.log"
)

type Account struct {
	Id        string   `json:"id"`
	FirstName string   `json:"first_name"`
	Devices   []string `json:"devices"`
}

func (a Account) Type() string {
	return "Account"
}

func (a Account) Pk() string {
	if a.Id == "" {
		// a.Id, _ := log.UUID()
		a.Id = "00000000-0000-0000-0000-000000000000"
	}
	return a.Id
}

func (a *Account) MarshalJSON() ([]byte, error) {
	a.Id = a.Pk()

	type alias Account
	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(a),
	})
}

func (a *Account) GetDevices() []*Device {
	return nil
}

func (a *Account) CreateDevice(device *Device) (*Device, error) {
	log.Debug.Printf("create device %v -> %v", a, device)
	device.UserId = a.Pk()
	entity, err := dataStorage.engine.Create(device)
	d := entity.(*Device)
	a.Devices = append(a.Devices, device.Pk())

	_, err = dataStorage.engine.Update(a)
	return d, err
}

func (a *Account) Device(device *Device) (*Device, error) {
	entity, err := dataStorage.engine.Read(device)
	return entity.(*Device), err
}

func (a *Account) UpdateDevice(device *Device) (*Device, error) {
	entity, err := dataStorage.engine.Update(device)
	d := entity.(*Device)
	return d, err
}

func (a *Account) DeleteDevice(device *Device) error {
	for i := range a.Devices {
		if a.Devices[i] == device.Pk() {
			a.Devices[i] = a.Devices[len(a.Devices)-1]
			a.Devices = a.Devices[:len(a.Devices)-1]
			return nil
		}
	}
	return ErrEntityNotFound
}
