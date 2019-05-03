CREATE DATABASE  IF NOT EXISTS `northernairport` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `northernairport`;
-- MySQL dump 10.13  Distrib 8.0.16, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: northernairport
-- ------------------------------------------------------
-- Server version	8.0.16

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `accountdetails`
--

DROP TABLE IF EXISTS `accountdetails`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `accountdetails` (
  `AccountDetailID` int(11) NOT NULL AUTO_INCREMENT,
  `ClientID` int(11) NOT NULL,
  `Username` varchar(25) NOT NULL,
  `Password` varchar(100) NOT NULL,
  `RoleID` int(11) NOT NULL,
  PRIMARY KEY (`AccountDetailID`),
  KEY `FK_148` (`ClientID`),
  KEY `FK_266` (`RoleID`),
  CONSTRAINT `FK_148` FOREIGN KEY (`ClientID`) REFERENCES `clients` (`ClientID`),
  CONSTRAINT `FK_266` FOREIGN KEY (`RoleID`) REFERENCES `roles` (`RoleID`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `accountdetails`
--

LOCK TABLES `accountdetails` WRITE;
/*!40000 ALTER TABLE `accountdetails` DISABLE KEYS */;
INSERT INTO `accountdetails` VALUES 
(1,3,'test','$2a$08$wy3Vf/Bhb27MVa7ZnVXqwujqxKkUtQ/rfHf9aryf84Tzr5zwXNa7q',1);

/*!40000 ALTER TABLE `accountdetails` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `airlines`
--

DROP TABLE IF EXISTS `airlines`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `airlines` (
  `AirlineID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(50) NOT NULL,
  `Terminal` int(11) NOT NULL,
  PRIMARY KEY (`AirlineID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `airlines`
--

LOCK TABLES `airlines` WRITE;
/*!40000 ALTER TABLE `airlines` DISABLE KEYS */;
INSERT INTO `airlines` VALUES (1,'Lufthansa',1),(2,'Alitalia',3),(3,'WestJet',3),(4,'Sunwing Airlines',3),(5,'KLM',3),(6,'Air France',3),(7,'Air Transat',3),(8,'Delta',3),(9,'EVA Air',1),(10,'Etihad Airways',3),(11,'Hainan Airlines',3),(12,'CanJet',3),(13,'LOT Polish Airlines',1),(14,'All Nippon Airways',1),(15,'American Airlines',3),(16,'Air New Zealand',1),(17,'Alaska Airlines',3),(18,'Arkefly',1),(19,'Air China',1),(20,'Caribbean Airlines',3),(21,'Copa Airlines',1),(22,'Condor',3),(23,'All Nippon Air',1),(24,'EL AL',3),(25,'Emirates',1),(26,'Finnair',3),(27,'Japan Airlines',3),(28,'Air Canada',1),(29,'Asiana Airlines',1),(30,'Air Wisconsin',1),(31,'Ethiopian Airlines',1),(32,'Aeroflot',3),(33,'Austrian',1),(34,'Avianca',1),(35,'British Airways',3),(36,'Air Jamaica',1),(37,'Cathay Pacific',3),(38,'Cubana',3),(39,'Comar',3),(40,'Fly Jamaica',3),(41,'Ibena',3),(42,'Icelandair',3),(43,'Jet Airways',1),(44,'Air India',1),(45,'Korean Air',3),(46,'LAN Airlines',3),(47,'Aerosvit Ukranian Air',3),(48,'Miami Air International',3),(49,'Pakistan International',3),(50,'Phillippine Airlines',3),(51,'Qantas',3),(52,'SAS Scandinavian',1),(53,'SATA International',3),(54,'SAUDIA',3),(55,'Singapore Airlines',1),(56,'Swiss International Air Lines',1),(57,'TACA',1),(58,'Thai Airways International',1),(59,'Transaero Airlines',3),(60,'Turkish Airlines',1),(61,'Aeromexico',3),(62,'United',1),(63,'US Airways',3),(64,'Commuter',3),(65,'Continental air',1),(66,'Continental Express',1),(67,'Jet Airways',3),(68,'WOW Air',3),(69,'Egypt Air',1),(70,'Aer LIngus',3),(71,'China Eastern',3),(72,'China Eastern Airlines',3),(73,'Brazilian Airlines',3),(74,'American Eagle',3),(75,'Union Pearson Express',1),(76,'Brussels',1),(77,'A Terminal 1  Not flying,',1),(78,'A Terminal 3, Not flying',3),(79,'China South Airlines',3),(80,'TAP Portugal Air',1),(81,'Flair Airline',3);
/*!40000 ALTER TABLE `airlines` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cities`
--

DROP TABLE IF EXISTS `cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `cities` (
  `CityID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(50) NOT NULL,
  PRIMARY KEY (`CityID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cities`
--

LOCK TABLES `cities` WRITE;
/*!40000 ALTER TABLE `cities` DISABLE KEYS */;
INSERT INTO `cities` VALUES (1,'North Bay'),(2,'Toronto Pearson and Hotels'),(3,'Bracebridge'),(4,'Burks Falls'),(5,'Callander'),(7,'Emsdale'),(8,'Gravenhurst'),(9,'Huntsville Deerhurst Resort'),(10,'Huntsville'),(11,'Katrine'),(14,'Novar'),(15,'Port Sydney'),(16,'Powassan'),(17,'South River'),(18,'Sundridge'),(20,'Trout Creek');
/*!40000 ALTER TABLE `cities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cityoffsets`
--

DROP TABLE IF EXISTS `cityoffsets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `cityoffsets` (
  `CityOffsetID` int(11) NOT NULL,
  `CityID` int(11) NOT NULL,
  `NorthOffset` int(11) NOT NULL,
  `SouthOffset` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cityoffsets`
--

LOCK TABLES `cityoffsets` WRITE;
/*!40000 ALTER TABLE `cityoffsets` DISABLE KEYS */;
INSERT INTO `cityoffsets` VALUES (1,1,0,240),(2,2,240,0),(3,3,115,125),(4,4,55,185),(5,5,10,230),(6,7,65,175),(7,8,135,105),(8,10,85,155),(9,11,60,180),(10,14,75,165),(11,15,90,150),(12,16,25,215),(13,17,40,200),(14,18,45,195),(15,20,30,210);
/*!40000 ALTER TABLE `cityoffsets` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clients`
--

DROP TABLE IF EXISTS `clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `clients` (
  `ClientID` int(11) NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(50) NOT NULL,
  `LastName` varchar(50) NOT NULL,
  `Phone` bigint(20) NOT NULL,
  `Email` varchar(100) NOT NULL,
  `StreetAddress` varchar(100) NOT NULL,
  `City` varchar(50) NOT NULL,
  `Province` varchar(50) NOT NULL,
  `PostalCode` char(7) NOT NULL,
  `Country` varchar(50) NOT NULL,
  `DefaultDepartureAddress` varchar(100) DEFAULT NULL,
  `DefaultDepartureCityID` int(11) DEFAULT NULL,
  `TravelAgentID` int(11) DEFAULT NULL,
  `Notes` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`ClientID`),
  KEY `FK_171` (`TravelAgentID`),
  KEY `FK_217` (`DefaultDepartureCityID`),
  CONSTRAINT `FK_171` FOREIGN KEY (`TravelAgentID`) REFERENCES `travelagents` (`TravelAgentID`),
  CONSTRAINT `FK_217` FOREIGN KEY (`DefaultDepartureCityID`) REFERENCES `cities` (`CityID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clients`
--

LOCK TABLES `clients` WRITE;
/*!40000 ALTER TABLE `clients` DISABLE KEYS */;
INSERT INTO `clients` VALUES (3,'test','test',1234567890,'test@test.com','test','test','ontario','p1b8p4','canada',NULL,1,NULL,NULL);
/*!40000 ALTER TABLE `clients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `customertypes`
--

DROP TABLE IF EXISTS `customertypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `customertypes` (
  `CustomerTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) NOT NULL,
  PRIMARY KEY (`CustomerTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customertypes`
--

LOCK TABLES `customertypes` WRITE;
/*!40000 ALTER TABLE `customertypes` DISABLE KEYS */;
INSERT INTO `customertypes` VALUES (1,'Adult'),(2,'Senior'),(3,'Student'),(4,'Child');
/*!40000 ALTER TABLE `customertypes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `customvenues`
--

DROP TABLE IF EXISTS `customvenues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `customvenues` (
  `CustomVenueID` int(11) NOT NULL AUTO_INCREMENT,
  `Address` varchar(255) NOT NULL,
  `ClientId` int(11) NOT NULL,
  PRIMARY KEY (`CustomVenueID`),
  KEY `FK_327` (`ClientId`),
  CONSTRAINT `FK_327` FOREIGN KEY (`ClientId`) REFERENCES `clients` (`ClientID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customvenues`
--

LOCK TABLES `customvenues` WRITE;
/*!40000 ALTER TABLE `customvenues` DISABLE KEYS */;
/*!40000 ALTER TABLE `customvenues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `departuretimes`
--

DROP TABLE IF EXISTS `departuretimes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `departuretimes` (
  `DepartureTimeID` int(11) NOT NULL AUTO_INCREMENT,
  `CityID` int(11) NOT NULL,
  `DepartureTime` int(11) NOT NULL,
  `Recurring` int(11) NOT NULL,
  `StartDate` date DEFAULT NULL,
  `EndDate` date DEFAULT NULL,
  PRIMARY KEY (`DepartureTimeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `departuretimes`
--

LOCK TABLES `departuretimes` WRITE;
/*!40000 ALTER TABLE `departuretimes` DISABLE KEYS */;
INSERT INTO `departuretimes` VALUES (1,1,630,1,NULL,NULL),(2,1,1120,1,NULL,NULL),(3,1,1200,1,NULL,NULL),(4,1,2240,1,NULL,NULL),(5,2,1330,1,NULL,NULL),(6,2,1130,1,NULL,NULL),(7,2,1630,1,NULL,NULL),(8,2,1800,1,NULL,NULL);
/*!40000 ALTER TABLE `departuretimes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `discountcodes`
--

DROP TABLE IF EXISTS `discountcodes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `discountcodes` (
  `DiscountCodeID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) NOT NULL,
  `Percentage` int(11) NOT NULL,
  `Amount` int(11) NOT NULL,
  `StartDate` date NOT NULL,
  `EndDate` date NOT NULL,
  PRIMARY KEY (`DiscountCodeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `discountcodes`
--

LOCK TABLES `discountcodes` WRITE;
/*!40000 ALTER TABLE `discountcodes` DISABLE KEYS */;
INSERT INTO `discountcodes` VALUES (1,'test code',10,0,'2000-01-01','2020-01-01');
/*!40000 ALTER TABLE `discountcodes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `drivers`
--

DROP TABLE IF EXISTS `drivers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `drivers` (
  `DriverID` int(11) NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(50) NOT NULL,
  `LastName` varchar(50) NOT NULL,
  PRIMARY KEY (`DriverID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `drivers`
--

LOCK TABLES `drivers` WRITE;
/*!40000 ALTER TABLE `drivers` DISABLE KEYS */;
INSERT INTO `drivers` VALUES (1,'Andy','Duggan'),(2,'Joe','Palangio'),(3,'Mike','Lok'),(4,'Blake','Gennoe'),(5,'Rob','Farris'),(6,'Carole','Tran'),(7,'Mike','Ianiro'),(8,'Cosmo','Ianiro'),(9,'Mike','Brideau'),(10,'Brad','Openshaw'),(11,'Jean Paul','Chartrand');
/*!40000 ALTER TABLE `drivers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prices`
--

DROP TABLE IF EXISTS `prices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `prices` (
  `PriceID` int(11) NOT NULL AUTO_INCREMENT,
  `CustomerTypeID` int(11) NOT NULL,
  `Price` float NOT NULL,
  `ReturnPrice` float NOT NULL,
  `DepartureCityID` int(11) NOT NULL,
  `DestinationCityID` int(11) NOT NULL,
  PRIMARY KEY (`PriceID`),
  KEY `FK_232` (`CustomerTypeID`),
  KEY `FK_236` (`DepartureCityID`),
  KEY `FK_239` (`DestinationCityID`),
  CONSTRAINT `FK_232` FOREIGN KEY (`CustomerTypeID`) REFERENCES `customertypes` (`CustomerTypeID`),
  CONSTRAINT `FK_236` FOREIGN KEY (`DepartureCityID`) REFERENCES `cities` (`CityID`),
  CONSTRAINT `FK_239` FOREIGN KEY (`DestinationCityID`) REFERENCES `cities` (`CityID`)
) ENGINE=InnoDB AUTO_INCREMENT=152 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prices`
--

LOCK TABLES `prices` WRITE;
/*!40000 ALTER TABLE `prices` DISABLE KEYS */;
INSERT INTO `prices` VALUES (1,1,119,0,2,5),(2,2,111,0,2,5),(3,4,49.5,0,17,2),(4,4,59.5,0,2,5),(5,4,49.5,0,18,2),(6,4,47,0,3,2),(7,3,109,0,2,5),(8,4,49.5,0,2,4),(9,1,99,0,2,4),(10,2,95,0,2,4),(11,3,95,0,2,4),(12,1,119,0,2,16),(13,2,111,0,2,16),(14,3,109,0,2,16),(15,4,59.5,0,2,16),(16,1,99,0,2,20),(17,4,49.5,0,11,2),(18,4,49.5,0,14,2),(19,4,49.5,0,2,20),(20,1,99,0,2,17),(21,2,95,0,4,2),(22,4,47,0,10,2),(23,1,94,0,10,2),(24,3,95,0,7,2),(25,2,95,0,2,17),(26,3,95,0,2,17),(27,4,49.5,0,2,17),(28,4,49.5,0,2,18),(29,1,99,0,20,2),(30,4,47.5,0,20,2),(31,1,119,0,1,2),(32,2,111,0,1,2),(33,3,109,0,1,2),(34,4,59.5,0,1,2),(35,1,119,0,5,2),(36,2,111,0,5,2),(37,3,109,0,5,2),(38,4,59.5,0,5,2),(39,1,119,0,16,2),(40,4,59.5,0,16,2),(41,3,109,0,16,2),(42,2,111,0,16,2),(43,1,94,0,2,10),(44,1,94,0,2,3),(45,3,95,169,20,2),(46,2,95,169,20,2),(47,1,99,0,17,2),(48,2,89,0,10,2),(49,3,95,169,17,2),(50,2,95,169,17,2),(51,1,99,0,18,2),(52,3,89,0,10,2),(53,3,95,169,18,2),(54,2,95,169,18,2),(55,1,99,0,4,2),(56,1,99,0,2,11),(57,3,95,0,4,2),(58,4,49.5,0,4,2),(59,1,99,0,11,2),(60,1,99,0,2,7),(61,3,95,0,11,2),(62,2,95,0,11,2),(63,1,99,0,7,2),(64,2,95,0,2,7),(65,2,95,0,7,2),(66,4,49.5,0,7,2),(67,1,99,0,14,2),(68,3,95,0,2,7),(69,2,95,0,14,2),(70,3,95,0,14,2),(71,4,49.5,0,2,7),(72,1,99,0,2,14),(73,2,95,0,2,11),(74,2,95,0,2,14),(75,1,94,0,3,2),(76,3,95,0,2,14),(77,2,89,0,3,2),(78,3,89,0,3,2),(79,1,94,0,15,2),(80,4,47,0,15,2),(81,3,89,159,15,2),(82,2,89,159,15,2),(83,1,94,0,8,2),(84,2,89,0,8,2),(85,4,47,0,8,2),(86,3,89,0,8,2),(87,4,49.5,0,2,14),(88,3,95,0,2,11),(89,4,47.5,0,2,11),(90,1,119,0,2,1),(91,2,111,0,2,1),(92,3,109,0,2,1),(93,2,89,0,2,10),(94,3,89,0,2,10),(95,4,47,0,2,10),(96,1,94,0,2,15),(97,4,47,0,2,15),(98,4,47,0,2,3),(99,1,94,0,2,8),(100,2,89,0,2,8),(101,3,89,0,2,8),(102,4,59.5,0,2,1),(103,4,47,0,2,8),(104,1,124,0,9,2),(105,2,119,0,9,2),(106,3,119,0,9,2),(107,4,47,0,9,2),(115,2,89,0,2,3),(119,3,89,0,2,3),(121,1,99,0,2,18),(132,3,89,159,2,15),(133,2,89,159,2,15),(144,2,95,169,2,18),(145,3,95,169,2,18),(150,2,95,169,2,20),(151,3,95,169,2,20);
/*!40000 ALTER TABLE `prices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reservations`
--

DROP TABLE IF EXISTS `reservations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `reservations` (
  `ReservationID` int(11) NOT NULL AUTO_INCREMENT,
  `ClientID` int(11) NOT NULL,
  `DepartureCityID` int(11) NOT NULL,
  `DepartureVenueID` int(11) NOT NULL,
  `DepartureTimeID` int(11) NOT NULL,
  `DestinationCityID` int(11) NOT NULL,
  `DestinationVenueID` int(11) NOT NULL,
  `ReturnDepartureCityID` int(11) DEFAULT NULL,
  `ReturnDepartureVenueID` int(11) DEFAULT NULL,
  `ReturnDepartureTimeID` int(11) DEFAULT NULL,
  `ReturnDestinationCityID` int(11) DEFAULT NULL,
  `ReturnDestinationVenueID` int(11) DEFAULT NULL,
  `DiscountCodeID` int(11) DEFAULT '0',
  `DepartureAirlineID` int(11) NOT NULL,
  `ReturnAirlineID` int(11) DEFAULT NULL,
  `DriverNotes` text,
  `InternalNotes` text,
  `DepartureNumAdults` int(11) NOT NULL,
  `DepartureNumStudents` int(11) NOT NULL,
  `DepartureNumSeniors` int(11) NOT NULL,
  `DepartureNumChildren` int(11) NOT NULL,
  `ReturnNumAdults` int(11) DEFAULT NULL,
  `ReturnNumStudents` int(11) DEFAULT NULL,
  `ReturnNumSeniors` int(11) DEFAULT NULL,
  `ReturnNumChildren` int(11) DEFAULT NULL,
  `Price` float NOT NULL,
  `Status` varchar(25) DEFAULT NULL,
  `Hash` varchar(255) DEFAULT NULL,
  `CustomDepartureID` int(11) DEFAULT NULL,
  `CustomDestinationID` int(11) DEFAULT NULL,
  `DepartureDate` date NOT NULL,
  `ReturnDate` date DEFAULT NULL,
  `TripTypeID` int(11) NOT NULL,
  `TripID` int(11) NOT NULL,
  `BalanceOwing` float DEFAULT NULL,
  `ElavonTransactionID` int(11) DEFAULT NULL,
  PRIMARY KEY (`ReservationID`),
  KEY `FK_163` (`DepartureCityID`),
  KEY `FK_181` (`DiscountCodeID`),
  KEY `FK_195` (`DepartureTimeID`),
  KEY `FK_203` (`DestinationCityID`),
  KEY `FK_211` (`DepartureVenueID`),
  KEY `FK_214` (`DestinationVenueID`),
  KEY `FK_273` (`DepartureAirlineID`),
  KEY `FK_311` (`ReturnAirlineID`),
  KEY `FK_330` (`CustomDepartureID`),
  KEY `FK_333` (`CustomDestinationID`),
  KEY `FK_341` (`TripTypeID`),
  KEY `FK_342` (`TripID`),
  CONSTRAINT `FK_163` FOREIGN KEY (`DepartureCityID`) REFERENCES `cities` (`CityID`),
  CONSTRAINT `FK_195` FOREIGN KEY (`DepartureTimeID`) REFERENCES `departuretimes` (`DepartureTimeID`),
  CONSTRAINT `FK_203` FOREIGN KEY (`DestinationCityID`) REFERENCES `cities` (`CityID`),
  CONSTRAINT `FK_211` FOREIGN KEY (`DepartureVenueID`) REFERENCES `venues` (`VenueID`),
  CONSTRAINT `FK_214` FOREIGN KEY (`DestinationVenueID`) REFERENCES `venues` (`VenueID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservations`
--

LOCK TABLES `reservations` WRITE;
/*!40000 ALTER TABLE `reservations` DISABLE KEYS */;
INSERT INTO `reservations` VALUES (1,3,1,48,1,2,16,NULL,NULL,NULL,NULL,NULL,1,1,NULL,'dnotes','inotes',0,0,0,0,NULL,NULL,NULL,NULL,120,'','',0,0,'2019-04-30',NULL,0,1,0,0);
/*!40000 ALTER TABLE `reservations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reservationtypes`
--

DROP TABLE IF EXISTS `reservationtypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `reservationtypes` (
  `ReservationTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) NOT NULL,
  PRIMARY KEY (`ReservationTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reservationtypes`
--

LOCK TABLES `reservationtypes` WRITE;
/*!40000 ALTER TABLE `reservationtypes` DISABLE KEYS */;
INSERT INTO `reservationtypes` VALUES (1,'one way'),(2,'round trip');
/*!40000 ALTER TABLE `reservationtypes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `roles` (
  `RoleID` int(11) NOT NULL AUTO_INCREMENT,
  `RoleName` varchar(25) NOT NULL,
  PRIMARY KEY (`RoleID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'test'),(2,'client'),(3,'staff'),(4,'admin'),(5,'driver');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `taxes`
--

DROP TABLE IF EXISTS `taxes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `taxes` (
  `TaxID` int(11) NOT NULL AUTO_INCREMENT,
  `Percentage` int(11) NOT NULL,
  `Name` varchar(25) NOT NULL,
  `Active` int(11) NOT NULL,
  PRIMARY KEY (`TaxID`,`Percentage`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `taxes`
--

LOCK TABLES `taxes` WRITE;
/*!40000 ALTER TABLE `taxes` DISABLE KEYS */;
INSERT INTO `taxes` VALUES (1,13,'HST',1);
/*!40000 ALTER TABLE `taxes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `transactions` (
  `TransactionID` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(15) NOT NULL,
  `Created` datetime NOT NULL,
  `ReservationID` int(11) NOT NULL,
  `Response` varchar(255) NOT NULL,
  PRIMARY KEY (`TransactionID`),
  KEY `FK_301` (`ReservationID`),
  CONSTRAINT `FK_301` FOREIGN KEY (`ReservationID`) REFERENCES `reservations` (`ReservationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `travelagencies`
--

DROP TABLE IF EXISTS `travelagencies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `travelagencies` (
  `TravelAgencyID` int(11) NOT NULL AUTO_INCREMENT,
  `TravelAgencyName` varchar(50) NOT NULL,
  PRIMARY KEY (`TravelAgencyID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `travelagencies`
--

LOCK TABLES `travelagencies` WRITE;
/*!40000 ALTER TABLE `travelagencies` DISABLE KEYS */;
/*!40000 ALTER TABLE `travelagencies` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `travelagents`
--

DROP TABLE IF EXISTS `travelagents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `travelagents` (
  `TravelAgentID` int(11) NOT NULL AUTO_INCREMENT,
  `TravelAgentName` varchar(100) NOT NULL,
  `IATANumber` varchar(100) NOT NULL,
  `TravelAgencyID` int(11) NOT NULL,
  PRIMARY KEY (`TravelAgentID`),
  KEY `FK_354` (`TravelAgencyID`),
  CONSTRAINT `FK_354` FOREIGN KEY (`TravelAgencyID`) REFERENCES `travelagencies` (`TravelAgencyID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `travelagents`
--

LOCK TABLES `travelagents` WRITE;
/*!40000 ALTER TABLE `travelagents` DISABLE KEYS */;
/*!40000 ALTER TABLE `travelagents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `trips`
--

DROP TABLE IF EXISTS `trips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `trips` (
  `TripID` int(11) NOT NULL AUTO_INCREMENT,
  `DepartureDate` date NOT NULL,
  `DepartureTimeID` int(11) NOT NULL,
  `NumPassengers` int(11) DEFAULT '0',
  `DriverID` int(11) DEFAULT '0',
  `VehicleID` int(11) DEFAULT '0',
  `OmitTrip` int(11) DEFAULT '0',
  `Postpone` int(11) DEFAULT '0',
  `RescheduleDate` date DEFAULT NULL,
  `RescheduleTime` int(11) DEFAULT '0',
  `Cancelled` int(11) DEFAULT '0',
  PRIMARY KEY (`TripID`),
  KEY `FK_254` (`DepartureTimeID`),
  KEY `FK_288` (`DriverID`),
  KEY `FK_291` (`VehicleID`),
  CONSTRAINT `FK_254` FOREIGN KEY (`DepartureTimeID`) REFERENCES `departuretimes` (`DepartureTimeID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `trips`
--

LOCK TABLES `trips` WRITE;
/*!40000 ALTER TABLE `trips` DISABLE KEYS */;
INSERT INTO `trips` VALUES (1,'2019-04-30',1,0,0,0,0,0,NULL,0,0);
/*!40000 ALTER TABLE `trips` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `triptypes`
--

DROP TABLE IF EXISTS `triptypes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `triptypes` (
  `TripTypeID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(25) NOT NULL,
  PRIMARY KEY (`TripTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `triptypes`
--

LOCK TABLES `triptypes` WRITE;
/*!40000 ALTER TABLE `triptypes` DISABLE KEYS */;
INSERT INTO `triptypes` VALUES (1,'one way'),(2,'return');
/*!40000 ALTER TABLE `triptypes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vehicles`
--

DROP TABLE IF EXISTS `vehicles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `vehicles` (
  `VehicleID` int(11) NOT NULL AUTO_INCREMENT,
  `LicensePlate` varchar(25) NOT NULL,
  `NumSeats` int(11) NOT NULL,
  `Make` varchar(25) NOT NULL,
  PRIMARY KEY (`VehicleID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vehicles`
--

LOCK TABLES `vehicles` WRITE;
/*!40000 ALTER TABLE `vehicles` DISABLE KEYS */;
INSERT INTO `vehicles` VALUES (1,'1110 lic. 7405BF',10,'Mercedes'),(2,'1109 lic. 8362BF',10,'Mercedes'),(3,'1512 lic. 8697BH',14,'Mercedes'),(4,'1411 lic. 8691BH',11,'Mercedes');
/*!40000 ALTER TABLE `vehicles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `venues`
--

DROP TABLE IF EXISTS `venues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `venues` (
  `VenueID` int(11) NOT NULL AUTO_INCREMENT,
  `CityID` int(11) NOT NULL,
  `Name` varchar(100) NOT NULL,
  `ExtraCost` float DEFAULT NULL,
  `Active` int(11) NOT NULL,
  `ExtraTime` int(11) NOT NULL,
  PRIMARY KEY (`VenueID`),
  KEY `FK_206` (`CityID`),
  CONSTRAINT `FK_206` FOREIGN KEY (`CityID`) REFERENCES `cities` (`CityID`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `venues`
--

LOCK TABLES `venues` WRITE;
/*!40000 ALTER TABLE `venues` DISABLE KEYS */;
INSERT INTO `venues` VALUES (1,8,'McDonald\'s Restaurant located at Bethune Dr. and Muskoka Rd',0,1,0),(2,3,'Tim Hortons coffee shop corner of Depot Rd. and Taylor Rd',0,1,0),(3,4,'Tim Horton\'s located at 27 Commercial Dr',0,1,0),(4,5,'Wasi Petro Canada located at the corner of Wasi Rd and Callander Bay Dr',0,1,0),(5,7,'Tourist information centre located on Emsdale Rd. Exit #244',0,1,0),(6,10,'Holiday Inn Express located at Howland Dr',0,1,0),(7,11,'The Lucky Dollar on old Hwy 11',0,1,0),(8,1,'Country Style coffee shop located at 1401 Seymour St. and Hwy 11',0,1,0),(9,14,'Foodland',0,1,0),(10,15,'Smiths Resaurant located on highway 11',0,1,0),(11,16,'East Himsworth Cafe, previous Ethels',0,1,0),(12,17,'Shell gas station highway 124',0,1,0),(13,18,'Blue Roof restaurant',0,1,0),(14,20,'TJ\'s Variety on highway 522',0,1,0),(15,2,'Nu Hotel-6465 Airport Rd.',0,1,0),(16,2,'Quality Inn 2180 Islington',0,1,0),(17,2,'YYZ No Flight TERMINAL 1,',0,0,0),(18,2,'YYZ No Flight  Terminal 3,',0,0,0),(19,2,'UP Express train Pearson Terminal 1',0,1,0),(25,2,'Sheraton Airport, 801 Dixon Rd',0,1,0),(26,2,'Comfort Inn, 6355 Airport Rd.',0,1,0),(27,2,'Carlingview, 221 Carlingview Dr.',0,1,0),(28,2,'Crown Plaza, 33 Carlson Crt.',0,1,0),(29,2,'Delta Hotels 655 Dixon Rd',0,1,0),(30,2,'Fairfield Inn, 3299 Caroga Dr.',0,1,0),(31,2,'Four Points Sheraton, 6257 Airport Rd.',0,1,0),(32,2,'Hampton Inn, 3279 Caroga Dr.',0,1,0),(33,2,'Hilton, 5875 Airport Rd',0,1,0),(34,2,'Holiday Inn , 970 Dixon Rd',0,1,0),(35,2,'Holiday Inn, 600 Dixon Rd.',0,1,0),(36,2,'Embassy Suites, 262 Carlingview Dr.',0,1,0),(37,2,'Hilton Garden Inn,  3311 Caroga Dr.',0,1,0),(38,2,'Radisson Suite, 640 Dixon Rd',0,1,0),(39,2,'Residence Inn Marriott, 17 Reading Crt.',0,1,0),(40,2,'Sheraton Gateway, PIA terminal 3',0,1,0),(41,2,'Marriott Toronto Airport, 901 Dixon Rd.',0,1,0),(42,2,'Double Tree,  925 Dixon Rd.',0,1,0),(43,2,'The Westin Toronto Airport, 950 Dixon Rd.',0,1,0),(44,2,'Best Western Premier, 135 Carlingview Dr.',0,1,0),(45,2,'Sandman Signature, 55 Reading Crt.',0,1,0),(46,2,'Alt Hotel,  6080 Viscount Rd.',0,1,0),(47,2,'Courtyard by Marriot-231 Carlingview',0,1,0),(48,1,'North Bay Office 191 Booth Rd #7',0,1,0),(49,10,'Deerhurst Resort',30,1,0),(99,2,'AIRPORT AIRLINES: CLICK HERE FOR AIRLINE LIST',0,1,0),(100,1,'Home Pickup or Dropoff',15,1,30);
/*!40000 ALTER TABLE `venues` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'northernairport'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-04-28 21:38:05
