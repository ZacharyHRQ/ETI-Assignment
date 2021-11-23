package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	TripStatus        int       `json:"tripstatus"`
	DateOfTrip        time.Time `json:"dateoftrip"`
}

// middleware for setting header to json only
func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Trip's service")
}

func getTripsById(db *sql.DB, id string) ([]Trip, error) {
	var tArr []Trip

	rows, err := db.Query("SELECT * FROM Trips WHERE PassengerId=?", id)
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

func fetchTrips(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, _ := getTripsById(db, params["passengerid"])
	// fmt.Print(fetchedTripData, err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
}

// request for trip

func requestTrip(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fetchedTripData, _ := getTripsById(db, params["passengerid"])
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedTripData)
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
	router.HandleFunc("/api/v1/trips/{passengerid}", fetchTrips).Methods(
		"GET")
	router.HandleFunc("/api/v1/trip/request/{passengerid}", requestTrip).Methods(
		"GET")
	// router.HandleFunc("/api/v1/availabledrivers", fetchAvailableDrivers).Methods(
	// 	"GET")
	// router.HandleFunc("/api/v1/driver/{driverid}", passenger).Methods(
	// 	"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

}
