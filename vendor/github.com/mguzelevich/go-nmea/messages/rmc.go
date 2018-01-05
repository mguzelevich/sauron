package messages

import (
	"fmt"
)

/*
RMC

NMEA has its own version of essential gps pvt (position, velocity, time) data.
It is called RMC, The Recommended Minimum, which will look similar to:

	$GPRMC,083559.00,A,4717.11437,N,00833.91522,E,0.004,77.52,091202,,,A*57

	$GPRMC,123519,A,4807.038,N,01131.000,E,022.4,084.4,230394,003.1,W*6A

	Where:
	     RMC          Recommended Minimum sentence C
	     123519       Fix taken at 12:35:19 UTC
	     A            Status A=active or V=Void.
	     4807.038,N   Latitude 48 deg 07.038' N
	     01131.000,E  Longitude 11 deg 31.000' E
	     022.4        Speed over the ground in knots
	     084.4        Track angle in degrees True
	     230394       Date - 23rd of March 1994
	     003.1,W      Magnetic Variation
	     *6A          The checksum data, always begins with *
*/

type Rmc struct {
	Time              string
	Status            string
	Latitude          string
	Longitude         string
	Speed             string
	Direction         string
	Date              string
	MagneticVariation string
	checksum          string
}

func (r Rmc) Marshal() ([]byte, error) {
	return nil, nil
}

func (r Rmc) Unmarshal(fields []string) error {
	r.Status = fields[2]
	return nil
}

func checkRmc(fields []string) error {
	if len(fields) != 13 {
		return fmt.Errorf("incorrect fields count")
	}
	return nil
}

func createRmc(fields []string) (NmeaMessage, error) {
	if err := checkRmc(fields); err != nil {
		return nil, err
	}
	msg := &Rmc{}
	err := msg.Unmarshal(fields)
	return msg, err
}
