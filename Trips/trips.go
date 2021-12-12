/*
	Trip service:
		handle request trip by passenger
		get all trips by passenger id
		get all trips by driver id
		change trip status
*/
package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/go-sql-driver/mysql"
)

type Trip struct {
	TripId            int       `json:"tripid"`
	PassengerId       string    `json:"passengerid"`
	DriverId          string    `json:"driverid"`
	PickUpPostalCode  string    `json:"pickuppostalcode"`
	DropOffPostalCode string    `json:"dropoffpostalcode"`
	TripStatus        int       `json:"tripstatus"` // 0 - pending, 1 - accepted, 2 - completed
	DateOfTrip        time.Time `json:"dateoftrip"`
}

// middleware for setting header to json only
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		next.ServeHTTP(w, r)
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Trip's service")
}

func getCompletedTripsByPassengerId(db *sql.DB, id string) ([]Trip, error) {
	var tArr []Trip

	rows, err := db.Query("SELECT * FROM Trips WHERE PassengerId=? AND TripStatus=2 ORDER BY DateofTrip DESC", id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripId, &trip.PassengerId, &trip.DriverId, &trip.PickUpPostalCode, &trip.DropOffPostalCode, &trip.TripStatus, &trip.DateOfTrip); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tArr = append(tArr, trip)
		fmt.Println(tArr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return tArr, nil

}

func fetchPassengerTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, _ := getCompletedTripsByPassengerId(db, params["passengerid"])
	// fmt.Print(fetchedTripData, err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

func getPendingTripsByDriverId(db *sql.DB, id string) ([]Trip, error) {
	var tArr []Trip

	rows, err := db.Query("SELECT * FROM Trips WHERE DriverId=? AND TripStatus=0", id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripId, &trip.PassengerId, &trip.DriverId, &trip.PickUpPostalCode, &trip.DropOffPostalCode, &trip.TripStatus, &trip.DateOfTrip); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tArr = append(tArr, trip)
		fmt.Println(tArr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return tArr, nil

}

func fetchDriverTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, err := getPendingTripsByDriverId(db, params["driverid"])
	fmt.Println(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

func getAcceptedTripsByDriverId(db *sql.DB, id string) ([]Trip, error) {
	var tArr []Trip

	rows, err := db.Query("SELECT * FROM Trips WHERE DriverId=? AND TripStatus=1", id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripId, &trip.PassengerId, &trip.DriverId, &trip.PickUpPostalCode, &trip.DropOffPostalCode, &trip.TripStatus, &trip.DateOfTrip); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tArr = append(tArr, trip)
		fmt.Println(tArr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return tArr, nil

}

func fetchDriverAcceptedTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, err := getAcceptedTripsByDriverId(db, params["driverid"])
	fmt.Println(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

// request for trip
func fetchFirstAvailableDriver() (driverId string, err error) {
	// get driverid
	const baseURL = "http://localhost:5001/api/v1/availabledrivers"
	resp, err := http.Get(baseURL)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	var drivers []string
	json.Unmarshal(data, &drivers)
	fmt.Println(drivers)
	return drivers[0], nil
}

func changeDriverStatus(driverid string, driverStatus int) (err error) {
	// get all drivers
	jsonValue, _ := json.Marshal(map[string]int{"driverstatus": driverStatus})
	baseURL := "http://localhost:5001/api/v1/driver/changeStatus/" + driverid
	request, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)

	if err != nil {
		log.Fatal(err)
		return err
	} else {
		fmt.Println(resp.StatusCode)
	}
	defer request.Body.Close()
	return nil
}

func getAllPendingTrips(db *sql.DB) (map[int]Trip, error) {
	tMap := make(map[int]Trip)

	rows, err := db.Query("SELECT * FROM Trips WHERE TripStatus=1")
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var trip Trip
		if err := rows.Scan(&trip.TripId, &trip.DriverId, &trip.PassengerId, &trip.PickUpPostalCode, &trip.DropOffPostalCode, &trip.TripStatus, &trip.DateOfTrip); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tMap[trip.TripId] = trip
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return tMap, nil
}

func createTrip(db *sql.DB, tripDetails Trip) (err error) {
	// insert into db
	stmt, err := db.Prepare("INSERT INTO Trips (PassengerId, DriverId, PickUpPostalCode, DropOffPostalCode, TripStatus ) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tripDetails.PassengerId, tripDetails.DriverId, tripDetails.PickUpPostalCode, tripDetails.DropOffPostalCode, tripDetails.TripStatus)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func requestTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "POST" {
		var newTrip Trip
		fmt.Println(r.Body)
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newTrip)
			driverId, _ := fetchFirstAvailableDriver() // get first available driver
			newTrip.DriverId = driverId
			newTrip.PassengerId = params["passengerid"]
			newTrip.TripStatus = 0
			createTrip(db, newTrip)
			changeDriverStatus(driverId, 0) // change to unavailable
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Trip added for passenger id: " +
				params["passengerid"]))
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply trip information " +
				"in JSON format"))
		}
	}

}

func updateTripStatus(db *sql.DB, tripId int, status int) (err error) {
	stmt, err := db.Prepare("UPDATE Trips SET TripStatus=? WHERE TripId=?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, tripId)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func changeTripStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "POST" {
		var newTrip Trip
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newTrip)
			tripId, _ := strconv.Atoi(params["tripid"])
			fmt.Println(newTrip.TripStatus, newTrip.TripId)
			updateTripStatus(db, tripId, newTrip.TripStatus)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Trip Status Updated: " +
				params["tripid"]))
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply trip information " +
				"in JSON format"))
		}

	}
}

func endTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	allTrips, _ := getAllPendingTrips(db)
	if r.Method == "POST" {
		var newTrip Trip
		fmt.Println(r.Body)
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &newTrip)
			tripId, _ := strconv.Atoi(params["tripid"])
			driverId := allTrips[tripId].DriverId
			fmt.Println(tripId, driverId)
			changeDriverStatus(driverId, 1) // change to available
			updateTripStatus(db, tripId, newTrip.TripStatus)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Trip added for passenger id: " +
				params["passengerid"]))
		} else {
			w.WriteHeader(
				http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply trip information " +
				"in JSON format"))
		}
	}

}

var db *sql.DB

func main() {
	// setting up db connection
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "123",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "Ridely",
		ParseTime: true, // to allow timestamp to parsed into time.Time object
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
	router.HandleFunc("/api/v1/trip/{passengerid}", fetchPassengerTrips).Methods(
		"GET")
	router.HandleFunc("/api/v1/getPendingTrips/{driverid}", fetchDriverTrips).Methods(
		"GET")
	router.HandleFunc("/api/v1/getAcceptedTrips/{driverid}", fetchDriverAcceptedTrips).Methods(
		"GET")
	router.HandleFunc("/api/v1/request/{passengerid}", requestTrip).Methods(
		"POST")
	router.HandleFunc("/api/v1/changeStatus/{tripid}", changeTripStatus).Methods(
		"POST")
	router.HandleFunc("/api/v1/endtrip/{tripid}", endTrip).Methods(
		"POST")

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

}
