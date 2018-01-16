package opengts

import (
	//"context"
	"net"
	"time"

	"github.com/mguzelevich/go-nmea"
	"github.com/mguzelevich/sauron/log"
	"github.com/mguzelevich/sauron/storage"
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
				raw := buf[0:n]
				log.Trace.Printf("[%s] udp: src [%s] body [%s]\n", timestamp, srcAddr, string(raw))

				if message, err := nmea.Unmarshal(raw); err != nil {
					log.Error.Printf("[%s] [%s]\n", message, err)
				} else {
					log.Debug.Printf("[%s] [%s]\n", message, err)
				}

				loc := &storage.Telemetry{}
				if err := loc.ParseUdpPacket(string(buf[0:n])); err != nil {
					//
				}
			}
		}
	}()

	log.Info.Printf("opengts udp logging server started [%s]\n", s.addr)
	<-shutdownChan
	// ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// server.Shutdown(ctx)
	close(closeChan)
	log.Debug.Printf("opengts udp logging server gracefully stopped\n")
	s.doneChan <- true
}
