package opengts

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/mguzelevich/sauron/storage"
)

func StartUDPServer(addr string, storageDb *storage.Storage, shutdownChan chan bool) {
	db = storageDb

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
				fmt.Fprintf(os.Stderr, "[%s] udp: src [%s] body [%s]\n", timestamp, string(buf[0:n]), srcAddr)

				loc := &storage.Telemetry{}
				if err := loc.ParseUdpPacket(string(buf[0:n])); err != nil {
					//
				}
			}
		}
	}()

	fmt.Fprintf(os.Stderr, "opengts udp logging server started [%s]\n", addr)
	<-shutdownChan
}
