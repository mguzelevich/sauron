package opengts

import (
	// "fmt"
	"net/http"
)

type Server struct {
	addr     string
	server   *http.Server
	doneChan chan bool
}

func (s Server) DoneChan() chan bool {
	return s.doneChan
}

func New(addr string) *Server {
	server := &Server{
		addr:     addr,
		doneChan: make(chan bool),
	}
	return server
}
