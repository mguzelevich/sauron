package main

import (
	"flag"
	stllog "log"
	"os"
	"os/signal"
	"time"

	"github.com/mguzelevich/go-log"

	"github.com/mguzelevich/sauron/api"
	"github.com/mguzelevich/sauron/loggers/custom"
	"github.com/mguzelevich/sauron/loggers/opengts"
	"github.com/mguzelevich/sauron/storage"
	"github.com/mguzelevich/sauron/storage/engines/bolt"
	"github.com/mguzelevich/sauron/storage/engines/mmap"
	"github.com/mguzelevich/sauron/ui"
	"github.com/mguzelevich/sauron/version"
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

	memory bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "debug mode")

	flag.StringVar(&apiServerAddr, "api", "localhost:8080", "http logger server address")

	flag.StringVar(&httpServerAddr, "http", "localhost:8081", "http logger server address")
	flag.StringVar(&udpServerAddr, "udp", ":8082", "udp logger server address")
	flag.StringVar(&uiServerAddr, "ui", "localhost:8083", "ui server address")

	flag.StringVar(&database, "db", "/tmp/sauron.db", "database file")
	flag.BoolVar(&memory, "memory", false, "in-memory mode")
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
			log.Warning.Printf("service shutdowned by timer")
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
	log.Stdout = stllog.New(os.Stdout, "", 0)

	log.Info.Printf(
		"starting the service... <commit: %s, build time: %s, release: %s>",
		version.Commit, version.BuildTime, version.Release,
	)

	shutdownChan := make(chan bool)
	services := []shutdownable{}

	var f storage.NewFunc
	params := map[string]string{}
	if memory {
		f = mmap.New
	} else {
		f = bolt.New
		params["filename"] = database
	}

	db := storage.Init(f, params, shutdownChan)
	services = append(services, db)

	apiServer := api.New(apiServerAddr)
	go apiServer.ListenAndServe(shutdownChan)

	customServer := custom.New(httpServerAddr)
	go customServer.ListenAndServe(shutdownChan)

	opengtsServer := opengts.New(udpServerAddr)
	go opengtsServer.ListenAndServeUdp(shutdownChan)

	uiServer := ui.New(uiServerAddr)
	go uiServer.ListenAndServe(shutdownChan)

	services = append(services, []shutdownable{apiServer, customServer, opengtsServer, uiServer}...)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	shutdown(shutdownChan, services...)
}
