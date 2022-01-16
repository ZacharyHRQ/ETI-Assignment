/*
	Passenger service:
		Create passenger
		Update passenger
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"

	"github.com/go-sql-driver/mysql"
)

type Passenger struct {
	PassengerId  string `json:"passengerid"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	MoblieNo     string `json:"moblieno"`
	EmailAddress string `json:"emailaddress"`
}

/*
setting content type to application/json and access control to allow all origins due
to cross origin resource sharing policy as request from fronted are blocked by the browser
as both the frontend server and passenger server are running on different ports but on
the same localhost.
*/

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "Welcome to Passenger's service")
}

/*
	gets all passengers from the db and returns a map of passenger ids and passenger objects
*/
func getPassengers(db *sql.DB) (map[string]Passenger, error) {
	pMap := make(map[string]Passenger)

	rows, err := db.Query("SELECT * FROM Passenger")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var person Passenger
		if err := rows.Scan(&person.PassengerId, &person.FirstName,
			&person.LastName, &person.MoblieNo, &person.EmailAddress); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		pMap[person.PassengerId] = person
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return pMap, nil
}

/*
	handler for route '/api/v1/passengers', return a map of all the passengers
*/
func allPassengers(w http.ResponseWriter, r *http.Request) {
	// fetch passenger map from db
	fetchedPassengerData, _ := getPassengers(db)
	fmt.Println(fetchedPassengerData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedPassengerData)

}

/*
	handler for route '/api/v1/passenger/{passengerid}', returns the passenger by the {passengerid}
*/
func getPassengerById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["passengerid"]
	fmt.Println(id)
	fetchedPassengerData, _ := getPassengers(db)
	fmt.Println(fetchedPassengerData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedPassengerData[id])
}

/*
	handler for route '/api/v1/passengersid', returns an slice of all passengers id
*/
func getPassengerIds(w http.ResponseWriter, r *http.Request) {
	pMap, _ := getPassengers(db)
	ids := make([]string, len(pMap))
	i := 0
	for k := range pMap {
		ids[i] = k
		i++
	}
	sort.Strings(ids)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ids)
}

/*
	inserting a new passenger into the db
*/
func insertPassenger(db *sql.DB, fN, lN, mN, eA string) {
	// insert passenger into db
	stmt, err := db.Prepare("INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES (?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(fN, lN, mN, eA)
	if err != nil {
		log.Fatal(err)
	}

}

/*
	updating an existing passenger in db
*/
func editPassenger(db *sql.DB, fN, lN, mN, eA, id string) {

	stmt, err := db.Prepare("UPDATE Passenger SET FirstName=?, LastName=?, MoblieNo=?, EmailAddress=? WHERE PassengerId=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(fN, lN, mN, eA, id)
	if err != nil {
		log.Fatal(err)
	}
}

/*
	handler for route '/api/v1/passengers/{passengerid}', return passenger object based on {passengerid}
*/
func passenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pMap, _ := getPassengers(db)

	if r.Method == "GET" {

		if _, ok := pMap[params["passengerid"]]; ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(pMap[params["passengerid"]])
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No course found"))
		}
	}

	if r.Method == "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Deletion is not allowed"))
	}
}

/*
	handler for route '/api/v1/passenger/createPassenger', returns the created passenger object and status code of the request
*/
func createPassenger(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// read the string sent to the service
		var newPassenger Passenger
		rb, err := ioutil.ReadAll(r.Body)
		fmt.Println(string(rb), err)

		if err == nil {
			// convert JSON to object
			json.Unmarshal(rb, &newPassenger)
			fmt.Println(newPassenger)
			insertPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newPassenger)
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply passenger information " +
				"in JSON format"))
		}
	}

}

/*
	handler for route '/api/v1/passenger/updatePassenger/{passengerid}', returns status code of the request
*/
func updatePassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "POST" {
		var newPassenger Passenger
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &newPassenger)
			editPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress, params["passengerid"])
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Passenger updated: " +
				params["passengerid"]))

		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply " +
				"Passenger information " +
				"in JSON format"))
		}
	}
}

var db *sql.DB

func main() {
	// setting up db connection
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "123",
		Net:                  "tcp",
		Addr:                 "mysql:3306",
		DBName:               "Ridely",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping() // check if the connection is alive
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!") // if connection is alive, print "Connected!"

	router := mux.NewRouter()
	router.Use(commonMiddleware) //setting context to "json" n cors error

	// routes
	router.HandleFunc("/", welcome)
	router.HandleFunc("/api/v1/passengers", allPassengers).Methods(
		"GET")
	router.HandleFunc("/api/v1/passenger/{passengerid}", getPassengerById).Methods(
		"GET")
	router.HandleFunc("/api/v1/passengersid/", getPassengerIds).Methods(
		"GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods(
		"GET", "Delete")
	router.HandleFunc("/api/v1/passenger/createPassenger", createPassenger).Methods(
		"POST")
	router.HandleFunc("/api/v1/passenger/updatePassenger/{passengerid}", updatePassenger).Methods(
		"POST")

	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
