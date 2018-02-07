package bolt

import (
	// "os"
	// "encoding/json"
	"fmt"
	// "time"

	"github.com/boltdb/bolt"
	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/storage"
)

func (s StorageBolt) Create(e storage.Entity) (storage.Entity, error) {
	return e, nil
}

func (s StorageBolt) Read(e storage.Entity) (storage.Entity, error) {
	return e, nil
}

func (s StorageBolt) Update(e storage.Entity) (storage.Entity, error) {
	return e, nil
}

func (s StorageBolt) Delete(e storage.Entity) (storage.Entity, error) {
	return e, nil
}

func (s StorageBolt) DumpAll() ([]byte, error) {
	return nil, nil
}

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

func (s StorageBolt) AddTelemetry(d *storage.Device, telemetry *storage.Telemetry) error {
	log.Trace.Printf("AddTelemetry %v: %v", d, telemetry)
	return nil
}
