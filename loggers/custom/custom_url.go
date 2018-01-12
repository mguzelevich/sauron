package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron"
	"github.com/mguzelevich/sauron/log"
)

var (
	server    *http.Server
	doneChan  chan bool
	statistic sauron.Stats

	locationsStorage *sauron.Storage
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Trace.Printf("\tBODY: %v\n", string(body))
	}

	if out, err := json.Marshal(statistic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Trace.Printf("out: %s\n", string(out))
		fmt.Fprintf(w, string(out))
	}
}

func logLocationHandler(w http.ResponseWriter, r *http.Request) {
	statistic.Requests++

	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		loc := &sauron.Telemetry{}
		if err := loc.ParseCustomUrl(string(body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		locationsStorage.Save(loc)
	}
}

func StartServer(addr string, storage *sauron.Storage, shutdownChan chan bool) {
	server = &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/log", logLocationHandler).Methods("POST")
	r.HandleFunc("/gts", handler).Methods("GET", "PUT")

	locationsStorage = storage

	server.Handler = r

	go server.ListenAndServe()
	log.Info.Printf("custom url logging server started [%s]\n", addr)

	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Debug.Printf("custom url logging server gracefully stopped\n")
	doneChan <- true
}

func DoneChan() chan bool {
	return doneChan
}

func init() {
	doneChan = make(chan bool)
}
