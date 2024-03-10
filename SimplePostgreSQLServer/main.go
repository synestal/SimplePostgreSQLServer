package main

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func newRouter() *mux.Router {

	r := mux.NewRouter()

	staticFileDirectory := http.Dir("./static/")

	staticFileServer := http.FileServer(staticFileDirectory)
	staticFileHandler := http.StripPrefix("/", staticFileServer)

	r.Handle("/", staticFileHandler).Methods("GET")

	r.HandleFunc("/person/get", getPersonHandler).Methods("GET")
	r.HandleFunc("/person/create", createPersonHandler).Methods("POST")
	r.HandleFunc("/person/delete", deletePersonHandler).Methods("POST")
	r.HandleFunc("/database/auth", Authorise).Methods("POST")

	return r
}

func main() {
	r := newRouter()
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
