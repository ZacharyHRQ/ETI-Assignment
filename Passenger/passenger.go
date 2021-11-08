package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Passenger struct {
	passengerId  int    `json:"id"`
	firstName    string `json:"firstname"`
	lastName     string `json:"lastname"`
	moblieNo     string `json:"moblieNo"`
	emailAddress string `json:"emailAddress"`
}

var passengers map[int]Passenger

func GetRecords(db *sql.DB) {
	results, err := db.Query("Select * FROM Ridely.Passenger")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// map this type to the record in the table
		var person Passenger
		err = results.Scan(&person.passengerId, &person.firstName, &person.lastName, &person.moblieNo, &person.emailAddress)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println(person.passengerId)
	}

}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Passenger's service")
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "List of all courses")
	// returns the key/value pairs in the query string as a map object
	kv := r.URL.Query()

	for k, v := range kv {
		fmt.Println(k, v) // print out the key/value pair
	}

	json.NewEncoder(w).Encode(passengers)
}

func passenger(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// open connection to database on port 3306
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Ridely")

	// handle error
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened")
	}

	GetRecords(db)
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", welcome)
	router.HandleFunc("/api/v1/passengers", allPassengers)
	router.HandleFunc("/api/v1/passengers/{passengerId}", passenger).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
