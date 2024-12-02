package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("hello")

	r := mux.NewRouter()

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:4242",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
