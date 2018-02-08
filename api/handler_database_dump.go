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
	log.Trace.Printf("url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)

	type request struct {
		format string
	}

	type response struct {
	}

	switch f := r.Header.Get(HeaderAccept); f {
	case HeaderAcceptGpx:
		//
	case HeaderAcceptGeoJson:
		//
	default:
		//
	}

	w.Header().Set("Content-Type", "application/json")

	dump, err := storage.DumpAll()

	if err != nil {
		w.WriteHeader(http.StatusTeapot)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(dump))
	}
}
