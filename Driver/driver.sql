Use Ridely;

CREATE TABLE Driver (driverid int  PRIMARY KEY auto_increment, FirstName VARCHAR(60), LastName VARCHAR(60), MoblieNo VARCHAR(8), EmailAddress VARCHAR(120), CarLicenseNo VARCHAR(20), DriverStatus TINYINT(1) DEFAULT 1 ); 

INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress,CarLicenseNo, DriverStatus) VALUES ("Dave", "Tan", "93950385","davetan@gmail.com", "SBA1334A", 1);
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress,CarLicenseNo, DriverStatus) VALUES ("Drake", "Cool", "90342784","drakecool@gmail.com", "SBA1434A", 1);
INSERT INTO Driver (FirstName, LastName, MoblieNo, EmailAddress, CarLicenseNo) VALUES ("Jeremy", "Lin", "89237593","linplaysballs@gmail.com","SBA1234A");

Select * From Driver;