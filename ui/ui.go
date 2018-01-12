package ui

import (
	"context"
	"net/http"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron/log"
)

//go:-generate go-bindata-assetfs -prefix "../../ui/dist" -pkg main -o assetfs_ui.go ../../ui/dist/...
//go:generate rice embed-go

var (
	server   *http.Server
	doneChan chan bool
)

func StartServer(addr string, shutdownChan chan bool) {
	r := mux.NewRouter()

	server = &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	box := rice.MustFindBox("frontend/dist").HTTPBox()
	r.Handle("/", http.StripPrefix("/", http.FileServer(box)))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(box)))

	server.Handler = r

	go server.ListenAndServe()
	log.Info.Printf("ui server started [%s]\n", addr)
	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	server.Shutdown(ctx)
	log.Debug.Printf("ui server gracefully stopped\n")
	doneChan <- true
}

func DoneChan() chan bool {
	return doneChan
}

func init() {
	doneChan = make(chan bool)
}
