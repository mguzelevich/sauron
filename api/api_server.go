package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron/log"
)

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
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/database/{command}", databaseHandler).Methods("POST")
	s.server.Handler = r

	go s.server.ListenAndServe()
	log.Info.Printf("api server started [%s]\n", s.addr)
	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	log.Info.Printf("api server gracefully stopped\n")
	close(s.doneChan)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Trace.Printf("\tBODY: %v\n", string(body))
	}
	fmt.Fprintf(w, string("OK!"))
}

func New(addr string) *Server {
	server := &Server{
		addr:     addr,
		doneChan: make(chan bool),
	}
	return server
}
