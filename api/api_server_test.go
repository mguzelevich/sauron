package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://medium.com/agrea-technogies/basic-testing-patterns-in-go-d8501e360197

func TestRouter(t *testing.T) {
	r := router()
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for / is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusOK)
	}

	res, err = http.Get(ts.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /not-exists is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusNotFound)
	}
}
