package nmea

import (
	"testing"
	"time"
)

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

func TestSplit(t *testing.T) {
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
		"", "",
		"*2E",
	}
	out, err := split(in)
	if err != nil {
		t.Error("For", string(in), "split failed")
	}
	if !compareStringsLists(expectedOut, out) {
		t.Error("For", string(in), "expected", expectedOut, "got", out)
	}
}

func TestToTime(t *testing.T) {
	tsts := []struct {
		date string
		time string
		out  time.Time
	}{
		{"050118", "083543", time.Date(2018, time.Month(1), 5, 8, 35, 43, 0, time.UTC)},
	}
	for idx, tst := range tsts {
		ts, err := toTime(tst.date, tst.time)
		if err != nil {
			t.Error("For", idx, tst, "toTime failed", err)
		}
		if !ts.Equal(tst.out) {
			t.Error("For", idx, tst, "toTime failed", err)
		}
	}
}

func TestStrToLatLon(t *testing.T) {
	tsts := []struct {
		loc  string
		sign string
		out  Location
	}{
		{"5355.67728", "N", Location{}},
	}

	for idx, tst := range tsts {
		value, err := strToLatLon(tst.loc, tst.sign)
		if err != nil {
			t.Error("For", idx, tst, "strToLatLon failed", err)
		}
		if value == 0 {
			t.Error("For", idx, tst, "strToLatLon failed", err)
		}
	}
}
