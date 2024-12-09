package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var TrDB TreeDB

func main() {
	log.Println("the tree daemon is up!")

	db_loc := os.Getenv("DB_LOCATION")
	if db_loc == "" {
		log.Println("warning: using development DB location. ideally set DB_LOCATION")
		db_loc = "./tree.db"
	}

	db, err := sql.Open("sqlite3", db_loc)
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

	TrDB = Prepare(conn)

	ReadDirectory()

	r := mux.NewRouter()

	r.Path("/upload").Methods("POST").HandlerFunc(UploadFile)
	r.Path("/events").HandlerFunc(Monitor)
	r.Path("/skip").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prog := ActiveProgram.name
		SkipProgram <- true

		w.Write([]byte("skipped " + prog))
	})

	r.Path("/name").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !CoolMode() {
			http.Error(w, "not working rn", http.StatusTeapot)
			return
		}

		name, err := GetName(r)
		if err != nil {
			http.Error(w, "name not found", http.StatusNotFound)
			log.Println("error finding name", r.RemoteAddr, err)
			return
		}

		w.Write([]byte(name))
	})

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:4242",
		// WriteTimeout: 15 * time.Second,
		// ReadTimeout:  15 * time.Second,
	}

	SetupLEDStrip()

	// listen to the internal socket
	go ListenSocket()

	// start queue
	go ManageQueue()

	log.Fatal(srv.ListenAndServe())
}
