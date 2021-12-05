# ETI-Assignment

## Design consideration

### Architecture of microservices

There are 3 microservices (Passenger, Driver and Trips). The Passenenger service will handle the creation and updating of passenger information. The Driver service will handle the creation and updating of driver information. The Trip service will handle the creation of trips and retrieving of trip information. These 3 services will have an MySQL database connected to fetch and store data.

### Separation of concerns for passenger and driver users

It is best practices to allow each portion of the program to handle one portion of application logic. In this case, the choice to separate the passenger and driver was to allow each microservice to handle one set of logic.

### Why use reactjs

Reactjs is a component based web framework which makes User interface resuable which is a good practice of Do not repeat yourself (DRY).

## Architecture Diagram

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
