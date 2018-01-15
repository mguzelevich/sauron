package main

import (
	//	"log"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron/log"
	"github.com/mguzelevich/sauron/loggers/custom"
	"github.com/mguzelevich/sauron/loggers/opengts"
	"github.com/mguzelevich/sauron/storage"
	"github.com/mguzelevich/sauron/ui"
)

var (
	debug bool

	httpServerAddr string
	udpServerAddr  string

	uiServerAddr string

	database string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")

	flag.StringVar(&httpServerAddr, "http", "localhost:8080", "http logger server address")
	flag.StringVar(&udpServerAddr, "udp", ":8022", "udp logger server address")
	flag.StringVar(&uiServerAddr, "ui", "localhost:8081", "ui server address")

	flag.StringVar(&database, "db", "/tmp/librarian.db", "database file")
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

	shutdownChan := make(chan bool)
	doneChan := make(chan bool)

	db, err := storage.Init(database, shutdownChan)
	if err != nil {
		log.Error.Printf("engine init error %s", err)
		os.Exit(1)
	}

	go custom.StartServer(httpServerAddr, db, shutdownChan)
	go opengts.StartUDPServer(udpServerAddr, db, shutdownChan)

	go ui.StartServer(uiServerAddr, shutdownChan)

	<-doneChan
	close(shutdownChan)
}
