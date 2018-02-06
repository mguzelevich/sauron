package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/loggers"
)

type Server struct {
	addr     string
	server   *http.Server
	doneChan chan bool

	logger *loggers.Logger
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
	r.HandleFunc("/", s.handler).Methods("GET")
	r.HandleFunc("/log", s.logLocationHandler).Methods("POST")
	s.server.Handler = r

	go s.server.ListenAndServe()
	log.Info.Printf("custom url logger server started [%s]\n", s.addr)
	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	log.Info.Printf("custom url logger server gracefully stopped\n")
	close(s.doneChan)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Trace.Printf("\tBODY: %v\n", string(body))
	}

	var statistic string
	if out, err := json.Marshal(statistic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		log.Trace.Printf("out: %s\n", string(out))
		fmt.Fprintf(w, string(out))
	}
}

func (s *Server) logLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		timestamp := time.Now().UTC().Format(time.RFC3339)
		raw := string(body)
		log.Stderr.Printf("[%s] http: [%q]", timestamp, raw)
		w.WriteHeader(http.StatusOK)
		s.logger.Log(raw)
	}
}

func parse(raw string) (loggers.Message, error) {
	msg := message{}
	if err := msg.ParseRaw(raw); err != nil {
		return nil, err
	}
	return msg, nil
}

func New(addr string) *Server {
	server := &Server{
		addr:     addr,
		doneChan: make(chan bool),
		logger:   loggers.New(parse),
	}
	return server
}
