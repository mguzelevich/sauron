package nmea

import ()

/*
http://www.gpsinformation.org/dale/nmea.htm
*/
func init() {
	/*
		register = map[string]CreateMessageFunc{
			"$GPAAM": nil,       // Waypoint Arrival Alarm
			"$GPALM": nil,       // Almanac data
			"$GPAPA": nil,       // Auto Pilot A sentence
			"$GPAPB": nil,       // Auto Pilot B sentence
			"$GPBOD": nil,       // Bearing Origin to Destination
			"$GPBWC": nil,       // Bearing using Great Circle route
			"$GPDTM": nil,       // Datum being used.
			"$GPGGA": nil,       // Fix information
			"$GPGLL": nil,       // Lat/Lon data
			"$GPGRS": nil,       // GPS Range Residuals
			"$GPGSA": nil,       // Overall Satellite data
			"$GPGST": nil,       // GPS Pseudorange Noise Statistics
			"$GPGSV": nil,       // Detailed Satellite data
			"$GPMSK": nil,       // send control for a beacon receiver
			"$GPMSS": nil,       // Beacon receiver status information.
			"$GPRMA": nil,       // recommended Loran data
			"$GPRMB": nil,       // recommended navigation data for gps
			"$GPRMC": createRmc, // recommended minimum data for gps
			"$GPRTE": nil,       // route message
			"$GPTRF": nil,       // Transit Fix Data
			"$GPSTN": nil,       // Multiple Data ID
			"$GPVBW": nil,       // dual Ground / Water Spped
			"$GPVTG": nil,       // Vector track an Speed over the Ground
			"$GPWCV": nil,       // Waypoint closure velocity (Velocity Made Good)
			"$GPWPL": nil,       // Waypoint Location information
			"$GPXTC": nil,       // cross track error
			"$GPXTE": nil,       // measured cross track error
			"$GPZTG": nil,       // Zulu (UTC) time and time to go (to destination)
			"$GPZDA": nil,       // Date and Time
		}
	*/
}

type Aam struct{}

type Alm struct{}

type Apa struct{}

type Apb struct{}

type Bod struct{}

type Bwc struct{}

type Dtm struct{}

type Gga struct{}

type Gll struct{}

type Grs struct{}

type Gsa struct{}

type Gst struct{}

type Gsv struct{}

type Msk struct{}

type Mss struct{}

type Rma struct{}

type Rmb struct{}

type Rte struct{}

type Trf struct{}

type Stn struct{}

type Vbw struct{}

type Vtg struct{}

type Wcv struct{}

type Wpl struct{}

type Xtc struct{}

type Xte struct{}

type Ztg struct{}

type Zda struct{}
