package nmea

import (
	"fmt"
	// "github.com/mguzelevich/go-nmea/messages"
)

// packet type -> payload fields count
var register map[string]int

/*
http://www.gpsinformation.org/dale/nmea.htm
*/
func init() {
	register = map[string]int{
		"$GPAAM": -1, // Waypoint Arrival Alarm
		"$GPALM": -1, // Almanac data
		"$GPAPA": -1, // Auto Pilot A sentence
		"$GPAPB": -1, // Auto Pilot B sentence
		"$GPBOD": -1, // Bearing Origin to Destination
		"$GPBWC": -1, // Bearing using Great Circle route
		"$GPDTM": -1, // Datum being used.
		"$GPGGA": -1, // Fix information
		"$GPGLL": -1, // Lat/Lon data
		"$GPGRS": -1, // GPS Range Residuals
		"$GPGSA": -1, // Overall Satellite data
		"$GPGST": -1, // GPS Pseudorange Noise Statistics
		"$GPGSV": -1, // Detailed Satellite data
		"$GPMSK": -1, // send control for a beacon receiver
		"$GPMSS": -1, // Beacon receiver status information.
		"$GPRMA": -1, // recommended Loran data
		"$GPRMB": -1, // recommended navigation data for gps
		"$GPRMC": 11, // recommended minimum data for gps
		"$GPRTE": -1, // route message
		"$GPTRF": -1, // Transit Fix Data
		"$GPSTN": -1, // Multiple Data ID
		"$GPVBW": -1, // dual Ground / Water Spped
		"$GPVTG": -1, // Vector track an Speed over the Ground
		"$GPWCV": -1, // Waypoint closure velocity (Velocity Made Good)
		"$GPWPL": -1, // Waypoint Location information
		"$GPXTC": -1, // cross track error
		"$GPXTE": -1, // measured cross track error
		"$GPZTG": -1, // Zulu (UTC) time and time to go (to destination)
		"$GPZDA": -1, // Date and Time
	}
}

type NmeaPacket struct {
	packetType string
	fields     []string
	crc        string
}

func NewPacket(data []byte) (*NmeaPacket, error) {
	p := &NmeaPacket{}

	if fields, err := split(data); err != nil {
		return nil, err
	} else {
		p.packetType = fields[0]
		p.fields = fields[1 : len(fields)-1]
		p.crc = fields[len(fields)-1]
	}

	cnt, ok := register[p.packetType]
	if !ok {
		return nil, fmt.Errorf("unknown message type [%s]", p.packetType)
	}
	if cnt != len(p.fields) {
		return nil, fmt.Errorf("expected %d fields for [%s], got %d ", cnt, p.packetType, len(p.fields))
	}
	if err := p.validateCrc(); err != nil {
		return nil, err
	}

	return p, nil
}

func (n *NmeaPacket) validateCrc() error {
	if n.crc[0] != '*' {
		return fmt.Errorf("checksum error (first symbol must be '*'")
	}
	return nil
}

func (n *NmeaPacket) AsAam() (*Aam, error) {
	return &Aam{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsAlm() (*Alm, error) {
	return &Alm{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsApa() (*Apa, error) {
	return &Apa{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsApb() (*Apb, error) {
	return &Apb{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsBod() (*Bod, error) {
	return &Bod{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsBwc() (*Bwc, error) {
	return &Bwc{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsDtm() (*Dtm, error) {
	return &Dtm{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGga() (*Gga, error) {
	return &Gga{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGll() (*Gll, error) {
	return &Gll{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGrs() (*Grs, error) {
	return &Grs{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGsa() (*Gsa, error) {
	return &Gsa{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGst() (*Gst, error) {
	return &Gst{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsGsv() (*Gsv, error) {
	return &Gsv{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsMsk() (*Msk, error) {
	return &Msk{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsMss() (*Mss, error) {
	return &Mss{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsRma() (*Rma, error) {
	return &Rma{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsRmb() (*Rmb, error) {
	return &Rmb{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsRmc() (*Rmc, error) {
	p := &Rmc{}
	if err := p.build(n); err != nil {
		return nil, err
	}
	return p, nil
}

func (n *NmeaPacket) AsRte() (*Rte, error) {
	return &Rte{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsTrf() (*Trf, error) {
	return &Trf{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsStn() (*Stn, error) {
	return &Stn{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsVbw() (*Vbw, error) {
	return &Vbw{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsVtg() (*Vtg, error) {
	return &Vtg{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsWcv() (*Wcv, error) {
	return &Wcv{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsWpl() (*Wpl, error) {
	return &Wpl{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsXtc() (*Xtc, error) {
	return &Xtc{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsXte() (*Xte, error) {
	return &Xte{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsZtg() (*Ztg, error) {
	return &Ztg{}, fmt.Errorf("temporary unsupported. in development")
}

func (n *NmeaPacket) AsZda() (*Zda, error) {
	return &Zda{}, fmt.Errorf("temporary unsupported. in development")
}
