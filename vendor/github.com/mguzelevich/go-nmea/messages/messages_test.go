package messages

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

func compareStringsLists(a, b []string) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func compareLists(a, b []int) bool {
	if &a == &b {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

func TestSplit_pos(t *testing.T) {
	in := []byte("$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E")
	expectedOut := []string{
		"$GPRMC",
		"083543",
		"A",
		"5355.67728", "N",
		"2738.62654", "E",
		"0.000000",
		"0.000000",
		"050118",
		"",
		"",
		"*2E",
	}
	out, err := split(in)
	if err != nil {
		t.Error("For", "TestSplit_pos", string(in), "split failed")
	}
	if !compareStringsLists(expectedOut, out) {
		t.Error("For", "TestSplit_pos", string(in), "expected", expectedOut, "got", out)
	}
}

func TestUnmarshal_pos(t *testing.T) {
	in := []byte("$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E")

	msg, err := ParseMessage(in)
	if err != nil {
		t.Error("For", "TestUnmarshal_pos", string(in), "unmarshal failed", err)
		return
	}
	rmcMsg := msg.(*Rmc)
	if rmcMsg.Status != "A" {
		t.Error("For", "TestUnmarshal_pos [", string(in), "] unmarshal (", msg, ") failed", err)
	}
}

func TestUnmarshal_unknown_type(t *testing.T) {
	in := []byte("$GPGSV,3,2,09,13,27,292,26,20,19,317,25,27,,,21,28,34,171,35*49")

	msg, err := ParseMessage(in)
	if err != nil {
		t.Error("For", "TestUnmarshal_unknown_type", string(in), "unmarshal ", msg, "failed", err)
	}
}
