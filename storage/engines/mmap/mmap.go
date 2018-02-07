package mmap

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/storage"
)

type StorageMemory struct {
	db map[string]map[string]json.RawMessage

	doneChan chan bool
}

func (s StorageMemory) DoneChan() chan bool {
	return s.doneChan
}

func (s *StorageMemory) shutdownLoop(shutdownChan chan bool) {
	<-shutdownChan
	// s.db.Close()
	if d, err := s.dump(); err != nil {
		log.Error.Printf("mmem storage dump error %v", err)
	} else {
		log.Stdout.Printf("%v", string(d))
	}
	log.Info.Printf("mmem storage engine gracefully stopped")
	close(s.doneChan)
}

func (s *StorageMemory) initSchema() error {
	//timestamp := time.Now().UTC().Format("20060102")
	timestamp := fmt.Sprintf("\"%s\"", time.Now().UTC().Format(time.RFC3339))

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

	for _, r := range meta {
		key := strings.Join(r.bucket, "/")
		s.db[key][r.key] = json.RawMessage(r.value)
	}
	return nil
}

func (s *StorageMemory) initBuckets() error {
	buckets := [][]string{
		[]string{".meta"},
		[]string{"accounts"},
		[]string{"devices"},
		[]string{"telemetry"},
	}
	for _, names := range buckets {
		key := strings.Join(names, "/")
		s.db[key] = make(map[string]json.RawMessage)
	}
	return nil
}

func (s *StorageMemory) dump() ([]byte, error) {
	buff, err := json.Marshal(s.db)
	return buff, err
}

// https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

func New(params map[string]string, shutdownChan chan bool) (storage.StorageEngine, error) {
	storage := &StorageMemory{
		db:       make(map[string]map[string]json.RawMessage),
		doneChan: make(chan bool),
	}

	storage.initSchema()
	log.Info.Printf("mmem storage engine inited")

	go storage.shutdownLoop(shutdownChan)

	return storage, nil
}
