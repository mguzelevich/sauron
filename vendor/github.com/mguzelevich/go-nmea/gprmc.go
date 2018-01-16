package nmea

import (
	"fmt"
	"strconv"
	"time"
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
	Timestamp         *time.Time
	Status            string
	Location          *Location
	Speed             float64
	Direction         float64
	MagneticVariation float64
	checksum          string
}

func (r Rmc) check(p *NmeaPacket) error {
	if p.packetType != "$GPRMC" {
		return fmt.Errorf("incorrect packet type [%s]. $GPRMC expected", p.packetType)
	}
	return nil
}

func (r *Rmc) build(p *NmeaPacket) error {
	if err := r.check(p); err != nil {
		return err
	}
	r.Timestamp, _ = toTime(p.fields[8], p.fields[0])
	r.Status = p.fields[1]
	r.Location, _ = nmeaToLocation(p.fields[2], p.fields[3], p.fields[4], p.fields[5])

	speed, _ := strconv.ParseFloat(p.fields[6], 64)
	r.Speed = speed

	direction, _ := strconv.ParseFloat(p.fields[7], 64)
	r.Direction = direction

	r.MagneticVariation, _ = strToLatLon(p.fields[8], p.fields[9])
	return nil
}
