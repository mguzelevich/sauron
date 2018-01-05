package loggers

import (
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "net/url"
	// "os"
	// "time"

	// "github.com/gorilla/mux"

	"github.com/mguzelevich/sauron"
)

type LoggerServer interface {
	StartServer(addr string, storage *sauron.Storage, shutdownChan chan bool)
}
