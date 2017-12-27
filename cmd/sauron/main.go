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
)

type stats struct {
	Requests int `json:"requests"`
	Errors   int `json:"errors"`
}

type location struct {
	Lat       string `json:"lat"`                 // lat=%LAT
	Lon       string `json:"lon"`                 // lon=%LON
	Sat       string `json:"sat,omitempty"`       // sat=%SAT
	Desc      string `json:"desc,omitempty"`      // desc=%DESC
	Alt       string `json:"alt,omitempty"`       // alt=%ALT
	Acc       string `json:"acc,omitempty"`       // acc=%ACC
	Dir       string `json:"dir,omitempty"`       // dir=%DIR
	Prov      string `json:"prov,omitempty"`      // prov=%PROV
	Spd       string `json:"spd,omitempty"`       // spd=%SPD
	Time      string `json:"time,omitempty"`      // time=%TIME
	Battery   string `json:"battery,omitempty"`   // battery=%BATT
	AndroidId string `json:"androidId,omitempty"` // androidId=%AID
	Serial    string `json:"serial,omitempty"`    // serial=%SER
	Activity  string `json:"activity,omitempty"`  // activity=%ACT
	Epoch     string `json:"epoch,omitempty"`     // epoch=%TIMESTAMP"
}

var (
	static bool
	host   string
	port   int

	locChan chan location

	statistic stats
)

func init() {
	flag.BoolVar(&static, "static", false, "static serve")

	flag.StringVar(&host, "h", "localhost", "host")
	flag.StringVar(&host, "host", "localhost", "host")

	flag.IntVar(&port, "p", 8080, "port")
	flag.IntVar(&port, "port", 8080, "port")
}

func execute() {
	if static {
		init_static("/files/", "/tmp")
	} else {
		http.HandleFunc("/", handler)
	}
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func init_static(urlPrefix string, root string) {
	http.Handle(urlPrefix, http.StripPrefix(urlPrefix, http.FileServer(http.Dir(root))))
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
		if out, err := json.Marshal(loc); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			fmt.Fprintf(os.Stderr, "location: %s\n", string(out))
			// fmt.Fprintf(w, string(out))
		}
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
			loc := location{
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
			return err
		}
		qt, err := route.GetQueriesTemplates()
		if err != nil {
			return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		// qr will contain a list of regular expressions with the same semantics as GetPathRegexp,
		// just applied to the Queries pairs instead, e.g., 'Queries("surname", "{surname}") will return
		// {"^surname=(?P<v0>.*)$}. Where each combined query pair will have an entry in the list.
		qr, err := route.GetQueriesRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(m, ","), strings.Join(qt, ","), strings.Join(qr, ","), t, p)
		return nil
	})
}

func main() {
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/log", logLocationHandler).Methods("POST")
	r.HandleFunc("/gts", handler).Methods("GET", "PUT")
	http.Handle("/", r)

	locChan = make(chan location)
	go locLoop()

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
