package sauron

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"
	//	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	files         map[string]*os.File
	db            *sql.DB
	telemetryChan chan *Telemetry
}

func (s *Storage) loop() {
	for loc := range s.telemetryChan {
		s.saveLocation(loc)
	}
}

func (s *Storage) saveLocation(t *Telemetry) {
	f, err := s.getFile(t.Device.Hash())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
	} else {
		if out, err := json.Marshal(t); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(f, "%s\n", string(out))
		}
	}
}

func (s *Storage) Save(telemetry *Telemetry) {
	go func() {
		s.telemetryChan <- telemetry
	}()
}

func (s *Storage) getFile(uid string) (*os.File, error) {
	if uid == "" {
		uid = "default"
	}
	file, ok := s.files[uid]
	if !ok {
		// s.files[uid] = nil
		timestamp := time.Now().UTC().Format("20060102")
		// filename := fmt.Sprintf("/tmp/sauron.%s.log", uid)
		filename := fmt.Sprintf("sauron.%s.%s.log", timestamp, uid)
		if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			return nil, fmt.Errorf("ERROR %v\n", err)
		} else {
			s.files[uid] = file
			fmt.Fprintf(os.Stderr, "file created: [%s]\n", filename)
		}
	} else {
		// fmt.Fprintf(os.Stderr, "use existed file: %s\n", uid)
	}
	return file, nil
}

func NewStorage() *Storage {
	// storageFile := ":memory:"

	storage := &Storage{
		files:         make(map[string]*os.File),
		telemetryChan: make(chan *Telemetry),
	}
	go storage.loop()

	// if db, err := sql.Open("sqlite3", storageFile); err != nil {
	// 	fmt.Fprintf(os.Stderr, "ERROR: %v", err)
	// } else {
	// 	storage.db = db
	// }

	return storage
}
