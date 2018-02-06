package storage

import (
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

func (a *Account) GetDevices() []*Device {
	return nil
}

func (a *Account) CreateDevice(device *Device) (*Device, error) {
	log.Debug.Printf("create device %v -> %v", a, device)
	device.Id = device.Pk()
	entity, err := dataStorage.engine.Create(device)
	d := entity.(*Device)
	a.Devices = append(a.Devices, device.Id)

	_, err = dataStorage.engine.Update(a)
	return d, err
}

func (a *Account) ReadDevice(device *Device) (*Device, error) {
	device, err := dataStorage.engine.ReadDevice(a, device)
	return device, err
}

func (a *Account) UpdateDevice(device *Device) (*Device, error) {
	entity, err := dataStorage.engine.Update(device)
	d := entity.(*Device)
	return d, err
}

func (a *Account) DeleteDevice(device *Device) error {
	return nil
}
