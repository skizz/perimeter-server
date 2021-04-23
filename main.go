package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Session struct {
	Id   int
	Time int64
}

var sessions = []Session{
	Session{Id: 1, Time: 999},
	Session{Id: 2, Time: 1000},
}

func allSessions(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(sessions)
}

func createSession(w http.ResponseWriter, r *http.Request) {

	dbUrl := "postgres://skizz@localhost:5432/perimeter_development?sslmode=disable"
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var id int
	err = db.QueryRow("INSERT INTO sessions DEFAULT VALUES RETURNING id").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ID = %d\n", id)

	json.NewEncoder(w).Encode(Session{Id: id, Time: 0})
}

func GetSession(w http.ResponseWriter, r *http.Request) {

	dbUrl := "postgres://skizz@localhost:5432/perimeter_development?sslmode=disable"
	db, err := sql.Open("postgres", dbUrl)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	var time int64
	err = db.QueryRow("SELECT id, time FROM sessions WHERE id = ?", id).Scan(&time)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(Session{Id: id, Time: time})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/sessions", allSessions).Methods("GET")
	myRouter.HandleFunc("/sessions", createSession).Methods("POST")
	myRouter.HandleFunc("/sessions/{id}", GetSession).Methods("POST").Name("session")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
