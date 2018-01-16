package bolt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

var storage *StorageBolt

func TestStorageBolt_init(t *testing.T) {
	meta, err := Storage.Meta()
	if err != nil {
		t.Error(err)
	}
	if meta.Version != "0.1" {
		t.Error("incorrect version", meta)
	}
}

func TestStorageBolt_SaveTelemetry(t *testing.T) {
	timestamp := time.Now().UTC()
	telemetry := &Telemetry{
		Timestamp: &timestamp,
	}
	f := func(tx *bolt.Tx) error {
		return storage.saveTelemetry(tx, telemetry)
	}
	if err := storage.db.Update(f); err != nil {
		t.Error(err)
	}
}

func setup() string {
	tmpDir := ""
	if dir, err := ioutil.TempDir("", "sauron"); err != nil {
		fmt.Fprintf(os.Stderr, "setup() %v\n", err)
		os.Exit(1)
	} else {
		tmpDir = dir
	}
	tmpDbFilename := filepath.Join(tmpDir, "tmpdb.db")

	shutdownChan := make(chan bool)
	storageBolt, err := NewBolt(tmpDbFilename, shutdownChan)
	if err != nil {
		shutdown(tmpDir)
		fmt.Fprintf(os.Stderr, "setup() storage init error %v\n", err)
		os.Exit(1)
	}
	storage = storageBolt
	Storage = storageBolt
	return tmpDir
}

func shutdown(tmpDir string) {
	if 1 == 1 {
		os.RemoveAll(tmpDir)
	}
}

func TestMain(m *testing.M) {
	// http://cs-guy.com/blog/2015/01/test-main/
	tmpDir := setup()
	code := m.Run()
	shutdown(tmpDir)
	os.Exit(code)
}
