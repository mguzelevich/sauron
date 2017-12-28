package main

import (
	//	"log"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron"
	"github.com/mguzelevich/sauron/ui"
)

var (
	debug bool

	host string
	port int

	locChan chan *sauron.Location

	statistic sauron.Stats
	storage   *sauron.Storage
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")

	flag.StringVar(&host, "h", "localhost", "host")
	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "p", 8080, "port")
	flag.IntVar(&port, "port", 8080, "port")

	// flag.BoolVar(&ui, "ui", false, "web ui")

	// flag.StringVar(&host, "h", "localhost", "host")
	// flag.StringVar(&host, "host", "localhost", "host")
	// flag.IntVar(&port, "p", 8080, "port")
	// flag.IntVar(&port, "port", 8080, "port")
}

func handler(w http.ResponseWriter, r *http.Request) {
	statistic.Requests++

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
	fmt.Fprintf(os.Stderr, "location loop started\n")
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

func walk(r *mux.Router) {
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			//return err
		}
		qt, err := route.GetQueriesTemplates()
		if err != nil {
			//fmt.Fprintf(os.Stderr, "err: %v\n", err)
			//return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			//return err
		}
		// qr will contain a list of regular expressions with the same semantics as GetPathRegexp,
		// just applied to the Queries pairs instead, e.g., 'Queries("surname", "{surname}") will return
		// {"^surname=(?P<v0>.*)$}. Where each combined query pair will have an entry in the list.
		qr, err := route.GetQueriesRegexp()
		if err != nil {
			//fmt.Fprintf(os.Stderr, "err: %v\n", err)
			//return err
		}
		m, err := route.GetMethods()
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %v\n", err)
			//return err
		}
		fmt.Fprintf(os.Stderr, "> m: %v\tqt: %v qr: %v t: %v p: %v\n", strings.Join(m, ","), strings.Join(qt, ","), strings.Join(qr, ","), t, p)
		return nil
	})
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/log", logLocationHandler).Methods("POST")
	r.HandleFunc("/gts", handler).Methods("GET", "PUT")
	ui.Init(r)

	locChan = make(chan *sauron.Location)
	storage = sauron.NewStorage()
	go locLoop()

	if debug {
		walk(r)
	}
	go http.ListenAndServe(fmt.Sprintf(":%d", port), r)

	doneChan := make(chan bool)
	<-doneChan
}
