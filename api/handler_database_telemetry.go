package api

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"

	"github.com/mguzelevich/go-log"
	//	"github.com/mguzelevich/sauron/storage"
)

func databaseTelemetryHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")

	log.Trace.Printf("get list of telemetry")
	// var req scanRequest
	// if err := json.Unmarshal(body, &req); err != nil {
	// 	log.Error.Printf("incorrect req")
	// }

	// if len(req.Sources) == 0 {
	// 	req.Sources = []string{"."}
	// }

	// taskId := "" // engine.TaskScanFs(req.Sources)
	// storage.GetTelemetry("")

	if buff, err := json.Marshal(&accountsResponse{}); err != nil {
		w.WriteHeader(http.StatusTeapot)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(buff))
	}

}
