package opengts

import (
	// "fmt"
	"net/http"

	// "github.com/mguzelevich/go-log"

	"github.com/mguzelevich/sauron/storage"
)

type Server struct {
	addr     string
	server   *http.Server
	doneChan chan bool

	processMsgChan chan udpMessage
}

func (s Server) DoneChan() chan bool {
	return s.doneChan
}

func (s *Server) processLoop() {
	for msg := range s.processMsgChan {
		telemetry := msg.telemetry()
		account, _ := storage.ReadAccount(&storage.Account{})
		device, _ := account.ReadDevice(&storage.Device{})
		device.AddTelemetry(telemetry)
	}
}

func New(addr string) *Server {
	server := &Server{
		addr:           addr,
		doneChan:       make(chan bool),
		processMsgChan: make(chan udpMessage),
	}
	go server.processLoop()
	return server
}
