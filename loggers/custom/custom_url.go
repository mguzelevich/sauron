package custom

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/go-log"
	"github.com/mguzelevich/sauron/storage"
)

type Server struct {
	addr     string
	server   *http.Server
	doneChan chan bool

	processTelemetryChan chan string
}

func (s Server) DoneChan() chan bool {
	return s.doneChan
}

func (s *Server) processLoop() {
	for raw := range s.processTelemetryChan {
		msg := &message{}
		if err := msg.ParseCustomUrl(raw); err != nil {
			log.Error.Printf("parse packet [%q] error [%s]", raw, err)
			continue
		}
		/*
			получить хеш девайса из телеметрии
			true:
				получить инстанс девайса из стораджа
				получить инстанс аккаунта
			false:
				создать инстанс девайса
				приааттачить к анонимусу (создать новый аккаунт)

			получить инстанс девайса из аккаунта
			записать телеметрию в девайс
		*/

		/*
			получить хеш девайса из телеметрии
			получить инстанс девайса из аккаунта
			записать телеметрию в девайс
		*/

		account := &storage.Account{}

		device, err := storage.GetDevice(&storage.Device{Id: msg.device().Hash()})
		if err != nil {
			switch err {
			case storage.ErrEntityNotFound:
				device, _ = account.CreateDevice(device)
			default:
				panic(err)
			}
		}

		account, err = storage.ReadAccount(&storage.Account{Id: device.UserId})
		if err != nil {
			switch err {
			case storage.ErrEntityNotFound:
				account, _ = storage.CreateAccount(account)
			default:
				panic(err)
			}
		}

		//account, _ := storage.ReadAccount(&storage.Account{Id: device.UserId})
		device.AddTelemetry(msg.telemetry())
	}
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
		log.Stdout.Printf("[%s] http: [%q]", timestamp, raw)
		w.WriteHeader(http.StatusOK)
		s.processTelemetryChan <- raw
	}
}

func New(addr string) *Server {
	server := &Server{
		addr:                 addr,
		doneChan:             make(chan bool),
		processTelemetryChan: make(chan string),
	}
	go server.processLoop()
	return server
}
