package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/go-sql-driver/mysql"
)

type Trip struct {
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
	fmt.Fprintf(w, "<h1>%s</h1>", "Welcome to Trip's service")
}

func allTrips(w http.ResponseWriter, r *http.Request) {

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
	router.HandleFunc("/api/v1/drivers", allTrips).Methods(
		"GET")
	// router.HandleFunc("/api/v1/availabledrivers", fetchAvailableDrivers).Methods(
	// 	"GET")
	// router.HandleFunc("/api/v1/driver/{driverid}", passenger).Methods(
	// 	"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

}
