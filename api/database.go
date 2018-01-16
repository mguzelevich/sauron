package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/mguzelevich/sauron/log"
	//	"github.com/mguzelevich/sauron/storage"
)

type scanRequest struct {
	Sources []string `json:"sources"`
}

type scanResponse struct {
	Id string `json:"id"`
}

func databaseHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch mux.Vars(r)["command"] {
	case "scan":
		scanHandler(w, r, body)
	case "show":
		log.Trace.Printf("show database stats")
	default:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(`{"status":"OK"}`))
	}

}

func scanHandler(w http.ResponseWriter, r *http.Request, body []byte) {
	log.Trace.Printf("scan database")
	var req scanRequest
	if err := json.Unmarshal(body, &req); err != nil {
		log.Error.Printf("incorrect req")
	}

	if len(req.Sources) == 0 {
		req.Sources = []string{"."}
	}

	taskId := "" // engine.TaskScanFs(req.Sources)

	if buff, err := json.Marshal(&scanResponse{Id: taskId}); err != nil {
		w.WriteHeader(http.StatusTeapot)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(buff))
	}

}
