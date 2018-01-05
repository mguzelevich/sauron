package opengts

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mguzelevich/go-nmea"
	"github.com/mguzelevich/sauron"
)

func StartUDPServer(addr string, storage *sauron.Storage, shutdownChan chan bool) {
	locationsStorage = storage

	/* Lets prepare a address at any address at port 10001*/
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic("StartUDPServer ResolveUDPAddr")
	}

	/* Now listen at selected port */
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic("StartUDPServer ListenUDP")
	}
	defer conn.Close()

	buf := make([]byte, 1024*4)

	go func() {
		for {
			if n, srcAddr, err := conn.ReadFromUDP(buf); err != nil {
				fmt.Fprintf(os.Stderr, "Error [%s]\n", err)
			} else {
				timestamp := time.Now().UTC().Format(time.RFC3339)
				raw := buf[0:n]
				fmt.Fprintf(os.Stderr, "[%s] udp: src [%s] body [%s]\n", timestamp, srcAddr, string(raw))

				if message, err := nmea.Unmarshal(raw); err != nil {
					fmt.Fprintf(os.Stderr, "[%s] [%s]\n", message, err)
				} else {
					fmt.Fprintf(os.Stderr, "[%s] [%s]\n", message, err)
				}

				loc := &sauron.Telemetry{}
				if err := loc.ParseUdpPacket(string(buf[0:n])); err != nil {
					//
				}
			}
		}
	}()

	fmt.Fprintf(os.Stderr, "opengts udp logging server started [%s]\n", addr)
	<-shutdownChan
}
