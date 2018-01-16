package opengts

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron/log"
	"github.com/mguzelevich/sauron/storage"
)

func (s *Server) ListenAndServe(shutdownChan chan bool) {
	s.server = &http.Server{
		Addr:           s.addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handler).Methods("GET")
	r.HandleFunc("/log", logLocationHandler).Methods("POST")
	r.HandleFunc("/gts", handler).Methods("GET", "PUT")
	s.server.Handler = r

	go s.server.ListenAndServe()
	log.Info.Printf("opengts http server started [%s]\n", s.addr)
	<-shutdownChan
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.server.Shutdown(ctx)
	log.Info.Printf("opengts http server gracefully stopped\n")
	close(s.doneChan)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(os.Stderr, "url: %s %s %d %v\n", r.Method, r.RequestURI, r.ContentLength, r.Header)
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(os.Stderr, "\tBODY: %v\n", string(body))
	}

	var statistic string
	if out, err := json.Marshal(statistic); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(os.Stderr, "out: %s\n", string(out))
		fmt.Fprintf(w, string(out))
	}
}

func logLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if body, err := ioutil.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		loc := &storage.Telemetry{}
		if err := loc.ParseCustomUrl(string(body)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		storage.Save(loc)
	}
}
