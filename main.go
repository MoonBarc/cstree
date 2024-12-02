package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.Println("the tree daemon is up!")

	db, err := sql.Open("sqlite3", "./tree.db")
	if err != nil {
		log.Fatalln("failed to open DB connection", err)
	}

	defer db.Close()

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalln("failed to connect to DB", err)
	}

	RunMigrations(conn)
	log.Println("doen")

	defer conn.Close()

	// treedb := Prepare(conn)

	directory := ReadDirectory()

	r := mux.NewRouter()

	r.Path("/name").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name, err := GetName(directory, r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte(name))
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:4242",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
