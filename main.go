package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

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

	json.NewEncoder(w).Encode(Session{Id: id, CreatedAt: time.Time{}})
}

func GetSession(w http.ResponseWriter, r *http.Request) {

	dbUrl := "postgres://skizz@localhost:5432/perimeter_development?sslmode=disable"
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	vars := mux.Vars(r)
	print(vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}

	var created_at time.Time
	err = db.QueryRow("SELECT created_at FROM sessions WHERE id = $1", id).Scan(&created_at)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(Session{Id: id, CreatedAt: created_at})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/sessions", createSession).Methods("POST")
	myRouter.HandleFunc("/sessions/{id}", GetSession).Methods("GET").Name("session")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	handleRequests()
}
