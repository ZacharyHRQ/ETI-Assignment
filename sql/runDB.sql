CREATE USER 'root'@'localhost' IDENTIFIED VIA mysql_native_password USING PASSWORD('root');
GRANT ALL PRIVILEGES ON *.* TO 'root'@'localhost';

CREATE DATABASE IF NOT EXISTS Ridely;

Use Ridely;

Drop Table if exists Ridely.Trips;
Drop Table if exists Ridely.Passenger;
Drop Table if exists Ridely.Driver;

CREATE TABLE Passenger (PassengerId int  PRIMARY KEY auto_increment, FirstName VARCHAR(60), LastName VARCHAR(60), MoblieNo VARCHAR(8), EmailAddress VARCHAR(120)); 

INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES ("Jake", "Lee", "87394723","jakelee@gmail.com");
INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES ("Frank", "Tan", "87394723","Frank@gmail.com");
INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES ("Patel", "Goh", "87394723","Patel@gmail.com");
INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES ("Mason", "Chong", "87394723","Mason@gmail.com");
INSERT INTO Passenger (FirstName, LastName, MoblieNo, EmailAddress) VALUES ("John", "Lee", "87394723","John@gmail.com");

Select * From Passenger;

CREATE TABLE Driver (DriverId int  PRIMARY KEY auto_increment, FirstName VARCHAR(60), LastName VARCHAR(60), MoblieNo VARCHAR(8), EmailAddress VARCHAR(120), CarLicenseNo VARCHAR(20), IdentificationNumber VARCHAR(20), DriverStatus TINYINT(1) DEFAULT 1 ); 

INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress,CarLicenseNo, IdentificationNumber ,DriverStatus) VALUES ("Dave", "Tan", "93950385","davetan@gmail.com", "SBA1334A", "S0000001I",1);
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress,CarLicenseNo, IdentificationNumber, DriverStatus) VALUES ("Drake", "Cool", "90342784","drakecool@gmail.com", "SBA1434A", "S0000002G", 1);
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress, CarLicenseNo, IdentificationNumber) VALUES ("Jeremy", "Lin", "89237593","linplaysballs@gmail.com","SBA1234A", "S0000003A");
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress, CarLicenseNo, IdentificationNumber, DriverStatus) VALUES ("Smith", "Hills", "90237593","SmithHills@gmail.com","SBA1234A","S0000004C", 0);
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress, CarLicenseNo, IdentificationNumber, DriverStatus) VALUES ("Cindy", "Keller", "90237593","Keller@gmail.com","SBA3234A", "S0000005A", 0);

Select * From Driver;

CREATE TABLE Trips (TripId int PRIMARY KEY auto_increment, PassengerId int, DriverId int, PickUpPostalCode VARCHAR(6), DropOffPostalCode VARCHAR(6), TripStatus TINYINT(1) DEFAULT 0 , DateofTrip TIMESTAMP DEFAULT NOW()); 

INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,TripStatus,DateofTrip) VALUES ("4", "1", "999999","000000", 2, NOW());
INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,TripStatus,DateofTrip) VALUES ("1", "2", "999993","000001", 2, NOW());
INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,TripStatus,DateofTrip) VALUES ("1", "5", "999992","000002", 2, NOW());
INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,TripStatus,DateofTrip) VALUES ("4", "2", "999991","000003", 2,NOW());
INSERT INTO Trips (PassengerId,DriverId,PickUpPostalCode,DropOffPostalCode,TripStatus,DateofTrip) VALUES ("5", "4", "999990","000004", 2,NOW());

Select * From Trips;