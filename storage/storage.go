package storage

import (
	//	"encoding/json"
	//	"os"
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/mguzelevich/sauron/log"
)

type Storage struct {
	db *bolt.DB

	doneChan chan bool

	saveTelemetryChan chan *Telemetry
}

func Init(filename string, shutdownChan chan bool) (*Storage, error) {
	storage := &Storage{
		doneChan:          make(chan bool),
		saveTelemetryChan: make(chan *Telemetry),
	}

	if boltDb, err := bolt.Open(filename, 0600, nil); err != nil {
		return nil, fmt.Errorf("cann't open database [%s] %s", filename, err)
	} else {
		log.Info.Printf("storage inited [%s]", filename)
		storage.db = boltDb
	}
	go storage.shutdownLoop(shutdownChan)
	go storage.saveLoop()
	storage.initSchema()
	return storage, nil
}

func (s Storage) DoneChan() chan bool {
	return s.doneChan
}

func (s *Storage) shutdownLoop(shutdownChan chan bool) {
	<-shutdownChan
	s.db.Close()
	log.Info.Printf("storage gracefully stopped")
	close(s.doneChan)
}

func (s *Storage) saveLoop() {
	for loc := range s.saveTelemetryChan {
		s.saveLocation(loc)
	}
}

func (s *Storage) initSchema() error {
	timestamp := time.Now().UTC().Format("20060102")

	s.initBuckets()

	meta := []struct {
		key   string
		value string
	}{
		{"version", "0.1"},
		{"created_at", timestamp},
		{"updated_at", timestamp},
	}

	s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(".meta"))
		for _, r := range meta {
			if err := bucket.Put([]byte(r.key), []byte(r.value)); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func (s *Storage) initBuckets() error {
	buckets := []string{
		".meta", "tasks", "files",
	}
	tasksBuckets := []string{
		"pending", "process", "finished", "failed",
	}
	if err := s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range buckets {
			bucket, err := tx.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return fmt.Errorf("create bucket: %s %s", bucket, err)
			}
		}
		tasksBucket := tx.Bucket([]byte("tasks"))
		for _, name := range tasksBuckets {
			bucket, err := tasksBucket.CreateBucketIfNotExists([]byte(name))
			if err != nil {
				return fmt.Errorf("create bucket: %s %s", bucket, err)
			}
		}
		return nil
	}); err != nil {
		log.Error.Printf("%s", err)
	}
	return nil
}

func (s *Storage) get(bucketName []string, key string) (string, error) {
	var response string
	s.db.View(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		for _, bn := range bucketName {
			if bucket == nil {
				bucket = tx.Bucket([]byte(bn))
			} else {
				bucket = bucket.Bucket([]byte(bn))
			}
		}
		response = string(bucket.Get([]byte(key)))
		return nil
	})
	return response, nil
}

func (s *Storage) Save(telemetry *Telemetry) {
	go func() {
		s.saveTelemetryChan <- telemetry
	}()
}

func (s *Storage) saveLocation(t *Telemetry) {
	// f, err := s.getFile(t.Device.Hash())
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "ERROR: %v", err)
	// } else {
	// 	if out, err := json.Marshal(t); err != nil {
	// 		// http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	} else {
	// 		fmt.Fprintf(f, "%s\n", string(out))
	// 	}
	// }
}
