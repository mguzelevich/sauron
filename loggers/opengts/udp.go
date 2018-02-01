package opengts

import (
	//"context"
	"net"
	"strings"
	"time"

	"github.com/mguzelevich/go-log"
)

func (s *Server) ListenAndServeUdp(shutdownChan chan bool) {
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

				msg := udpMessage{}
				if err := msg.ParseUdpPacket(raw); err != nil {
					log.Error.Printf("parse packet [%q] error [%s]", raw, err)
					continue
				}
				s.processMsgChan <- msg
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
