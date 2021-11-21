Use Ridely;

CREATE TABLE JWT (Id VARCHAR(38) PRIMARY KEY auto_increment, PassengerId int, DriverId int, PickUpPostalCode VARCHAR(6), DropOffPostalCode VARCHAR(6), TripStatus TINYINT(1) DEFAULT 1 , DateofTrip TIMESTAMP ,FOREIGN KEY (PassengerId) REFERENCES Passenger(PassengerId), FOREIGN KEY (DriverId) REFERENCES Driver(DriverId) ); 

INSERT INTO JWT () VALUES ("1", "1", "999999","000000", NOW());

Select * From JWT;