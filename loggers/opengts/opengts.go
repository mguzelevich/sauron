package opengts

import (
	// "fmt"
	"net/http"

	"github.com/mguzelevich/sauron/storage"
)

var (
	server *http.Server
	db     *storage.Storage

	statistic string
)
