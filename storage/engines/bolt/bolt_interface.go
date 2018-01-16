package bolt

import (
	// "os"
	// "encoding/json"
	"fmt"
	// "time"

	"github.com/boltdb/bolt"
	"github.com/mguzelevich/go-log"

	"github.com/mguzelevich/sauron/storage"
)

func (s StorageBolt) Accounts() ([]*storage.Account, error) {
	accounts := []*storage.Account{}
	log.Debug.Printf("accounts")
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("accounts"))

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			accounts = append(accounts, &storage.Account{})
		}

		return nil
	})
	return accounts, nil
}

func (s StorageBolt) CreateAccount(account *storage.Account) (*storage.Account, error) {
	return nil, nil
}

func (s StorageBolt) ReadAccount(account *storage.Account) (*storage.Account, error) {
	return nil, nil
}

func (s StorageBolt) GetDevice(d *storage.Device) (*storage.Device, error) {
	return nil, nil
}

func (s StorageBolt) CreateDevice(account *storage.Account, device *storage.Device) (*storage.Device, error) {
	return device, nil
}

func (s StorageBolt) GetDevices(a *storage.Account) []*storage.Device {
	devices := []*storage.Device{}
	return devices
}

func (s StorageBolt) ReadDevice(account *storage.Account, device *storage.Device) (*storage.Device, error) {
	return device, nil
}