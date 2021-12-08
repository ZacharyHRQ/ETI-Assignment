# ETI-Assignment

## Design consideration

### Architecture of microservices

There are 3 microservices (Passenger, Driver and Trips). The Passenenger service will handle the creation and updating of passenger information. The Driver service will handle the creation and updating of driver information. The Trip service will handle the creation of trips and retrieving of trip information. These 3 services will have an MySQL database connected to fetch and store data.

### Separation of concerns for passenger and driver users

It is best practices to allow each portion of the program to handle one portion of application logic. In this case, the choice to separate the passenger and driver was to allow each microservice to handle one set of logic. In addition, this gives better isolation of logic which aids in maintainability as logic is associated to a particual domain.

### Trip service

This service will used to store trip data and will act as a link between passenger and driver. Passenger will interact with this service to request for a trip as a result the trip service will call the driver service for a list of available driver to assign for this trip. The driver will then call the trip service to confirm that the driver has accepted the trip.

### Why use reactjs

Reactjs is a component based web framework which makes User interface (UI) resuable which is a good practice of Do not repeat yourself (DRY).

### Data strucuture

> Passenger struct

| Name          | Description                    | type   |
| ------------- | ------------------------------ | ------ |
| PassengerId   | unique indentifer of passenger | string |
| First Name    | first name of passenger        | string |
| Last Name     | last name of passenger         | string |
| Moblie Number | moblie number of passenger     | string |
| email address | email address of passenger     | string |

> Driver struct

| Name                 | Description                                               | type   |
| -------------------- | --------------------------------------------------------- | ------ |
| DriverId             | unique indentifer of Driver                               | string |
| First Name           | first name of Driver                                      | string |
| Last Name            | last name of Driver                                       | string |
| Moblie Number        | moblie number of Driver                                   | string |
| email address        | email address of Driver                                   | string |
| CarLicenseNo         | Car License Number of Driver                              | string |
| IdentificationNumber | Identification Number of Driver                           | string |
| DriverStatus         | status of driver , 1 means available, 0 means unavailable | int    |

> Trip struct

| Name              | Description                                         | type      |
| ----------------- | --------------------------------------------------- | --------- |
| TripId            | unique indentifer of trip                           | string    |
| PassengerId       | unique indentifer of passenger                      | string    |
| DriverId          | unique indentifer of driver                         | string    |
| PickUpPostalCode  | pick up location postal code                        | string    |
| DropOffPostalCode | drop off location postal code                       | string    |
| TripStatus        | status of trip , 1 means completed, 0 means pending | int       |
| DateOfTrip        | timestamp of trip                                   | time.Time |

## Architecture Diagram

![Architecture Diagram](https://github.com/ZacharyHRQ/ETI-Assignment/blob/main/Architecture.png)

### Features

> Passenger service

#### Create Passenger

![Architecture Diagram](https://github.com/ZacharyHRQ/ETI-Assignment/blob/main/Architecture.png)

## Instructions for setting up and running

### Prerequites

- [NodeJs](https://nodejs.org/en/)
- Golang

### Ports used

- Passenger service - 5000
- Driver service - 5001
- Trips service - 5002

- Reactjs Frontend - 3000

### How to run

Before running the services do remember to run the runDB.sql sql script located in sql folder to load database.

#### Run Passenger service on Port 5000

```bash
cd Passenger
go run .
```

#### Run Driver service on Port 5001

```bash
cd Driver
go run .
```

#### Run Trip service on Port 5002

```bash
cd Trips
go run .
```

#### Run frontend webapp on Port 3000

```bash
cd frontend
npm install
npm run dev
```
