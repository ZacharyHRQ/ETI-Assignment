# ETI-Assignment

## Design consideration

### Separation of concerns for passenger and driver users

### Domain controller

## Architecture Diagram

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
cd frontend/nextjs-blog
npm run dev
```
