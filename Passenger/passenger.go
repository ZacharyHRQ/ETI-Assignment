package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
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

func getPassengerById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["passengerid"]
	fmt.Println(id)
	fetchedPassengerData, _ := getPassengers(db)
	fmt.Println(fetchedPassengerData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedPassengerData[id])
}

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
}

func createPassenger(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create passenger")
	pMap, _ := getPassengers(db)
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
			newPassenger.PassengerId = strconv.Itoa(len(pMap) + 1)
			fmt.Println(newPassenger)
			json.NewEncoder(w).Encode(newPassenger)
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply passenger information " +
				"in JSON format"))
		}
	}

}

func updatePassenger(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pMap, _ := getPassengers(db)
	if r.Method == "PUT" {
		var newPassenger Passenger
		reqBody, err := ioutil.ReadAll(r.Body)
		if err == nil {
			json.Unmarshal(reqBody, &newPassenger)

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
	router.HandleFunc("/api/v1/passenger/{passengerid}", getPassengerById).Methods(
		"GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", passenger).Methods(
		"GET", "Delete")
	router.HandleFunc("/api/v1/passenger/createPassenger", createPassenger).Methods(
		"POST")
	router.HandleFunc("/api/v1/passenger/updatePassenger/{passengerid}", updatePassenger).Methods(
		"PUT")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))

}
