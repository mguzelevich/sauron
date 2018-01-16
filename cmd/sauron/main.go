package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/mguzelevich/sauron/api"
	"github.com/mguzelevich/sauron/log"
	"github.com/mguzelevich/sauron/loggers/custom"
	"github.com/mguzelevich/sauron/loggers/opengts"
	"github.com/mguzelevich/sauron/storage"
	"github.com/mguzelevich/sauron/ui"
)

type shutdownable interface {
	DoneChan() chan bool
}

var (
	debug bool

	apiServerAddr string
	database      string

	httpServerAddr string
	udpServerAddr  string

	uiServerAddr string
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")

	flag.StringVar(&apiServerAddr, "api", "localhost:8080", "http logger server address")

	flag.StringVar(&httpServerAddr, "http", "localhost:8081", "http logger server address")
	flag.StringVar(&udpServerAddr, "udp", ":8082", "udp logger server address")
	flag.StringVar(&uiServerAddr, "ui", "localhost:8083", "ui server address")

	flag.StringVar(&database, "db", "/tmp/sauron.db", "database file")
}

func shutdown(shutdownChan chan bool, services ...shutdownable) {
	close(shutdownChan)

	chans := []chan bool{}
	for _, s := range services {
		chans = append(chans, s.DoneChan())
	}

	for {
		timeout := time.After(10 * time.Second)
		select {
		case <-chans[0]:
			chans[0] = nil
		case <-chans[1]:
			chans[1] = nil
		case <-chans[2]:
			chans[2] = nil
		case <-chans[3]:
			chans[3] = nil
		case <-chans[4]:
			chans[4] = nil
		case <-timeout:
			log.Warning.Printf("shutdowned by timer")
			return
		default:
			allDone := true
			for _, ch := range chans {
				allDone = allDone && ch == nil
			}
			if allDone {
				return
			}
		}
	}
}

func main() {
	flag.Parse()

	log.InitLoggers(&log.Logger{
		os.Stdout, // ioutil.Discard,
		os.Stdout, // ioutil.Discard,
		os.Stdout,
		os.Stdout,
		os.Stderr,
	})

	shutdownChan := make(chan bool)

	db, err := storage.Init(database, shutdownChan)
	if err != nil {
		log.Error.Printf("engine init error %s", err)
		os.Exit(1)
	}

	apiServer := api.New(apiServerAddr)
	go apiServer.ListenAndServe(shutdownChan)

	customServer := custom.New(httpServerAddr)
	go customServer.ListenAndServe(shutdownChan)

	opengtsServer := opengts.New(udpServerAddr)
	go opengtsServer.ListenAndServeUdp(shutdownChan)

	uiServer := ui.New(uiServerAddr)
	go uiServer.ListenAndServe(shutdownChan)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	shutdown(shutdownChan, db, apiServer, customServer, opengtsServer, uiServer)
}
