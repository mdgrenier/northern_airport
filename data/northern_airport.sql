--
-- Table structure for table `accountdetails`
--

CREATE TABLE `accountdetails` (
  `AccountDetailID` int(11) NOT NULL,
  `ClientID` int(11) NOT NULL,
  `Username` varchar(25) NOT NULL,
  `Password` varchar(100) NOT NULL,
  `RoleID` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `accountdetails`
--

INSERT INTO `accountdetails` (`AccountDetailID`, `ClientID`, `Username`, `Password`, `RoleID`) VALUES
(8, 3, 'test', '$2a$08$wy3Vf/Bhb27MVa7ZnVXqwujqxKkUtQ/rfHf9aryf84Tzr5zwXNa7q', 1),
(18, 3, 'mgrenier', '$2a$08$OTWXGQu8Cb77EaT9HuL3iOLGXPDJDEft8HAlfjEWYt7BEmAJMwi2q', 1);

-- --------------------------------------------------------

--
-- Table structure for table `airlines`
--

CREATE TABLE `airlines` (
  `AirlineID` int(11) NOT NULL,
  `Name` varchar(50) NOT NULL,
  `Terminal` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `airlines`
--

INSERT INTO `airlines` (`airlineid`, `name`, `terminal`) VALUES
(1, 'Lufthansa', 1),
(2, 'Alitalia', 3),
(3, 'WestJet', 3),
(4, 'Sunwing Airlines', 3),
(5, 'KLM', 3),
(6, 'Air France', 3),
(7, 'Air Transat', 3),
(8, 'Delta', 3),
(9, 'EVA Air', 1),
(10, 'Etihad Airways', 3),
(11, 'Hainan Airlines', 3),
(12, 'CanJet', 3),
(13, 'LOT Polish Airlines', 1),
(14, 'All Nippon Airways', 1),
(15, 'American Airlines', 3),
(16, 'Air New Zealand', 1),
(17, 'Alaska Airlines', 3),
(18, 'Arkefly', 1),
(19, 'Air China', 1),
(20, 'Caribbean Airlines', 3),
(21, 'Copa Airlines', 1),
(22, 'Condor', 3),
(23, 'All Nippon Air', 1),
(24, 'EL AL', 3),
(25, 'Emirates', 1),
(26, 'Finnair', 3),
(27, 'Japan Airlines', 3),
(28, 'Air Canada', 1),
(29, 'Asiana Airlines', 1),
(30, 'Air Wisconsin', 1),
(31, 'Ethiopian Airlines', 1),
(32, 'Aeroflot', 3),
(33, 'Austrian', 1),
(34, 'Avianca', 1),
(35, 'British Airways', 3),
(36, 'Air Jamaica', 1),
(37, 'Cathay Pacific', 3),
(38, 'Cubana', 3),
(39, 'Comar', 3),
(40, 'Fly Jamaica', 3),
(41, 'Ibena', 3),
(42, 'Icelandair', 3),
(43, 'Jet Airways', 1),
(44, 'Air India', 1),
(45, 'Korean Air', 3),
(46, 'LAN Airlines', 3),
(47, 'Aerosvit Ukranian Air', 3),
(48, 'Miami Air International', 3),
(49, 'Pakistan International', 3),
(50, 'Phillippine Airlines', 3),
(51, 'Qantas', 3),
(52, 'SAS Scandinavian', 1),
(53, 'SATA International', 3),
(54, 'SAUDIA', 3),
(55, 'Singapore Airlines', 1),
(56, 'Swiss International Air Lines', 1),
(57, 'TACA', 1),
(58, 'Thai Airways International', 1),
(59, 'Transaero Airlines', 3),
(60, 'Turkish Airlines', 1),
(61, 'Aeromexico', 3),
(62, 'United', 1),
(63, 'US Airways', 3),
(64, 'Commuter', 3),
(65, 'Continental air', 1),
(66, 'Continental Express', 1),
(67, 'Jet Airways', 3),
(68, 'WOW Air', 3),
(69, 'Egypt Air', 1),
(70, 'Aer LIngus', 3),
(71, 'China Eastern', 3),
(72, 'China Eastern Airlines', 3),
(73, 'Brazilian Airlines', 3),
(74, 'American Eagle', 3),
(75, 'Union Pearson Express', 1),
(76, 'Brussels', 1),
(77, 'A Terminal 1  Not flying,', 1),
(78, 'A Terminal 3, Not flying', 3),
(79, 'China South Airlines', 3),
(80, 'TAP Portugal Air', 1),
(81, 'Flair Airline', 3);

--
-- Table structure for table `cities`
--

CREATE TABLE `cities` (
  `CityID` int(11) NOT NULL,
  `Name` varchar(50) NOT NULL
  
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `cities`
--

INSERT INTO `cities` (`cityid`, `name`) VALUES
(1, 'North Bay'),
(2, 'Toronto Pearson and Hotels'),
(3, 'Bracebridge'),
(4, 'Burks Falls'),
(5, 'Callander'),
(7, 'Emsdale'),
(8, 'Gravenhurst'),
(9, 'Huntsville Deerhurst Resort'),
(10, 'Huntsville'),
(11, 'Katrine'),
(14, 'Novar'),
(15, 'Port Sydney'),
(16, 'Powassan'),
(17, 'South River'),
(18, 'Sundridge'),
(20, 'Trout Creek');

-- --------------------------------------------------------

CREATE TABLE `cityoffsets` (
  `CityOffsetID` int(11) NOT NULL,
  `CityID` int(11) NOT NULL,
  `NorthOffset` int(11) NOT NULL,
  `SouthOffset` int(11) NOT NULL
  
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `cities`
--

INSERT INTO `cityoffsets` (`cityoffsetid`, `cityid`, `northoffset`, `southoffset`) VALUES
(1, 1, 0, 240),
(2, 2, 240, 0),
(3, 3, 115, 125),
(4, 4, 55, 185),
(5, 5, 10, 230),
(6, 7, 65, 175),
(7, 8, 135, 105),
(8, 10, 85, 155),
(9, 11, 60, 180),
(10, 14, 75, 165),
(11, 15, 90, 150),
(12, 16, 25, 215),
(13, 17, 40, 200),
(14, 18, 45, 195),
(15, 20, 30, 210);

--
-- Table structure for table `clients`
--

CREATE TABLE `clients` (
  `ClientID` int(11) NOT NULL,
  `FirstName` varchar(50) NOT NULL,
  `LastName` varchar(50) NOT NULL,
  `Phone` bigint(20) NOT NULL,
  `Email` varchar(100) NOT NULL,
  `StreetAddress` varchar(100) NOT NULL,
  `City` varchar(50) NOT NULL,
  `Province` varchar(50) NOT NULL,
  `PostalCode` char(6) NOT NULL,
  `Country` varchar(50) NOT NULL,
  `DefaultDepartureAddress` varchar(100) DEFAULT NULL,
  `DefaultDepartureCityID` int(11) DEFAULT NULL,
  `TravelAgentID` int(11) DEFAULT NULL,
  `Notes` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `clients`
--

INSERT INTO `clients` (`ClientID`, `FirstName`, `LastName`, `Phone`, `Email`, `StreetAddress`, `City`, `Province`, `PostalCode`, `Country`, `DefaultDepartureAddress`, `DefaultDepartureCityID`, `TravelAgentID`, `Notes`) VALUES
(3, 'test', 'test', 1234567890, 'test@test.com', 'test', 'test', 'ontario', 'p1b8p4', 'canada', NULL, 1, NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `customertypes`
--

CREATE TABLE `customertypes` (
  `CustomerTypeID` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `clients`
--

INSERT INTO `customertypes` (`CustomerTypeID`, `Name`) VALUES
(1, 'Adult'),
(2, 'Senior'),
(3, 'Student'),
(4, 'Child');


--
-- Table structure for table `customvenues`
--

CREATE TABLE `customvenues` (
  `CustomVenueID` int(11) NOT NULL,
  `Address` varchar(255) NOT NULL,
  `ClientId` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `departuretimes`
--

CREATE TABLE `departuretimes` (
  `DepartureTimeID` int(11) NOT NULL,
  `CityID` int(11) NOT NULL,
  `DepartureTime` int(11) NOT NULL,
  `Recurring` int(11) NOT NULL,
  `StartDate` date DEFAULT NULL,
  `EndDate` date DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `times`
--

INSERT INTO `departuretimes` (`DepartureTimeID`, `CityID`, `DepartureTime`, `Recurring`, `StartDate`, `EndDate`) VALUES
(1, 1, 0630, 1, null, null),
(2, 1, 1120, 1, null, null),
(3, 1, 1200, 1, null, null),
(4, 1, 2240, 1, null, null),
(6, 2, 1130, 1, null, null),
(5, 2, 1330, 1, null, null),
(7, 2, 1630, 1, null, null),
(8, 2, 1800, 1, null, null);

--
-- Table structure for table `discountcodes`
--

CREATE TABLE `discountcodes` (
  `DiscountCodeID` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL,
  `Percentage` int(11) NOT NULL,
  `Amount` int(11) NOT NULL,
  `StartDate` date NOT NULL,
  `EndDate` date NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `drivers`
--

CREATE TABLE `drivers` (
  `DriverID` int(11) NOT NULL,
  `FirstName` varchar(50) NOT NULL,
  `LastName` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `drivers`
--

INSERT INTO `drivers` (`driverid`, `firstname`, `lastname`) VALUES
(1, 'Andy', 'Duggan'),
(2, 'Joe', 'Palangio'),
(3, 'Mike', 'Lok'),
(4, 'Blake', 'Gennoe'),
(5, 'Rob', 'Farris'),
(6, 'Carole', 'Tran'),
(7, 'Mike', 'Ianiro'),
(8, 'Cosmo', 'Ianiro'),
(9, 'Mike', 'Brideau'),
(10, 'Brad', 'Openshaw'),
(11, 'Jean Paul', 'Chartrand');

--
-- Table structure for table `prices`
--

CREATE TABLE `prices` (
  `PriceID` int(11) NOT NULL,
  `CustomerTypeID` int(11) NOT NULL,
  `Price` float NOT NULL,
  `ReturnPrice` float NOT NULL,
  `DepartureCityID` int(11) NOT NULL,
  `DestinationCityID` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `prices`
--

INSERT INTO `prices` (`priceid`, `departurecityid`, `destinationcityid`, `customertypeid`, `price`, `returnprice`) VALUES
(23, 10, 2, 1, 94.00,0),
(31, 1, 2, 1, 119.00,0),
(32, 1, 2, 2, 111.00,0),
(33, 1, 2, 3, 109.00,0),
(35, 5, 2, 1, 119.00,0),
(37, 5, 2, 3, 109.00,0),
(36, 5, 2, 2, 111.00,0),
(39, 16, 2, 1, 119.00,0),
(40, 16, 2, 4, 59.50,0),
(41, 16, 2, 3, 109.00,0),
(29, 20, 2, 1, 99.00,0),
(30, 20, 2, 4, 47.50,0),
(45, 20, 2, 3, 95.00,169.00),
(46, 20, 2, 2, 95.00,169.00),
(47, 17, 2, 1, 99.00,0),
(3, 17, 2, 4, 49.50,0),
(49, 17, 2, 3, 95.00,169.00),
(50, 17, 2, 2, 95.00,169.00),
(51, 18, 2, 1, 99.00,0),
(5, 18, 2, 4, 49.50,0),
(53, 18, 2, 3, 95.00,169.00),
(54, 18, 2, 2, 95.00,169.00),
(55, 4, 2, 1, 99.00,0),
(21, 4, 2, 2, 95.00,0),
(57, 4, 2, 3, 95.00,0),
(59, 11, 2, 1, 99.00,0),
(17, 11, 2, 4, 49.50,0),
(61, 11, 2, 3, 95.00,0),
(62, 11, 2, 2, 95.00,0),
(63, 7, 2, 1, 99.00,0),
(65, 7, 2, 2, 95.00,0),
(24, 7, 2, 3, 95.00,0),
(67, 14, 2, 1, 99.00,0),
(18, 14, 2, 4, 49.50,0),
(69, 14, 2, 2, 95.00,0),
(48, 10, 2, 2, 89.00,0),
(22, 10, 2, 4, 47.00,0),
(75, 3, 2, 1, 94.00,0),
(6, 3, 2, 4, 47.00,0),
(77, 3, 2, 2, 89.00,0),
(79, 15, 2, 1, 94.00,0),
(80, 15, 2, 4, 47.00,0),
(81, 15, 2, 3, 89.00,159.00),
(82, 15, 2, 2, 89.00,159.00),
(85, 8, 2, 4, 47.00,0),
(83, 8, 2, 1, 94.00,0),
(84, 8, 2, 2, 89.00,0),
(89, 2, 11, 4, 47.50,0),
(88, 2, 11, 3, 95.00,0),
(73, 2, 11, 2, 95.00,0),
(92, 2, 1, 3, 109.00,0),
(91, 2, 1, 2, 111.00,0),
(90, 2, 1, 1, 119.00,0),
(9, 2, 4, 1, 99.00,0),
(8, 2, 4, 4, 49.50,0),
(2, 2, 5, 2, 111.00,0),
(1, 2, 5, 1, 119.00,0),
(95, 2, 10, 4, 47.00,0),
(56, 2, 11, 1, 99.00,0),
(102, 2, 1, 4, 59.50,0),
(10, 2, 4, 2, 95.00,0),
(4, 2, 5, 4, 59.50,0),
(71, 2, 7, 4, 49.50,0),
(68, 2, 7, 3, 95.00,0),
(64, 2, 7, 2, 95.00,0),
(60, 2, 7, 1, 99.00,0),
(103, 2, 8, 4, 47.00,0),
(101, 2, 8, 3, 89.00,0),
(100, 2, 8, 2, 89.00,0),
(99, 2, 8, 1, 94.00,0),
(107, 9, 2, 4, 47.00,0),
(94, 2, 10, 3, 89.00,0),
(93, 2, 10, 2, 89.00,0),
(43, 2, 10, 1, 94.00,0),
(98, 2, 3, 4, 47.00,0),
(119, 2, 3, 3, 89.00,0),
(44, 2, 3, 1, 94.00,0),
(106, 9, 2, 3, 119.00,0),
(105, 9, 2, 2, 119.00,0),
(72, 2, 14, 1, 99.00,0),
(74, 2, 14, 2, 95.00,0),
(97, 2, 15, 4, 47.00,0),
(132, 2, 15, 3, 89.00,159.00),
(133, 2, 15, 2, 89.00,159.00),
(96, 2, 15, 1, 94.00,0),
(15, 2, 16, 4, 59.50,0),
(13, 2, 16, 2, 111.00,0),
(14, 2, 16, 3, 109.00,0),
(12, 2, 16, 1, 119.00,0),
(26, 2, 17, 3, 95.00,0),
(25, 2, 17, 2, 95.00,0),
(20, 2, 17, 1, 99.00,0),
(121, 2, 18, 1, 99.00,0),
(144, 2, 18, 2, 95.00,169.00),
(145, 2, 18, 3, 95.00,169.00),
(28, 2, 18, 4, 49.50,0),
(104, 9, 2, 1, 124.00,0),
(19, 2, 20, 4, 49.50,0),
(150, 2, 20, 2, 95.00,169.00),
(151, 2, 20, 3, 95.00,169.00),
(16, 2, 20, 1, 99.00,0),
(115, 2, 3, 2, 89.00,0),
(87, 2, 14, 4, 49.50,0),
(76, 2, 14, 3, 95.00,0),
(27, 2, 17, 4, 49.50,0),
(11, 2, 4, 3, 95.00,0),
(7, 2, 5, 3, 109.00,0),
(78, 3, 2, 3, 89.00,0),
(58, 4, 2, 4, 49.50,0),
(38, 5, 2, 4, 59.50,0),
(66, 7, 2, 4, 49.50,0),
(86, 8, 2, 3, 89.00,0),
(52, 10, 2, 3, 89.00,0),
(34, 1, 2, 4, 59.50,0),
(70, 14, 2, 3, 95.00,0),
(42, 16, 2, 2, 111.00,0);


--
-- Table structure for table `reservations`
--

CREATE TABLE `reservations` (
  `ReservationID` int(11) NOT NULL,
  `ClientID` int(11) NOT NULL,
  `ReservationTypeID` int(11) NOT NULL,
  `DepartureCityID` int(11) NOT NULL,
  `DepartureVenueID` int(11) NOT NULL,
  `DepartureTimeID` int(11) NOT NULL,
  `DestinationCityID` int(11) NOT NULL,
  `DestinationVenueID` int(11) NOT NULL,
  `DiscountCodeID` int(11) NOT NULL,
  `DepartureAirlineID` int(11) NOT NULL,
  `ReturnAirlineID` int(11) NOT NULL,
  `DriverNotes` text,
  `InternalNotes` text,
  `NumAdults` int(11) NOT NULL,
  `NumStudents` int(11) NOT NULL,
  `NumSeniors` int(11) NOT NULL,
  `Price` float NOT NULL,
  `Status` varchar(25) NOT NULL,
  `Hash` varchar(255) NOT NULL,
  `CustomDepartureID` int(11) DEFAULT NULL,
  `CustomDestinationID` int(11) DEFAULT NULL,
  `ReturnDate` date DEFAULT NULL,
  `TripTypeID` int(11) NOT NULL,
  `BalanceOwing` float NOT NULL,
  `ElvaonTransactionID` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `reservationtypes`
--

CREATE TABLE `reservationtypes` (
  `ReservationTypeID` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `RoleID` int(11) NOT NULL,
  `RoleName` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`roleid`, `RoleName`) VALUES
(1, 'test'),
(2, 'client'),
(3, 'staff'),
(4, 'admin');

-- --------------------------------------------------------

--
-- Table structure for table `taxes`
--

CREATE TABLE `taxes` (
  `TaxID` int(11) NOT NULL,
  `Percentage` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL,
  `Active` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `roles`
--

INSERT INTO `taxes` (`TaxID`, `Percentage`, `Name`, `Active`) VALUES
(1, '13', 'HST', 1);


--
-- Table structure for table `transactions`
--

CREATE TABLE `transactions` (
  `TransactionID` int(11) NOT NULL,
  `Type` varchar(15) NOT NULL,
  `Created` datetime NOT NULL,
  `ReservationID` int(11) NOT NULL,
  `Response` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `travelagencies`
--

CREATE TABLE `travelagencies` (
  `TravelAgencyID` int(11) NOT NULL,
  `TravelAgencyName` varchar(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `travelagents`
--

CREATE TABLE `travelagents` (
  `TravelAgentID` int(11) NOT NULL,
  `TravelAgentName` varchar(100) NOT NULL,
  `IATANumber` varchar(100) NOT NULL,
  `TravelAgencyID` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `trips`
--

CREATE TABLE `trips` (
  `TripID` int(11) NOT NULL,
  `DepartureDate` date NOT NULL,
  `DepartureTimeID` int(11) NOT NULL,
  `ReservationID` int(11) NOT NULL,
  `NumPassengers` int(11) NOT NULL,
  `DriverID` int(11) DEFAULT NULL,
  `VehicleID` int(11) DEFAULT NULL,
  `OmitTrip` int(11) DEFAULT NULL,
  `Postpone` int(11) DEFAULT NULL,
  `RescheduleDate` date DEFAULT NULL,
  `RescheduleTime` int(11) DEFAULT NULL,
  `Cancelled` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `triptypes`
--

CREATE TABLE `triptypes` (
  `TripTypeID` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- --------------------------------------------------------

--
-- Table structure for table `vehicles`
--

CREATE TABLE `vehicles` (
  `VehicleID` int(11) NOT NULL,
  `LicensePlate` varchar(25) NOT NULL,
  `NumSeats` int(11) NOT NULL,
  `Make` varchar(25) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `vehicles`
--

INSERT INTO `vehicles` (`vehicleid`, `licenseplate`, `make`, `numseats`) VALUES
(2, '1109 lic. 8362BF', 'Mercedes', 10),
(3, '1512 lic. 8697BH', 'Mercedes', 14),
(1, '1110 lic. 7405BF', 'Mercedes', 10),
(4, '1411 lic. 8691BH', 'Mercedes', 11);

--
-- Table structure for table `venues`
--

CREATE TABLE `venues` (
  `VenueID` int(11) NOT NULL,
  `CityID` int(11) NOT NULL,
  `Name` varchar(100) NOT NULL,
  `ExtraCost` float DEFAULT NULL,
  `Active` int(11) NOT NULL,
  `ExtraTime` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- Dumping data for table `venues`
--

INSERT INTO `venues` (`venueid`, `cityid`, `name`, `extracost`, `active`, `extratime`) VALUES
(1, 8, 'McDonald''s Restaurant located at Bethune Dr. and Muskoka Rd', 0.00, 1, 0),
(2, 3, 'Tim Hortons coffee shop corner of Depot Rd. and Taylor Rd', 0.00, 1, 0),
(3, 4, 'Tim Horton''s located at 27 Commercial Dr', 0.00, 1, 0),
(4, 5, 'Wasi Petro Canada located at the corner of Wasi Rd and Callander Bay Dr', 0.00, 1, 0),
(5, 7, 'Tourist information centre located on Emsdale Rd. Exit #244', 0.00, 1, 0),
(6, 10, 'Holiday Inn Express located at Howland Dr', 0.00, 1, 0),
(7, 11, 'The Lucky Dollar on old Hwy 11', 0.00, 1, 0),
(8, 1, 'Country Style coffee shop located at 1401 Seymour St. and Hwy 11', 0.00, 1, 0),
(9, 14, 'Foodland', 0.00, 1, 0),
(10, 15, 'Smiths Resaurant located on highway 11', 0.00, 1, 0),
(11, 16, 'East Himsworth Cafe, previous Ethels', 0.00, 1, 0),
(12, 17, 'Shell gas station highway 124', 0.00, 1, 0),
(13, 18, 'Blue Roof restaurant', 0.00, 1, 0),
(14, 20, 'TJ''s Variety on highway 522', 0.00, 1, 0),
(99, 2, 'AIRPORT AIRLINES: CLICK HERE FOR AIRLINE LIST', 0.00, 1, 0),
(100, 1, 'Home Pickup or Dropoff', 15.00, 1, 30),
(15, 2, 'Nu Hotel-6465 Airport Rd.', 0.00, 1, 0),
(19, 2, 'UP Express train Pearson Terminal 1', 0.00, 1, 0),
(17, 2, 'YYZ No Flight TERMINAL 1,', 0.00, 0, 0),
(49, 10, 'Deerhurst Resort', 30.00, 1, 0),
(18, 2, 'YYZ No Flight  Terminal 3,', 0.00, 0, 0),
(48, 1, 'North Bay Office 191 Booth Rd #7', 0.00, 1, 0),
(16, 2, 'Quality Inn 2180 Islington', 0.00, 1, 0),
(25, 2, 'Sheraton Airport, 801 Dixon Rd', 0.00, 1, 0),
(26, 2, 'Comfort Inn, 6355 Airport Rd.', 0.00, 1, 0),
(27, 2, 'Carlingview, 221 Carlingview Dr.', 0.00, 1, 0),
(28, 2, 'Crown Plaza, 33 Carlson Crt.', 0.00, 1, 0),
(29, 2, 'Delta Hotels 655 Dixon Rd', 0.00, 1, 0),
(30, 2, 'Fairfield Inn, 3299 Caroga Dr.', 0.00, 1, 0),
(31, 2, 'Four Points Sheraton, 6257 Airport Rd.', 0.00, 1, 0),
(32, 2, 'Hampton Inn, 3279 Caroga Dr.', 0.00, 1, 0),
(33, 2, 'Hilton, 5875 Airport Rd', 0.00, 1, 0),
(34, 2, 'Holiday Inn , 970 Dixon Rd', 0.00, 1, 0),
(35, 2, 'Holiday Inn, 600 Dixon Rd.', 0.00, 1, 0),
(36, 2, 'Embassy Suites, 262 Carlingview Dr.', 0.00, 1, 0),
(37, 2, 'Hilton Garden Inn,  3311 Caroga Dr.', 0.00, 1, 0),
(38, 2, 'Radisson Suite, 640 Dixon Rd', 0.00, 1, 0),
(39, 2, 'Residence Inn Marriott, 17 Reading Crt.', 0.00, 1, 0),
(40, 2, 'Sheraton Gateway, PIA terminal 3', 0.00, 1, 0),
(41, 2, 'Marriott Toronto Airport, 901 Dixon Rd.', 0.00, 1, 0),
(42, 2, 'Double Tree,  925 Dixon Rd.', 0.00, 1, 0),
(43, 2, 'The Westin Toronto Airport, 950 Dixon Rd.', 0.00, 1, 0),
(44, 2, 'Best Western Premier, 135 Carlingview Dr.', 0.00, 1, 0),
(45, 2, 'Sandman Signature, 55 Reading Crt.', 0.00, 1, 0),
(46, 2, 'Alt Hotel,  6080 Viscount Rd.', 0.00, 1, 0),
(47, 2, 'Courtyard by Marriot-231 Carlingview', 0.00, 1, 0);

--
-- Indexes for table `accountdetails`
--
ALTER TABLE `accountdetails`
  ADD PRIMARY KEY (`AccountDetailID`),
  ADD KEY `FK_148` (`ClientID`),
  ADD KEY `FK_266` (`RoleID`);

--
-- Indexes for table `airlines`
--
ALTER TABLE `airlines`
  ADD PRIMARY KEY (`AirlineID`);

--
-- Indexes for table `cities`
--
ALTER TABLE `cities`
  ADD PRIMARY KEY (`CityID`);

--
-- Indexes for table `clients`
--
ALTER TABLE `clients`
  ADD PRIMARY KEY (`ClientID`),
  ADD UNIQUE KEY `NONCLUSTERED` (`FirstName`),
  ADD KEY `FK_171` (`TravelAgentID`),
  ADD KEY `FK_217` (`DefaultDepartureCityID`);

--
-- Indexes for table `customertypes`
--
ALTER TABLE `customertypes`
  ADD PRIMARY KEY (`CustomerTypeID`);

--
-- Indexes for table `customvenues`
--
ALTER TABLE `customvenues`
  ADD PRIMARY KEY (`CustomVenueID`),
  ADD KEY `FK_327` (`ClientId`);

--
-- Indexes for table `departuretimes`
--
ALTER TABLE `departuretimes`
  ADD PRIMARY KEY (`DepartureTimeID`);

--
-- Indexes for table `discountcodes`
--
ALTER TABLE `discountcodes`
  ADD PRIMARY KEY (`DiscountCodeID`);

--
-- Indexes for table `drivers`
--
ALTER TABLE `drivers`
  ADD PRIMARY KEY (`DriverID`);

--
-- Indexes for table `prices`
--
ALTER TABLE `prices`
  ADD PRIMARY KEY (`PriceID`),
  ADD KEY `FK_232` (`CustomerTypeID`),
  ADD KEY `FK_236` (`DepartureCityID`),
  ADD KEY `FK_239` (`DestinationCityID`);

--
-- Indexes for table `reservations`
--
ALTER TABLE `reservations`
  ADD PRIMARY KEY (`ReservationID`),
  ADD KEY `FK_155` (`ReservationTypeID`),
  ADD KEY `FK_163` (`DepartureCityID`),
  ADD KEY `FK_181` (`DiscountCodeID`),
  ADD KEY `FK_195` (`DepartureTimeID`),
  ADD KEY `FK_203` (`DestinationCityID`),
  ADD KEY `FK_211` (`DepartureVenueID`),
  ADD KEY `FK_214` (`DestinationVenueID`),
  ADD KEY `FK_273` (`DepartureAirlineID`),
  ADD KEY `FK_311` (`ReturnAirlineID`),
  ADD KEY `FK_330` (`CustomDepartureID`),
  ADD KEY `FK_333` (`CustomDestinationID`),
  ADD KEY `FK_341` (`TripTypeID`);

--
-- Indexes for table `reservationtypes`
--
ALTER TABLE `reservationtypes`
  ADD PRIMARY KEY (`ReservationTypeID`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`RoleID`);

--
-- Indexes for table `taxes`
--
ALTER TABLE `taxes`
  ADD PRIMARY KEY (`TaxID`,`Percentage`);

--
-- Indexes for table `transactions`
--
ALTER TABLE `transactions`
  ADD PRIMARY KEY (`TransactionID`),
  ADD KEY `FK_301` (`ReservationID`);

--
-- Indexes for table `travelagencies`
--
ALTER TABLE `travelagencies`
  ADD PRIMARY KEY (`TravelAgencyID`);

--
-- Indexes for table `travelagents`
--
ALTER TABLE `travelagents`
  ADD PRIMARY KEY (`TravelAgentID`),
  ADD KEY `FK_354` (`TravelAgencyID`);

--
-- Indexes for table `trips`
--
ALTER TABLE `trips`
  ADD PRIMARY KEY (`TripID`),
  ADD KEY `FK_254` (`DepartureTimeID`),
  ADD KEY `FK_259` (`ReservationID`),
  ADD KEY `FK_288` (`DriverID`),
  ADD KEY `FK_291` (`VehicleID`);

--
-- Indexes for table `triptypes`
--
ALTER TABLE `triptypes`
  ADD PRIMARY KEY (`TripTypeID`);

--
-- Indexes for table `vehicles`
--
ALTER TABLE `vehicles`
  ADD PRIMARY KEY (`VehicleID`);

--
-- Indexes for table `venues`
--
ALTER TABLE `venues`
  ADD PRIMARY KEY (`VenueID`),
  ADD KEY `FK_206` (`CityID`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `accountdetails`
--
ALTER TABLE `accountdetails`
  MODIFY `AccountDetailID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=19;

--
-- AUTO_INCREMENT for table `airlines`
--
ALTER TABLE `airlines`
  MODIFY `AirlineID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `cities`
--
ALTER TABLE `cities`
  MODIFY `CityID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `clients`
--
ALTER TABLE `clients`
  MODIFY `ClientID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `customertypes`
--
ALTER TABLE `customertypes`
  MODIFY `CustomerTypeID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `customvenues`
--
ALTER TABLE `customvenues`
  MODIFY `CustomVenueID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `departuretimes`
--
ALTER TABLE `departuretimes`
  MODIFY `DepartureTimeID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `discountcodes`
--
ALTER TABLE `discountcodes`
  MODIFY `DiscountCodeID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `drivers`
--
ALTER TABLE `drivers`
  MODIFY `DriverID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `prices`
--
ALTER TABLE `prices`
  MODIFY `PriceID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `reservations`
--
ALTER TABLE `reservations`
  MODIFY `ReservationID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `reservationtypes`
--
ALTER TABLE `reservationtypes`
  MODIFY `ReservationTypeID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `RoleID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `taxes`
--
ALTER TABLE `taxes`
  MODIFY `TaxID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `transactions`
--
ALTER TABLE `transactions`
  MODIFY `TransactionID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `travelagencies`
--
ALTER TABLE `travelagencies`
  MODIFY `TravelAgencyID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `travelagents`
--
ALTER TABLE `travelagents`
  MODIFY `TravelAgentID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `trips`
--
ALTER TABLE `trips`
  MODIFY `TripID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `triptypes`
--
ALTER TABLE `triptypes`
  MODIFY `TripTypeID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `vehicles`
--
ALTER TABLE `vehicles`
  MODIFY `VehicleID` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `venues`
--
ALTER TABLE `venues`
  MODIFY `VenueID` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `accountdetails`
--
ALTER TABLE `accountdetails`
  ADD CONSTRAINT `FK_148` FOREIGN KEY (`ClientID`) REFERENCES `clients` (`ClientID`),
  ADD CONSTRAINT `FK_266` FOREIGN KEY (`RoleID`) REFERENCES `roles` (`RoleID`);

--
-- Constraints for table `clients`
--
ALTER TABLE `clients`
  ADD CONSTRAINT `FK_171` FOREIGN KEY (`TravelAgentID`) REFERENCES `travelagents` (`TravelAgentID`),
  ADD CONSTRAINT `FK_217` FOREIGN KEY (`DefaultDepartureCityID`) REFERENCES `cities` (`CityID`);

--
-- Constraints for table `customvenues`
--
ALTER TABLE `customvenues`
  ADD CONSTRAINT `FK_327` FOREIGN KEY (`ClientId`) REFERENCES `clients` (`ClientID`);

--
-- Constraints for table `prices`
--
ALTER TABLE `prices`
  ADD CONSTRAINT `FK_232` FOREIGN KEY (`CustomerTypeID`) REFERENCES `customertypes` (`CustomerTypeID`),
  ADD CONSTRAINT `FK_236` FOREIGN KEY (`DepartureCityID`) REFERENCES `cities` (`CityID`),
  ADD CONSTRAINT `FK_239` FOREIGN KEY (`DestinationCityID`) REFERENCES `cities` (`CityID`);

--
-- Constraints for table `reservations`
--
ALTER TABLE `reservations`
  ADD CONSTRAINT `FK_155` FOREIGN KEY (`ReservationTypeID`) REFERENCES `reservationtypes` (`ReservationTypeID`),
  ADD CONSTRAINT `FK_163` FOREIGN KEY (`DepartureCityID`) REFERENCES `cities` (`CityID`),
  ADD CONSTRAINT `FK_181` FOREIGN KEY (`DiscountCodeID`) REFERENCES `discountcodes` (`DiscountCodeID`),
  ADD CONSTRAINT `FK_195` FOREIGN KEY (`DepartureTimeID`) REFERENCES `departuretimes` (`DepartureTimeID`),
  ADD CONSTRAINT `FK_203` FOREIGN KEY (`DestinationCityID`) REFERENCES `cities` (`CityID`),
  ADD CONSTRAINT `FK_211` FOREIGN KEY (`DepartureVenueID`) REFERENCES `venues` (`VenueID`),
  ADD CONSTRAINT `FK_214` FOREIGN KEY (`DestinationVenueID`) REFERENCES `venues` (`VenueID`),
  ADD CONSTRAINT `FK_273` FOREIGN KEY (`DepartureAirlineID`) REFERENCES `airlines` (`AirlineID`),
  ADD CONSTRAINT `FK_311` FOREIGN KEY (`ReturnAirlineID`) REFERENCES `trips` (`TripID`),
  ADD CONSTRAINT `FK_330` FOREIGN KEY (`CustomDepartureID`) REFERENCES `customvenues` (`CustomVenueID`),
  ADD CONSTRAINT `FK_333` FOREIGN KEY (`CustomDestinationID`) REFERENCES `customvenues` (`CustomVenueID`),
  ADD CONSTRAINT `FK_341` FOREIGN KEY (`TripTypeID`) REFERENCES `triptypes` (`TripTypeID`);

--
-- Constraints for table `transactions`
--
ALTER TABLE `transactions`
  ADD CONSTRAINT `FK_301` FOREIGN KEY (`ReservationID`) REFERENCES `reservations` (`ReservationID`);

--
-- Constraints for table `travelagents`
--
ALTER TABLE `travelagents`
  ADD CONSTRAINT `FK_354` FOREIGN KEY (`TravelAgencyID`) REFERENCES `travelagencies` (`TravelAgencyID`);

--
-- Constraints for table `trips`
--
ALTER TABLE `trips`
  ADD CONSTRAINT `FK_254` FOREIGN KEY (`DepartureTimeID`) REFERENCES `departuretimes` (`DepartureTimeID`),
  ADD CONSTRAINT `FK_259` FOREIGN KEY (`ReservationID`) REFERENCES `reservations` (`ReservationID`),
  ADD CONSTRAINT `FK_288` FOREIGN KEY (`DriverID`) REFERENCES `drivers` (`DriverID`),
  ADD CONSTRAINT `FK_291` FOREIGN KEY (`VehicleID`) REFERENCES `vehicles` (`VehicleID`);

--
-- Constraints for table `venues`
--
ALTER TABLE `venues`
  ADD CONSTRAINT `FK_206` FOREIGN KEY (`CityID`) REFERENCES `cities` (`CityID`);
--
-- Database: `phpmyadmin`
--
CREATE DATABASE IF NOT EXISTS `phpmyadmin` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
USE `phpmyadmin`;

-- --------------------------------------------------------

--
-- Table structure for table `pma__bookmark`
--

CREATE TABLE `pma__bookmark` (
  `id` int(11) NOT NULL,
  `dbase` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `user` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `label` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `query` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Bookmarks';

-- --------------------------------------------------------

--
-- Table structure for table `pma__central_columns`
--

CREATE TABLE `pma__central_columns` (
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `col_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `col_type` varchar(64) COLLATE utf8_bin NOT NULL,
  `col_length` text COLLATE utf8_bin,
  `col_collation` varchar(64) COLLATE utf8_bin NOT NULL,
  `col_isNull` tinyint(1) NOT NULL,
  `col_extra` varchar(255) COLLATE utf8_bin DEFAULT '',
  `col_default` text COLLATE utf8_bin
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Central list of columns';

-- --------------------------------------------------------

--
-- Table structure for table `pma__column_info`
--

CREATE TABLE `pma__column_info` (
  `id` int(5) UNSIGNED NOT NULL,
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `column_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `comment` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `mimetype` varchar(255) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `transformation` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `transformation_options` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `input_transformation` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT '',
  `input_transformation_options` varchar(255) COLLATE utf8_bin NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Column information for phpMyAdmin';

-- --------------------------------------------------------

--
-- Table structure for table `pma__designer_settings`
--

CREATE TABLE `pma__designer_settings` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `settings_data` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Settings related to Designer';

-- --------------------------------------------------------

--
-- Table structure for table `pma__export_templates`
--

CREATE TABLE `pma__export_templates` (
  `id` int(5) UNSIGNED NOT NULL,
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `export_type` varchar(10) COLLATE utf8_bin NOT NULL,
  `template_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `template_data` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Saved export templates';

--
-- Dumping data for table `pma__export_templates`
--

INSERT INTO `pma__export_templates` (`id`, `username`, `export_type`, `template_name`, `template_data`) VALUES
(1, 'root', 'database', 'northern_airport_empty', '{\"quick_or_custom\":\"custom\",\"what\":\"sql\",\"structure_or_data_forced\":\"0\",\"table_select[]\":[\"airlines\",\"cities\",\"clients\",\"customertypes\",\"customvenues\",\"departuretimes\",\"discountcodes\",\"drivers\",\"reservationtypes\",\"roles\",\"taxes\",\"travelagencies\",\"travelagents\",\"triptypes\",\"vehicles\",\"venues\"],\"table_structure[]\":[\"airlines\",\"cities\",\"clients\",\"customertypes\",\"customvenues\",\"departuretimes\",\"discountcodes\",\"drivers\",\"reservationtypes\",\"roles\",\"taxes\",\"travelagencies\",\"travelagents\",\"triptypes\",\"vehicles\",\"venues\"],\"table_data[]\":[\"airlines\",\"cities\",\"clients\",\"customertypes\",\"customvenues\",\"departuretimes\",\"discountcodes\",\"drivers\",\"reservationtypes\",\"roles\",\"taxes\",\"travelagencies\",\"travelagents\",\"triptypes\",\"vehicles\",\"venues\"],\"aliases_new\":\"\",\"output_format\":\"sendit\",\"filename_template\":\"@DATABASE@\",\"remember_template\":\"on\",\"charset\":\"utf-8\",\"compression\":\"none\",\"maxsize\":\"\",\"codegen_structure_or_data\":\"data\",\"codegen_format\":\"0\",\"csv_separator\":\",\",\"csv_enclosed\":\"\\\"\",\"csv_escaped\":\"\\\"\",\"csv_terminated\":\"AUTO\",\"csv_null\":\"NULL\",\"csv_structure_or_data\":\"data\",\"excel_null\":\"NULL\",\"excel_columns\":\"something\",\"excel_edition\":\"win\",\"excel_structure_or_data\":\"data\",\"json_structure_or_data\":\"data\",\"json_unicode\":\"something\",\"latex_caption\":\"something\",\"latex_structure_or_data\":\"structure_and_data\",\"latex_structure_caption\":\"Structure of table @TABLE@\",\"latex_structure_continued_caption\":\"Structure of table @TABLE@ (continued)\",\"latex_structure_label\":\"tab:@TABLE@-structure\",\"latex_relation\":\"something\",\"latex_comments\":\"something\",\"latex_mime\":\"something\",\"latex_columns\":\"something\",\"latex_data_caption\":\"Content of table @TABLE@\",\"latex_data_continued_caption\":\"Content of table @TABLE@ (continued)\",\"latex_data_label\":\"tab:@TABLE@-data\",\"latex_null\":\"\\\\textit{NULL}\",\"mediawiki_structure_or_data\":\"structure_and_data\",\"mediawiki_caption\":\"something\",\"mediawiki_headers\":\"something\",\"htmlword_structure_or_data\":\"structure_and_data\",\"htmlword_null\":\"NULL\",\"ods_null\":\"NULL\",\"ods_structure_or_data\":\"data\",\"odt_structure_or_data\":\"structure_and_data\",\"odt_relation\":\"something\",\"odt_comments\":\"something\",\"odt_mime\":\"something\",\"odt_columns\":\"something\",\"odt_null\":\"NULL\",\"pdf_report_title\":\"\",\"pdf_structure_or_data\":\"structure_and_data\",\"phparray_structure_or_data\":\"data\",\"sql_include_comments\":\"something\",\"sql_header_comment\":\"\",\"sql_use_transaction\":\"something\",\"sql_compatibility\":\"NONE\",\"sql_structure_or_data\":\"structure_and_data\",\"sql_create_table\":\"something\",\"sql_auto_increment\":\"something\",\"sql_create_view\":\"something\",\"sql_procedure_function\":\"something\",\"sql_create_trigger\":\"something\",\"sql_backquotes\":\"something\",\"sql_type\":\"INSERT\",\"sql_insert_syntax\":\"both\",\"sql_max_query_size\":\"50000\",\"sql_hex_for_binary\":\"something\",\"sql_utc_time\":\"something\",\"texytext_structure_or_data\":\"structure_and_data\",\"texytext_null\":\"NULL\",\"xml_structure_or_data\":\"data\",\"xml_export_events\":\"something\",\"xml_export_functions\":\"something\",\"xml_export_procedures\":\"something\",\"xml_export_tables\":\"something\",\"xml_export_triggers\":\"something\",\"xml_export_views\":\"something\",\"xml_export_contents\":\"something\",\"yaml_structure_or_data\":\"data\",\"\":null,\"lock_tables\":null,\"as_separate_files\":null,\"csv_removeCRLF\":null,\"csv_columns\":null,\"excel_removeCRLF\":null,\"json_pretty_print\":null,\"htmlword_columns\":null,\"ods_columns\":null,\"sql_dates\":null,\"sql_relation\":null,\"sql_mime\":null,\"sql_disable_fk\":null,\"sql_views_as_tables\":null,\"sql_metadata\":null,\"sql_create_database\":null,\"sql_drop_table\":null,\"sql_if_not_exists\":null,\"sql_truncate\":null,\"sql_delayed\":null,\"sql_ignore\":null,\"texytext_columns\":null}');

-- --------------------------------------------------------

--
-- Table structure for table `pma__favorite`
--

CREATE TABLE `pma__favorite` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `tables` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Favorite tables';

-- --------------------------------------------------------

--
-- Table structure for table `pma__history`
--

CREATE TABLE `pma__history` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `username` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `db` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `table` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `timevalue` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sqlquery` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='SQL history for phpMyAdmin';

-- --------------------------------------------------------

--
-- Table structure for table `pma__navigationhiding`
--

CREATE TABLE `pma__navigationhiding` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `item_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `item_type` varchar(64) COLLATE utf8_bin NOT NULL,
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Hidden items of navigation tree';

-- --------------------------------------------------------

--
-- Table structure for table `pma__pdf_pages`
--

CREATE TABLE `pma__pdf_pages` (
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `page_nr` int(10) UNSIGNED NOT NULL,
  `page_descr` varchar(50) CHARACTER SET utf8 NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='PDF relation pages for phpMyAdmin';

-- --------------------------------------------------------

--
-- Table structure for table `pma__recent`
--

CREATE TABLE `pma__recent` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `tables` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Recently accessed tables';

--
-- Dumping data for table `pma__recent`
--

INSERT INTO `pma__recent` (`username`, `tables`) VALUES
('root', '[{\"db\":\"northernairport\",\"table\":\"accountdetails\"},{\"db\":\"northernairport\",\"table\":\"roles\"},{\"db\":\"northernairport\",\"table\":\"taxes\"},{\"db\":\"northernairport\",\"table\":\"clients\"},{\"db\":\"northernairport\",\"table\":\"cities\"}]');

-- --------------------------------------------------------

--
-- Table structure for table `pma__relation`
--

CREATE TABLE `pma__relation` (
  `master_db` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `master_table` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `master_field` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `foreign_db` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `foreign_table` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `foreign_field` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Relation table';

-- --------------------------------------------------------

--
-- Table structure for table `pma__savedsearches`
--

CREATE TABLE `pma__savedsearches` (
  `id` int(5) UNSIGNED NOT NULL,
  `username` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `search_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `search_data` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Saved searches';

-- --------------------------------------------------------

--
-- Table structure for table `pma__table_coords`
--

CREATE TABLE `pma__table_coords` (
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `pdf_page_number` int(11) NOT NULL DEFAULT '0',
  `x` float UNSIGNED NOT NULL DEFAULT '0',
  `y` float UNSIGNED NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Table coordinates for phpMyAdmin PDF output';

-- --------------------------------------------------------

--
-- Table structure for table `pma__table_info`
--

CREATE TABLE `pma__table_info` (
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT '',
  `display_field` varchar(64) COLLATE utf8_bin NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Table information for phpMyAdmin';

-- --------------------------------------------------------

--
-- Table structure for table `pma__table_uiprefs`
--

CREATE TABLE `pma__table_uiprefs` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `prefs` text COLLATE utf8_bin NOT NULL,
  `last_update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Tables'' UI preferences';

-- --------------------------------------------------------

--
-- Table structure for table `pma__tracking`
--

CREATE TABLE `pma__tracking` (
  `db_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `table_name` varchar(64) COLLATE utf8_bin NOT NULL,
  `version` int(10) UNSIGNED NOT NULL,
  `date_created` datetime NOT NULL,
  `date_updated` datetime NOT NULL,
  `schema_snapshot` text COLLATE utf8_bin NOT NULL,
  `schema_sql` text COLLATE utf8_bin,
  `data_sql` longtext COLLATE utf8_bin,
  `tracking` set('UPDATE','REPLACE','INSERT','DELETE','TRUNCATE','CREATE DATABASE','ALTER DATABASE','DROP DATABASE','CREATE TABLE','ALTER TABLE','RENAME TABLE','DROP TABLE','CREATE INDEX','DROP INDEX','CREATE VIEW','ALTER VIEW','DROP VIEW') COLLATE utf8_bin DEFAULT NULL,
  `tracking_active` int(1) UNSIGNED NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Database changes tracking for phpMyAdmin';

-- --------------------------------------------------------

--
-- Table structure for table `pma__userconfig`
--

CREATE TABLE `pma__userconfig` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `timevalue` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `config_data` text COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='User preferences storage for phpMyAdmin';

--
-- Dumping data for table `pma__userconfig`
--

INSERT INTO `pma__userconfig` (`username`, `timevalue`, `config_data`) VALUES
('root', '2019-03-01 03:01:37', '{\"Console\\/Mode\":\"collapse\"}');

-- --------------------------------------------------------

--
-- Table structure for table `pma__usergroups`
--

CREATE TABLE `pma__usergroups` (
  `usergroup` varchar(64) COLLATE utf8_bin NOT NULL,
  `tab` varchar(64) COLLATE utf8_bin NOT NULL,
  `allowed` enum('Y','N') COLLATE utf8_bin NOT NULL DEFAULT 'N'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='User groups with configured menu items';

-- --------------------------------------------------------

--
-- Table structure for table `pma__users`
--

CREATE TABLE `pma__users` (
  `username` varchar(64) COLLATE utf8_bin NOT NULL,
  `usergroup` varchar(64) COLLATE utf8_bin NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin COMMENT='Users and their assignments to user groups';

--
-- Indexes for dumped tables
--

--
-- Indexes for table `pma__bookmark`
--
ALTER TABLE `pma__bookmark`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `pma__central_columns`
--
ALTER TABLE `pma__central_columns`
  ADD PRIMARY KEY (`db_name`,`col_name`);

--
-- Indexes for table `pma__column_info`
--
ALTER TABLE `pma__column_info`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `db_name` (`db_name`,`table_name`,`column_name`);

--
-- Indexes for table `pma__designer_settings`
--
ALTER TABLE `pma__designer_settings`
  ADD PRIMARY KEY (`username`);

--
-- Indexes for table `pma__export_templates`
--
ALTER TABLE `pma__export_templates`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `u_user_type_template` (`username`,`export_type`,`template_name`);

--
-- Indexes for table `pma__favorite`
--
ALTER TABLE `pma__favorite`
  ADD PRIMARY KEY (`username`);

--
-- Indexes for table `pma__history`
--
ALTER TABLE `pma__history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `username` (`username`,`db`,`table`,`timevalue`);

--
-- Indexes for table `pma__navigationhiding`
--
ALTER TABLE `pma__navigationhiding`
  ADD PRIMARY KEY (`username`,`item_name`,`item_type`,`db_name`,`table_name`);

--
-- Indexes for table `pma__pdf_pages`
--
ALTER TABLE `pma__pdf_pages`
  ADD PRIMARY KEY (`page_nr`),
  ADD KEY `db_name` (`db_name`);

--
-- Indexes for table `pma__recent`
--
ALTER TABLE `pma__recent`
  ADD PRIMARY KEY (`username`);

--
-- Indexes for table `pma__relation`
--
ALTER TABLE `pma__relation`
  ADD PRIMARY KEY (`master_db`,`master_table`,`master_field`),
  ADD KEY `foreign_field` (`foreign_db`,`foreign_table`);

--
-- Indexes for table `pma__savedsearches`
--
ALTER TABLE `pma__savedsearches`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `u_savedsearches_username_dbname` (`username`,`db_name`,`search_name`);

--
-- Indexes for table `pma__table_coords`
--
ALTER TABLE `pma__table_coords`
  ADD PRIMARY KEY (`db_name`,`table_name`,`pdf_page_number`);

--
-- Indexes for table `pma__table_info`
--
ALTER TABLE `pma__table_info`
  ADD PRIMARY KEY (`db_name`,`table_name`);

--
-- Indexes for table `pma__table_uiprefs`
--
ALTER TABLE `pma__table_uiprefs`
  ADD PRIMARY KEY (`username`,`db_name`,`table_name`);

--
-- Indexes for table `pma__tracking`
--
ALTER TABLE `pma__tracking`
  ADD PRIMARY KEY (`db_name`,`table_name`,`version`);

--
-- Indexes for table `pma__userconfig`
--
ALTER TABLE `pma__userconfig`
  ADD PRIMARY KEY (`username`);

--
-- Indexes for table `pma__usergroups`
--
ALTER TABLE `pma__usergroups`
  ADD PRIMARY KEY (`usergroup`,`tab`,`allowed`);

--
-- Indexes for table `pma__users`
--
ALTER TABLE `pma__users`
  ADD PRIMARY KEY (`username`,`usergroup`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `pma__bookmark`
--
ALTER TABLE `pma__bookmark`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `pma__column_info`
--
ALTER TABLE `pma__column_info`
  MODIFY `id` int(5) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `pma__export_templates`
--
ALTER TABLE `pma__export_templates`
  MODIFY `id` int(5) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `pma__history`
--
ALTER TABLE `pma__history`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `pma__pdf_pages`
--
ALTER TABLE `pma__pdf_pages`
  MODIFY `page_nr` int(10) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `pma__savedsearches`
--
ALTER TABLE `pma__savedsearches`
  MODIFY `id` int(5) UNSIGNED NOT NULL AUTO_INCREMENT;
COMMIT;