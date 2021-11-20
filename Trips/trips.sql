Use Ridely;

CREATE TABLE Trips (TripId int PRIMARY KEY auto_increment, PassengerId int, DriverId int, PickUpPostalCode VARCHAR(6), DropOffPostalCode VARCHAR(6), TripStatus TINYINT(1) DEFAULT 1 , DateofTrip TIMESTAMP ,FOREIGN KEY (PassengerId) REFERENCES Passenger(PassengerId), FOREIGN KEY (DriverId) REFERENCES Driver(DriverId) ); 

INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,DateofTrip) VALUES ("1", "1", "999999","000000", NOW());



Select * From Trips;