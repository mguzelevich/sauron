package opengts

import (
	// "fmt"
	"net/http"

	"github.com/mguzelevich/sauron/storage"
)

var (
	server   *http.Server
	doneChan chan bool
	db       *storage.Storage

	statistic string
)

func DoneChan() chan bool {
	return doneChan
}

func init() {
	doneChan = make(chan bool)
}
