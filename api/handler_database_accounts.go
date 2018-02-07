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

	w.Header().Set("Content-Type", "application/json")

	log.Trace.Printf("get list of accounts")
	if accs, err := storage.Accounts(); err != nil {
		log.Error.Printf("get list of accounts %v", err)
	} else {
		if buff, err := json.Marshal(&accountsResponse{Accounts: accs}); err != nil {
			w.WriteHeader(http.StatusTeapot)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, string(buff))
		}
	}
}
