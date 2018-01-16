package storage

import ()

type Account struct {
	Id        string   `json:"id"`
	FirstName string   `json:"first_name"`
	Devices   []Device `json:"devices"`
}

func (a *Account) GetDevices() []*Device {
	return nil
}

func (a *Account) CreateDevice(device *Device) (*Device, error) {
	device, err := dataStorage.engine.CreateDevice(a, device)
	return device, err
}

func (a *Account) ReadDevice(device *Device) (*Device, error) {
	device, err := dataStorage.engine.ReadDevice(a, device)
	return device, err
}

func (a *Account) UpdateDevice(device *Device) (*Device, error) {
	return nil, nil

}

func (a *Account) DeleteDevice(device *Device) error {
	return nil
}
