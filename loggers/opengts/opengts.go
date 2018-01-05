package opengts

import (
	// "fmt"
	"net/http"

	"github.com/mguzelevich/sauron"
)

var (
	server    *http.Server
	statistic sauron.Stats

	locationsStorage *sauron.Storage
)
