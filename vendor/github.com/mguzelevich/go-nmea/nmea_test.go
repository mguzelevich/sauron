package nmea

import (
	// "fmt"
	// "strings"
	"testing"
)

var tests = []string{
	"$GPGGA,195433.00,5357.23210,N,02741.19650,E,1,06,1.21,234.0,M,24.7,M,,*58",
	"$GPGLL,5357.23210,N,02741.19650,E,195433.00,A,A*6D",
	"$GPGSA,A,3,07,28,30,08,13,20,,,,,,,2.07,1.21,1.68*0C",
	"$GPGSV,3,1,09,04,,,23,05,,,20,07,61,073,36,08,18,076,30*72",
	"$GPGSV,3,2,09,13,27,292,26,20,19,317,25,27,,,21,28,34,171,35*49",
	"$GPGSV,3,3,09,30,82,218,34*45",
	"$GPRMC,195433.00,A,5357.23210,N,02741.19650,E,0.017,,051017,,,A*70",
	"$GPVTG,,T,,M,0.017,N,0.031,K,A*27",
}

func TestNewPacket_pos(t *testing.T) {
	in := []byte("$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E")

	if _, err := NewPacket(in); err != nil {
		t.Error("For", string(in), "test failed", err)
		return
	}
}

func TestNewPacket_unknown_type(t *testing.T) {
	in := []byte("$GPAAA,3,2,09,13,27,292,26,20,19,317,25,27,,,21,28,34,171,35*49")

	if _, err := NewPacket(in); err == nil {
		t.Error("For", string(in), "test failed", err)
		return
	}
}

func TestNewPacket_incorrect_fields_cnt(t *testing.T) {
	in := []byte("$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,*2E")

	if _, err := NewPacket(in); err == nil {
		t.Error("For", string(in), "test failed", err)
		return
	}
}

func TestAsRmc_pos(t *testing.T) {
	in := []byte("$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E")

	p, err := NewPacket(in)
	if err != nil {
		t.Error("For", string(in), "test failed", err)
		return
	}

	if rmc, err := p.AsRmc(); err != nil {
		t.Error("For", string(in), "test failed", err)
	} else {
		if rmc.Status != "A" {
		}
	}
}
