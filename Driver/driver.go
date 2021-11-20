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

type Driver struct {
	DriverId     string `json:"driverid"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	MoblieNo     string `json:"moblieno"`
	EmailAddress string `json:"emailaddress"`
	CarLicenseNo string `json:"carlicenseno"`
	DriverStatus int    `json:"driverstatus"` // 1 means available , 0 means unavailable
}

// middleware for setting header to json only
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Driver's service")
}

func getDrivers(db *sql.DB) (map[string]Driver, error) {
	pMap := make(map[string]Driver)

	rows, err := db.Query("SELECT * FROM Driver")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var person Driver
		if err := rows.Scan(&person.DriverId, &person.FirstName,
			&person.LastName, &person.EmailAddress, &person.MoblieNo, &person.CarLicenseNo, &person.DriverStatus); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		pMap[person.DriverId] = person
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return pMap, nil
}

func allDrivers(w http.ResponseWriter, r *http.Request) {
	// fetch passenger map from db
	fetchedPassengerData, _ := getDrivers(db)
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
	pMap, _ := getDrivers(db)

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
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.DriverId == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply course " + "information " + "in JSON format"))
					return
				}

				// check if course exists; add only if
				// course does not exist
				if _, ok := pMap[params["driverid"]]; !ok {
					insertPassenger(db, newDriver.FirstName, newDriver.LastName, newDriver.MoblieNo, newDriver.EmailAddress)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["driverid"]))
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
			var newPassenger Driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &newPassenger)
				newPassenger.DriverId = params["driverid"]
				if newPassenger.DriverId == "" {
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
				if _, ok := pMap[params["driverid"]]; !ok {
					//insertPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("201 - Course added: " +
						params["driverid"]))
				} else {
					// update course
					// to fix-
					//editPassenger(db, newPassenger.FirstName, newPassenger.LastName, newPassenger.MoblieNo, newPassenger.EmailAddress, newPassenger.PassengerId)
					w.WriteHeader(http.StatusAccepted)
					w.Write([]byte("202 - Course updated: " +
						params["driverid"]))
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

func getAvailableDrivers(db *sql.DB) (map[string]Driver, error) {
	pMap := make(map[string]Driver)

	rows, err := db.Query("SELECT * FROM Driver WHERE DriverStatus=1")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var person Driver
		if err := rows.Scan(&person.DriverId, &person.FirstName,
			&person.LastName, &person.EmailAddress, &person.MoblieNo, &person.CarLicenseNo, &person.DriverStatus); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		pMap[person.DriverId] = person
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return pMap, nil
}

func fetchAvailableDrivers(w http.ResponseWriter, r *http.Request) {
	fetchedDriverData, _ := getAvailableDrivers(db)
	fmt.Println(fetchedDriverData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedDriverData)
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
	router.HandleFunc("/api/v1/drivers", allDrivers).Methods(
		"GET")
	router.HandleFunc("/api/v1/availabledrivers", fetchAvailableDrivers).Methods(
		"GET")
	router.HandleFunc("/api/v1/driver/{driverid}", passenger).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))

}
