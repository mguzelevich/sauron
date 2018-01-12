package opengts

import (
	// "fmt"
	"net/http"

	"github.com/mguzelevich/sauron"
)

var (
	server    *http.Server
	doneChan  chan bool
	statistic sauron.Stats

	locationsStorage *sauron.Storage
)

func DoneChan() chan bool {
	return doneChan
}

func init() {
	doneChan = make(chan bool)
}
