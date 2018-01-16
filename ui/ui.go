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

type Server struct {
	addr     string
	server   *http.Server
	doneChan chan bool
}

func (s Server) DoneChan() chan bool {
	return s.doneChan
}

func (s *Server) ListenAndServe(shutdownChan chan bool) {
	s.server = &http.Server{
		Addr:           s.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	r := mux.NewRouter()
	box := rice.MustFindBox("frontend/dist").HTTPBox()
	r.Handle("/", http.StripPrefix("/", http.FileServer(box)))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(box)))
	s.server.Handler = r

	go s.server.ListenAndServe()
	log.Info.Printf("ui server started [%s]\n", s.addr)
	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	log.Info.Printf("ui server gracefully stopped\n")
	close(s.doneChan)
}

func New(addr string) *Server {
	server := &Server{
		addr:     addr,
		doneChan: make(chan bool),
	}
	return server
}
