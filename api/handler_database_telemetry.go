package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mguzelevich/go.log"
	//	"github.com/mguzelevich/sauron/storage"
)

func databaseTelemetryHandler(w http.ResponseWriter, r *http.Request) {

	type request struct {
		format string
	}

	// application/json
	// application/gpx+xml
	// application/vnd.geo+json

	type response struct {
	}

	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	log.Trace.Printf("get list of telemetry")
	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		log.Error.Printf("incorrect req")
	}

	if buff, err := json.Marshal(&response{}); err != nil {
		w.WriteHeader(http.StatusTeapot)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(buff))
	}

}
