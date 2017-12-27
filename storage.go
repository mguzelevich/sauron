package sauron

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Storage struct {
	files map[string]*os.File
}

func (s *Storage) Save(l *Location) {
	f, err := s.getFile(l.Serial)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
	} else {
		if out, err := json.Marshal(l); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(f, "%s\n", string(out))
		}
	}
}

func (s *Storage) getFile(uid string) (*os.File, error) {
	if uid == "" {
		uid = "default"
	}
	file, ok := s.files[uid]
	if !ok {
		fmt.Fprintf(os.Stderr, "create file: %s\n", uid)

		// s.files[uid] = nil
		timestamp := time.Now().UTC().Format("20060102")
		// filename := fmt.Sprintf("/tmp/sauron.%s.log", uid)
		filename := fmt.Sprintf("sauron.%s.%s.log", timestamp, uid)
		if file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			return nil, fmt.Errorf("ERROR %v\n", err)
		} else {
			s.files[uid] = file
		}
	} else {
		// fmt.Fprintf(os.Stderr, "use existed file: %s\n", uid)
	}
	return file, nil
}

func NewStorage() *Storage {
	return &Storage{
		files: make(map[string]*os.File),
	}
}
