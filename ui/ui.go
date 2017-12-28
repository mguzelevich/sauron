package ui

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
)

//go:-generate go-bindata-assetfs -prefix "../../ui/dist" -pkg main -o assetfs_ui.go ../../ui/dist/...
//go:generate rice embed-go

func Init(r *mux.Router) {
	fmt.Fprintf(os.Stderr, "ui inited\n")
	box := rice.MustFindBox("frontend/dist").HTTPBox()
	r.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(box)))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(box)))
}
