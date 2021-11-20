package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Passenger's service")
}

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
			&person.LastName, &person.EmailAddress, &person.MoblieNo); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		pMap[person.PassengerId] = person
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return pMap, nil
}

func allPassengers(w http.ResponseWriter, r *http.Request) {
	// fetch passenger map from db
	fetchedPassengerData, _ := getPassengers(db)
	fmt.Println(fetchedPassengerData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedPassengerData)

}

func insertPassenger(db *sql.DB, fN, lN, mN, eA string) {
	query := fmt.Sprintf("INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES (%s, '%s', '%s', %s)",
		fN, lN, mN, eA)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func editPassenger(db *sql.DB, fN, lN, mN, eA, id string) {
	query := fmt.Sprintf("UPDATE Passenger SET FirstName='%s', LastName='%s', MoblieNo='%s', EmailAddress='%s' WHERE PassengerId=%s",
		fN, lN, mN, eA, id)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

func passenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pMap, _ := getPassengers(db)

	if r.Method == "GET" {

		if _, ok := pMap[params["passengerid"]]; ok {
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

	if r.Header.Get("Content-type") == "application/json" {
		// POST is for creating new course
		if r.Method == "POST" {

			// read the string sent to the service
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newPassenger)

				if newPassenger.PassengerId == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " + "information " + "in JSON format"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := pMap[params["passengerid"]]; !ok {
					insertPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["passengerid"]))
				} else {
					w.WriteHeader(http.StatusConflict)
					w.Write([]byte(
						"409 - Duplicate course ID"))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information " +
					"in JSON format"))
			}
		}

		if r.Method == "PUT" {
			var newPassenger Passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newPassenger)
				// fmt.Println(newPassenger)
				newPassenger.PassengerId = params["passengerid"]

				// fmt.Println(newPassenger.PassengerId)
				if newPassenger.PassengerId == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " +
							" information " +
							"in JSON format"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := pMap[params["passengerid"]]; !ok {
					insertPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["passengerid"]))
				} else {
					// update course
					// to fix-
					editPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress, newPassenger.PassengerId)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Course updated: " +
						params["passengerid"]))
				}
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply " +
					"course information " +
					"in JSON format"))
			}
		}
	}

}

var db *sql.DB

func main() {
	// setting up db connection
	cfg := mysql.Config{
		User:   "root",
		Passwd: "123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "Ridely",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	router := mux.NewRouter()
	router.Use(commonMiddleware) //setting context to "json"
	router.HandleFunc("/api/v1/", welcome)
	router.HandleFunc("/api/v1/passengers", allPassengers).Methods(
		"GET")
	router.HandleFunc("/api/v1/passenger/{passengerid}", passenger).Methods(
		"GET", "PUT", "POST", "DELETE")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
