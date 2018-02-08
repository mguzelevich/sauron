package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDatabaseDump(t *testing.T) {
	s := &Server{}
	r := s.router()
	ts := httptest.NewServer(r)
	defer ts.Close()
	apiUrl := "/database/dump"

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
