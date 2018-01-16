package nmea

import (
	"testing"
)

func TestNmeaToLocation(t *testing.T) {
	tsts := []struct {
		in  []string
		out Location
	}{
		{[]string{"5355.67728", "N", "5355.67728", "E"}, Location{539279546666, 539279546666}},
	}

	for idx, tst := range tsts {
		value, err := nmeaToLocation(tst.in[0], tst.in[1], tst.in[2], tst.in[3])
		if err != nil {
			t.Error("For", idx, tst, "strToLatLon failed", err)
		}
		if !tst.out.Equal(value) {
			t.Error("For", idx, tst.in, tst.out, "!=", value)
		}
	}
}
