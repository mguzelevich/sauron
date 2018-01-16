package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDatabaseTelemetry(t *testing.T) {
	r := router()
	ts := httptest.NewServer(r)
	defer ts.Close()
	apiUrl := "/database/telemetry"

	res, err := http.Get(ts.URL + apiUrl)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status code for %s is wrong. Have: %d, want: %d.", apiUrl, res.StatusCode, http.StatusMethodNotAllowed)
	}

	/*
		res, err = http.Post(ts.URL+apiUrl, "text/plain", nil)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status code for %s is wrong. Have: %d, want: %d.", apiUrl, res.StatusCode, http.StatusOK)
		}
	*/
}
