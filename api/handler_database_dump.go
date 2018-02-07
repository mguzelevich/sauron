package api

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/storage"
)

func databaseDumpHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		format string
	}

	// application/json
	// application/gpx+xml
	// application/vnd.geo+json

	type response struct {
	}

	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	w.Header().Set("Content-Type", "application/json")

	dump, err := storage.DumpAll()

	if err != nil {
		w.WriteHeader(http.StatusTeapot)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(dump))
	}
}
