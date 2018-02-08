package custom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/loggers"
	"github.com/mguzelevich/sauron/storage"
	"github.com/mguzelevich/sauron/storage/engines/mmap"
)

// https://medium.com/agrea-technogies/basic-testing-patterns-in-go-d8501e360197

var (
	server     *Server
	httpserver *httptest.Server
)

func getUrlParams(p map[string]string) string {
	params := url.Values{}
	for k, v := range p {
		params.Add(k, v)
	}

	return params.Encode()
}

func TestRouter(t *testing.T) {
	url := httpserver.URL + "/"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("Status code for [%v] is wrong. Have: %d, want: %d.", url, res.StatusCode, http.StatusOK)
	}

	res, err = http.Get(httpserver.URL + "/not-exists")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Status code for /not-exists is wrong. Have: %d, want: %d.", res.StatusCode, http.StatusNotFound)
	}
}

func TestTelemetry(t *testing.T) {
	rurl := httpserver.URL + "/log"

	pp := map[string]string{
		"lat":  "53.9279421",
		"lon":  "27.6437863",
		"ts":   "2017-12-27T12:30:30.338Z",
		"s":    "0.0",
		"prov": "network",
		//"aid":  "1de1a4a0e296ef63",
		"acc": "21.795000076293945",
		"ser": "123456",
	}

	for _, id := range []string{"1de1a4a0e296ef63", "2de1a4a0e296ef63", "3de1a4a0e296ef63"} {
		pp["aid"] = id
		for lat, lon := range map[float64]float64{
			53.9279421: 27.6437863,
			53.9279422: 27.6437864,
			53.9279423: 27.6437865,
		} {
			pp["lat"] = fmt.Sprintf("%f", lat)
			pp["lon"] = fmt.Sprintf("%f", lon)
			params := getUrlParams(pp)
			log.Stderr.Printf("http: [%q]", params)
			res, err := http.Post(rurl, ContentTypeUrlEncoder, strings.NewReader(params))

			if err != nil {
				t.Fatal(err)
			}
			if res.StatusCode != http.StatusOK {
				t.Errorf("Status code for [%s] is wrong. Have: %d, want: %d.", rurl, res.StatusCode, http.StatusOK)
			}
		}
	}
	// d, err := storage.DumpAll()
	// t.Error(string(d), err)
}

/*
	curl -v -X POST --data "lat=53.9279421&lon=27.6437863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=1de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log
	curl -v -X POST --data "lat=53.9279421&lon=27.6437863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=2de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log
	curl -v -X POST --data "lat=53.9279421&lon=27.6437863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=3de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log

	curl -v -X POST --data "lat=53.9179421&lon=27.6537863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=1de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log
	curl -v -X POST --data "lat=53.9179421&lon=27.6537863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=2de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log
	curl -v -X POST --data "lat=53.9179421&lon=27.6537863&ts=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=3de1a4a0e296ef63&acc=21.795000076293945&ser=123456" localhost:8081/log

	bash -c "echo -e 'mgu/mi5s.1/\$GPRMC,083543,A,5356.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"
	bash -c "echo -e 'mgu/mi5s.2/\$GPRMC,083543,A,5356.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"
	bash -c "echo -e 'mgu/mi5s.3/\$GPRMC,083543,A,5356.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"

	bash -c "echo -e 'mgu/mi5s.1/\$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"
	bash -c "echo -e 'mgu/mi5s.2/\$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"
	bash -c "echo -e 'mgu/mi5s.3/\$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,050118,,*2E' > /dev/udp/127.0.0.1/8082"

	curl -X POST --data '{}' localhost:8080/database/dump | jq '.'
*/

func setup() {
	server := &Server{
		doneChan: make(chan bool),
		logger:   loggers.New(parse),
	}

	storage.Init(mmap.New, map[string]string{}, make(chan bool))

	httpserver = httptest.NewServer(server.router())
}

func shutdown() {
	httpserver.Close()
}

func TestMain(m *testing.M) {
	// http://cs-guy.com/blog/2015/01/test-main/
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
