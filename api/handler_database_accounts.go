package api

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"

	"github.com/mguzelevich/go.log"
	"github.com/mguzelevich/sauron/storage"
)

// type scanRequest struct {
// 	Sources []string `json:"sources"`
// }

type accountsResponse struct {
	Accounts []*storage.Account `json:"accounts"`
}

func databaseAccountsHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")

	log.Trace.Printf("get list of accounts")
	// var req scanRequest
	// if err := json.Unmarshal(body, &req); err != nil {
	// 	log.Error.Printf("incorrect req")
	// }

	// if len(req.Sources) == 0 {
	// 	req.Sources = []string{"."}
	// }

	// taskId := "" // engine.TaskScanFs(req.Sources)
	if accs, err := storage.Accounts(); err != nil {
		log.Error.Printf("get list of accounts %v", err)
	} else {
		// for _, a := range accs {
		// 	log.Trace.Printf("%v", a)
		// }
		if buff, err := json.Marshal(&accountsResponse{Accounts: accs}); err != nil {
			w.WriteHeader(http.StatusTeapot)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(buff))
		}
	}
}
