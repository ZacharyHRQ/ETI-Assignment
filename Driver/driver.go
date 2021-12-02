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

type Driver struct {
	DriverId             string `json:"driverid"`
	FirstName            string `json:"firstname"`
	LastName             string `json:"lastname"`
	MoblieNo             string `json:"moblieno"`
	EmailAddress         string `json:"emailaddress"`
	CarLicenseNo         string `json:"carlicenseno"`
	IdentificationNumber string `json:"identificationnumber"`
	DriverStatus         int    `json:"driverstatus"` // 1 means available , 0 means unavailable
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
			&person.LastName, &person.MoblieNo, &person.EmailAddress, &person.CarLicenseNo, &person.IdentificationNumber, &person.DriverStatus); err != nil {
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
	// fetch driver map from db
	fetchedDriverData, _ := getDrivers(db)
	// fmt.Println(fetchedPassengerData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedDriverData)

}

func getDriverById(w http.ResponseWriter, r *http.Request) {
	// fetch driver map from db
	params := mux.Vars(r)
	fetchedDriverData, _ := getDrivers(db)
	fmt.Println(fetchedDriverData)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fetchedDriverData[params["driverid"]])

}

func insertDriver(db *sql.DB, fN, lN, mN, eA, cA, iN string) {
	stmt, err := db.Prepare("INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress, CarLicenseNo, IdentificationNumber) VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(fN, lN, mN, eA, cA, iN)
	if err != nil {
		log.Fatal(err)
	}
}

func editDriver(db *sql.DB, fN, lN, mN, eA, cA, id string) {
	stmt, err := db.Prepare("UPDATE Driver SET FirstName=?, LastName=?, MoblieNo=?, EmailAddress=?, CarLicenseNo=? WHERE DriverId=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(fN, lN, mN, eA, cA, id)
	if err != nil {
		log.Fatal(err)
	}
}

func createDriver(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422- Please supply driver information in JSON format"))
		}
		var driver Driver
		json.Unmarshal(body, &driver)
		insertDriver(db, driver.FirstName, driver.LastName, driver.MoblieNo, driver.EmailAddress, driver.CarLicenseNo, driver.IdentificationNumber)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(driver)
	}
}

func updateDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422- Please supply driver information in JSON format"))
		}
		var driver Driver
		json.Unmarshal(body, &driver)
		editDriver(db, driver.FirstName, driver.LastName, driver.MoblieNo, driver.EmailAddress, driver.CarLicenseNo, params["driverid"])
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(driver)
	}
}

func getDriver(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pMap, _ := getDrivers(db)

	if r.Method == "GET" {

		if _, ok := pMap[params["driverid"]]; ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(pMap[params["driverid"]])
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
			&person.LastName, &person.MoblieNo, &person.EmailAddress, &person.CarLicenseNo, &person.IdentificationNumber, &person.DriverStatus); err != nil {
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
	driverIds := make([]string, 0)
	for k := range fetchedDriverData {
		driverIds = append(driverIds, k)
	}
	fmt.Println(driverIds)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(driverIds)
}

func updateDriverStatus(db *sql.DB, driverId string, status int) (err error) {
	stmt, err := db.Prepare("UPDATE Driver SET DriverStatus=? WHERE DriverId=?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, driverId)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func changeDriverStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if r.Header.Get("Content-type") == "application/json" {
		if r.Method == "POST" {
			var newDriver Driver
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				json.Unmarshal(reqBody, &newDriver)
				updateDriverStatus(db, params["driverid"], newDriver.DriverStatus)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Driver Status Updated: " +
					params["driverid"]))
			} else {
				w.WriteHeader(
					http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply driver information " +
					"in JSON format"))
			}
		}

	}
}

func getDriverIds(w http.ResponseWriter, r *http.Request) {
	dMap, _ := getDrivers(db)
	driverIds := make([]string, 0)
	for k := range dMap {
		fmt.Println(k)
		driverIds = append(driverIds, k)
	}
	sort.Strings(driverIds)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(driverIds)
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
	router.HandleFunc("/api/v1/driver/{driverid}", getDriverById).Methods(
		"GET")
	router.HandleFunc("/api/v1/availabledrivers", fetchAvailableDrivers).Methods(
		"GET")
	router.HandleFunc("/api/v1/fetchAllIds", getDriverIds).Methods(
		"GET")
	router.HandleFunc("/api/v1/driver/{driverid}", getDriver).Methods(
		"GET", "DELETE")
	router.HandleFunc("/api/v1/driver/createDriver", createDriver).Methods(
		"POST")
	router.HandleFunc("/api/v1/driver/updateDriver/{driverid}", updateDriver).Methods(
		"POST")
	router.HandleFunc("/api/v1/driver/changeStatus/{driverid}", changeDriverStatus).Methods(
		"POST")

	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))

}
