package api

import (
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	"net/http"

	"github.com/mguzelevich/go.log"

	"github.com/mguzelevich/sauron/storage"
)

func databaseDumpHandler(w http.ResponseWriter, r *http.Request) {
	log.Trace.Printf("url: %s %s %d\n", r.Method, r.RequestURI, r.ContentLength)

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")

	log.Trace.Printf("dump")
	storage.Dump()

	// if buff, err := json.Marshal(&accauntsResponse{Id: "taskId"}); err != nil {
	// 	w.WriteHeader(http.StatusTeapot)
	// } else {
	// 	w.WriteHeader(http.StatusOK)
	// 	fmt.Fprintf(w, string(buff))
	// }

}
