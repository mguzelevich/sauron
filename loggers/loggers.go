package loggers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron"
)

var (
	server *http.Server

	locChan chan *sauron.Location

	storage   *sauron.Storage
	statistic sauron.Stats
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(os.Stderr, "\tBODY: %v\n", string(body))
	}

	if out, err := json.Marshal(statistic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(os.Stderr, "out: %s\n", string(out))
		fmt.Fprintf(w, string(out))
	}
}

func locLoop() {
	for loc := range locChan {
		storage.Save(loc)
	}
}

func logLocationHandler(w http.ResponseWriter, r *http.Request) {
	statistic.Requests++

	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if values, err := url.ParseQuery(string(body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			loc := &sauron.Location{
				Lat:       values.Get("lat"),
				Lon:       values.Get("lon"),
				Sat:       values.Get("sat"),
				Desc:      values.Get("desc"),
				Alt:       values.Get("alt"),
				Acc:       values.Get("acc"),
				Dir:       values.Get("dir"),
				Prov:      values.Get("prov"),
				Spd:       values.Get("spd"),
				Time:      values.Get("time"),
				Battery:   values.Get("battery"),
				AndroidId: values.Get("androidId"),
				Serial:    values.Get("serial"),
				Activity:  values.Get("activity"),
				Epoch:     values.Get("epoch"),
			}
			w.WriteHeader(http.StatusOK)
			locChan <- loc
		}
	}
}

func StartServer(addr string, shutdownChan chan bool) {
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

	locChan = make(chan *sauron.Location)
	storage = sauron.NewStorage()
	go locLoop()

	server.Handler = r

	fmt.Fprintf(os.Stderr, "logging server started [%s]\n", addr)
	go server.ListenAndServe()
	<-shutdownChan
}
