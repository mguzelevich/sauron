package ui

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
)

//go:-generate go-bindata-assetfs -prefix "../../ui/dist" -pkg main -o assetfs_ui.go ../../ui/dist/...
//go:generate rice embed-go

var (
	server *http.Server
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

	fmt.Fprintf(os.Stderr, "ui server started [%s]\n", addr)
	go server.ListenAndServe()
	<-shutdownChan
}
