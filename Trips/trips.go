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

/*
setting content type to application/json and access control to allow all origins due
to cross origin resource sharing policy as request from fronted are blocked by the browser
as both the frontend server and passenger server are running on different ports but on
the same localhost.
*/
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

/*
	fetch all completed trips by passenger id from db, returns a slice of trip object
*/
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

/*
	handler for route '/api/v1/trip/{passengerid}', return a slice of all completed trips by passenger id
*/
func fetchPassengerTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, _ := getCompletedTripsByPassengerId(db, params["passengerid"]) // fetch all completed trips by passenger id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

/*
	fetching all pending trips by driver id from db, returns a slice of trip objects
*/
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

/*
	handler for route '/api/v1/getPendingTrips/{driverid}', return a slice of all pending trips by {driverid}
*/
func fetchDriverTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, err := getPendingTripsByDriverId(db, params["driverid"])
	fmt.Println(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

/*
	fetching all accepted trips by driver id from db, returns a slice of trip objects
*/
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

/*
	handler for route '/api/v1/getAcceptedTrips/{driverid}', return a slice of all accepted trips by {driverid}
*/
func fetchDriverAcceptedTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, err := getAcceptedTripsByDriverId(db, params["driverid"])
	fmt.Println(err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

/*
	send a request to driver service to call for first available driver id
*/
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

/*
	send a request to driver service to update driverStatus
*/
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

/*
	fetch all pending trips from database, return a map of trip objects
*/
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

/*
	insert a new trip into database
*/
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

/*
	handler for route '/api/v1/request/{passengerid}'
	calls driver service to get first available driver id as part of driver assignment
	then createTrip function to insert a new trip into database
	then calls driver service to update driver status to unavailable

	returns http status of 201 if successful
*/
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

/*
	update trip status of existing trip in database
*/
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

/*
	handler for route '/api/v1/changeStatus/{tripid}', returns http status of 201 if successful
*/
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

/*
	handler for route '/api/v1/endtrip/{tripid}'
	calls driver service to get update driver to available
	then updateTripStatus function to update trip status in database

	returns http status of 201 if successful
*/
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
		User:                 "root",
		Passwd:               "123",
		Net:                  "tcp",
		Addr:                 "mysql:3306",
		DBName:               "Ridely",
		AllowNativePasswords: true,
		ParseTime:            true, // to allow timestamp to parsed into time.Time object
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
