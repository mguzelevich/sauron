package bolt

import (
	// "os"
	// "encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/mguzelevich/go-log"

	"github.com/mguzelevich/sauron/storage"
)

type StorageBolt struct {
	db *bolt.DB

	doneChan chan bool
}

func (s StorageBolt) DoneChan() chan bool {
	return s.doneChan
}

func (s *StorageBolt) shutdownLoop(shutdownChan chan bool) {
	<-shutdownChan
	s.db.Close()
	log.Info.Printf("bolt storage engine gracefully stopped")
	close(s.doneChan)
}

func (s *StorageBolt) initSchema() error {
	timestamp := time.Now().UTC().Format("20060102")

	s.initBuckets()

	meta := []struct {
		bucket []string
		key    string
		value  string
	}{
		{[]string{".meta"}, "version", "0.1"},
		{[]string{".meta"}, "created_at", timestamp},
		{[]string{".meta"}, "updated_at", timestamp},
	}

	s.db.Update(func(tx *bolt.Tx) error {
		txc := NewTx(tx)
		for _, r := range meta {
			if err := txc.put(r.bucket, r.key, r.value); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (s *StorageBolt) initBuckets() error {
	buckets := [][]string{
		[]string{".meta"},
		[]string{"accounts"},
		[]string{"devices"},
		[]string{"telemetry"},
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		txc := NewTx(tx)
		for _, names := range buckets {
			if _, err := txc.getOrCreateBucket(names); err != nil {
				return fmt.Errorf("create bucket: %v %v", names, err)
			}
		}
		return nil
	}); err != nil {
		log.Error.Printf("%s", err)
	}
	return nil
}

func (s *StorageBolt) Meta() (*storage.MetaInfo, error) {
	meta := &storage.MetaInfo{}
	f := func(tx *bolt.Tx) error {
		data, err := NewTx(tx).get([]string{".meta"}, "version")
		if err != nil {
			return err
		}
		if data == "" {
			return fmt.Errorf("no data")
		}
		meta.Version = data
		return nil
	}
	err := s.db.View(f)
	return meta, err
}

func (s StorageBolt) Dump() error {
	s.db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("telemetry"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
	return nil
}

// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

func New(params map[string]string, shutdownChan chan bool) (storage.StorageEngine, error) {
	storage := &StorageBolt{
		doneChan: make(chan bool),
	}

	filename := params["filename"]
	if boltDb, err := bolt.Open(filename, 0600, nil); err != nil {
		return nil, fmt.Errorf("cann't open database [%s] %s", filename, err)
	} else {
		storage.db = boltDb
		storage.initSchema()
		log.Info.Printf("bolt storage engine inited [%s]", filename)
	}

	go storage.shutdownLoop(shutdownChan)

	return storage, nil
}
