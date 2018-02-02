package bolt

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	// "time"

	"github.com/mguzelevich/sauron/storage"
)

var engine storage.StorageEngine

func TestStorageBolt_init(t *testing.T) {
	var meta *storage.MetaInfo

	boltDb := engine.(*StorageBolt)
	meta, err := boltDb.Meta()
	if err != nil {
		t.Error(err)
	}
	if meta.Version != "0.1" {
		t.Error("incorrect version", meta)
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

	storageBolt, err := New(map[string]string{"filename": tmpDbFilename}, shutdownChan)
	if err != nil {
		shutdown(tmpDir)
		fmt.Fprintf(os.Stderr, "setup() storage init error %v\n", err)
		os.Exit(1)
	}
	engine = storageBolt
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
