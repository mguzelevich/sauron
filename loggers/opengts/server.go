package opengts

import (
	"net"
	"strings"
	"time"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/loggers"
)

type Server struct {
	// server   *http.Server

	addr     string
	doneChan chan bool

	logger *loggers.Logger
}

func (s Server) DoneChan() chan bool {
	return s.doneChan
}

func (s *Server) ListenAndServe(shutdownChan chan bool) {
	/* Lets prepare a address at any address at port 10001*/
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		panic("StartUDPServer ResolveUDPAddr")
	}

	/* Now listen at selected port */
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic("StartUDPServer ListenUDP")
	}
	buf := make([]byte, 1024*4)
	closeChan := make(chan bool)
	go func() {
		for {
			select {
			case <-closeChan:
				conn.Close()
				return
			default:
			}
			if n, srcAddr, err := conn.ReadFromUDP(buf); err != nil {
				log.Error.Printf("Error [%s]\n", err)
			} else {
				timestamp := time.Now().UTC().Format(time.RFC3339)
				raw := strings.TrimSpace(string(buf[0:n]))
				log.Stderr.Printf("[%s] udp[%s]: [%q]", timestamp, srcAddr, raw)
				s.logger.Log(raw)
			}
		}
	}()

	log.Info.Printf("opengts udp logging server started [%s]\n", s.addr)
	<-shutdownChan
	// ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// server.Shutdown(ctx)
	close(closeChan)
	log.Info.Printf("opengts udp logging server gracefully stopped\n")
	s.doneChan <- true
}

func parse(raw string) (loggers.Message, error) {
	msg := udpMessage{}
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
