# go.nmea [![GoDoc](https://godoc.org/github.com/mguzelevich/go.nmea?status.svg)](http://godoc.org/github.com/mguzelevich/go.nmea) [![Build Status](https://travis-ci.org/mguzelevich/go.nmea.svg?branch=master)](https://travis-ci.org/mguzelevich/go.nmea)

nmea protocol parser


# usage

```
import (
	"fmt"

	"github.com/mguzelevich/go.nmea"
)

func main() {
	packet := "$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E"
	if msg, err := nmea.Unmarshal([]byte(packet)); err != nil {
		fmt.Printf("nmea.Unmarshal error %s", err)
	} else {
		rmcMsg := msg.(*messages.Rmc)

		fmt.Printf("nmea message %v", rmcMsg)
	}
}

```