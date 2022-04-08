package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//Store offers an interface for various db function to the rest of the application
type Store interface {
	CreateUser(client *Client) (int, error)
	SignInUser(client *Client) error
	CreateReservation(reservation *Reservation) error
	GetClientInfo(client *Client) error
	GetDriverIDFromClient(clientid int) (int, error)
	GetVenues() []Venues
	GetVenueName(int) string
	GetVenueCount() int
	AddVenue(venue Venues) error
	UpdateVenue(venue *Venues) error
	DeleteVenue(venue int) error
	GetCities() []Cities
	GetCityCount() int
	GetCityID(string) int
	GetCityName(int) string
	AddCity(city Cities) error
	UpdateCity(city *Cities) error
	DeleteCity(city int) error
	GetDepartureTimesCount() int
	GetDepartureTimes() []DepartureTimes
	AddDepartureTime(departuretime *DepartureTimes) error
	UpdateDepartureTime(departuretime *DepartureTimes) error
	GetAirlines() []Airlines
	GetAirlineCount() int
	GetOrAddTrip(reservation *Reservation) error
	GetTrips() []Trips
	UpdateTrip(trip *Trips) error
	OmitTrip(trip *Trips) error
	SearchReservations(name string, phone int, email string) []SearchReservations
	SearchTrips(tripdate time.Time, reportType int) []Trips
	PostponeReservation(searchreservation *SearchReservations) error
	CancelReservation(searchreservation *SearchReservations) error
	GetDrivers(driverid int) []Drivers
	AddDriver(driver Drivers) error
	UpdateDriver(driver *Drivers) error
	DeleteDriver(driver int) error
	GetVehicles() []Vehicles
	AddVehicle(vehicle Vehicles) error
	UpdateVehicle(vehicle *Vehicles) error
	DeleteVehicle(vehicle int) error
	GetPrice(departurecityid int, destinationcityid int, retdeparturecityid int, retdestinationcityid int, customertypeid int, triptypeid int, discountcode string) float32
	GetPrices() []Prices
	UpdatePrice(price *Prices) error
	GetDiscount(discountcode *DiscountCode) error
	AddVenueFee(departurevenueid int, destinationvenueid int, retdeparturevenueid int, retdestinationvenueid int) float32
	DriverReservations(driverid int, departuredate time.Time) DriverReportForm
	GetTravelAgencyCount() int
	GetTravelAgencies() []TravelAgencies
	GetClientFromReservation(reservationid int) Client
	TravelAgencyReports(month int, year int) []TravelAgencyReport
	AGTAQueryReport(startdate time.Time, enddate time.Time) []AGTAReport
	//UpdateStatus(reservationid int, elavontransactionid int, status string) error
	UpdateStatus(reservationid int, elavontransactionid string, status string) error
	GetReservationAmount(reservationid int) float32
	MigrateDB()
}

//The `dbStore` struct will implement the `Store` interface it also takes the sql
//DB connection object, which represents the database connection.
type dbStore struct {
	db *sql.DB
}

//CreateUser - store new client in database
func (store *dbStore) CreateUser(client *Client) (int, error) {

	result, err := store.db.Exec("INSERT INTO clients(firstname, lastname, phone, email, "+
		"streetaddress, city, province, postalcode, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		client.Firstname, client.Lastname, client.Phone, client.Email, client.StreetAddress,
		client.City, client.Province, client.PostalCode, client.Country)

	if err != nil {
		return 0, err
	}

	//get id from client insertion transaction
	id, _ := result.LastInsertId()

	//create account details record linked to client record
	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) "+
		"VALUES (?, ?, ?, ?)",
		id, client.Password, client.RoleID, client.Username)

	return int(id), err
}

//CreateReservation - store new reservation in database
func (store *dbStore) CreateReservation(reservation *Reservation) error {

	//temporarily hardcoding some values until they are either implemented or removed
	reservation.DepartureAirlineID = 1
	reservation.ReturnAirlineID = 1
	reservation.Status = ""
	reservation.Hash = ""
	reservation.ElavonTransactionID = 0

	var departuredetails string
	var returndetails string

	departuredetails = ""
	returndetails = ""

	log.Print("Inserting reservation into DB")

	log.Printf("TripID: %d", reservation.TripID)
	log.Printf("ClientID: %d", reservation.ClientDetails.ClientID)

	var insertError error
	var id int64

	/*****************************************************
	*	must ensure trip is not omitted!
	******************************************************/

	if reservation.TripTypeID == 2 {
		result, err := store.db.Exec("INSERT INTO reservations("+
			"clientid, departurecityid, departurevenueid, departuretimeid, "+
			"destinationcityid, destinationvenueid, returndeparturecityid, returndeparturevenueid, returndeparturetimeid, "+
			"returndestinationcityid, returndestinationvenueid, discountcodeid, departureairlineid, returnairlineid, "+
			"drivernotes, internalnotes, departurenumadults, departurenumstudents, departurenumseniors, "+
			"departurenumchildren, returnnumadults, returnnumstudents, returnnumseniors, returnnumchildren, "+
			"price, status, hash, customdepartureid, customdestinationid, "+
			"departuredate, returndate, triptypeid, tripid, returntrip, balanceowing, elavontransactionid, flightnumber, flighttime) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?)",
			reservation.ClientDetails.ClientID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.ReturnDepartureCityID, reservation.ReturnDepartureVenueID, reservation.ReturnDepartureTimeID,
			reservation.ReturnDestinationCityID, reservation.ReturnDestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID, reservation.ReturnAirlineID,
			reservation.DriverNotes, reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors,
			reservation.DepartureNumChildren, reservation.ReturnNumAdults, reservation.ReturnNumStudents, reservation.ReturnNumSeniors, reservation.ReturnNumChildren,
			reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID, reservation.CustomDepartureID,
			reservation.DepartureDate, reservation.ReturnDate, 2, reservation.TripID, reservation.ReturnTripID, reservation.BalanceOwing, reservation.ElavonTransactionID, reservation.FlightNumber, reservation.FlightTime)

		if err != nil {
			log.Printf("Error creating return reservation: %s", err)
		}

		insertError = err

		id, _ = result.LastInsertId()
		reservation.ReservationID = int(id)

		returndeparturecity := store.GetCityName(reservation.ReturnDepartureCityID)
		returndeparturevenue := store.GetVenueName(reservation.ReturnDepartureVenueID)
		returndeparturedate := reservation.ReturnDate.Format("2006-01-02")
		returndeparturetime := strconv.Itoa(store.GetDepartureTime(reservation.ReturnDepartureTimeID))
		returndestinationcity := store.GetCityName(reservation.ReturnDestinationCityID)
		returndestinationvenue := store.GetCityName(reservation.ReturnDestinationVenueID)
		returnnumadults := strconv.Itoa(reservation.ReturnNumAdults)
		returnnumseniors := strconv.Itoa(reservation.ReturnNumSeniors)
		returnnumstudents := strconv.Itoa(reservation.ReturnNumStudents)
		returnnumchildren := strconv.Itoa(reservation.ReturnNumChildren)

		returndetails = FormatTripDetails(returndeparturecity, returndeparturevenue, returndeparturedate,
			returndeparturetime, returndestinationcity, returndestinationvenue, returnnumadults, returnnumseniors,
			returnnumstudents, returnnumchildren)

	} else {
		result, err := store.db.Exec("INSERT INTO reservations("+
			"clientid, departurecityid, departurevenueid, departuretimeid, "+
			"destinationcityid, destinationvenueid, discountcodeid, departureairlineid, drivernotes, "+
			"internalnotes, departurenumadults, departurenumstudents, departurenumseniors, departurenumchildren, "+
			"price, status, hash, customdepartureid, customdestinationid, "+
			"departuredate, triptypeid, tripid, balanceowing, elavontransactionid, flightnumber, flighttime) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?)",
			reservation.ClientID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID, reservation.DriverNotes,
			reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors, reservation.DepartureNumChildren,
			reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID, reservation.CustomDepartureID,
			reservation.DepartureDate, 1, reservation.TripID, reservation.BalanceOwing, reservation.ElavonTransactionID, reservation.FlightNumber, reservation.FlightTime)

		if err != nil {
			log.Printf("Error creating one-way reservation: %s", err)
		}

		insertError = err

		id, _ = result.LastInsertId()
		reservation.ReservationID = int(id)
	}

	log.Printf("New Reservation ID: %d", reservation.ReservationID)

	departurecity := store.GetCityName(reservation.DepartureCityID)
	departurevenue := store.GetVenueName(reservation.DepartureVenueID)
	departuredate := reservation.DepartureDate.Format("2006-01-02")
	departuretime := strconv.Itoa(store.GetDepartureTime(reservation.DepartureTimeID))
	destinationcity := store.GetCityName(reservation.DestinationCityID)
	destinationvenue := store.GetVenueName(reservation.DestinationVenueID)
	numadults := strconv.Itoa(reservation.DepartureNumAdults)
	numseniors := strconv.Itoa(reservation.DepartureNumSeniors)
	numstudents := strconv.Itoa(reservation.DepartureNumStudents)
	numchildren := strconv.Itoa(reservation.DepartureNumChildren)

	log.Printf("Total Reservation Cost: $%f", reservation.Price)

	departuredetails = FormatTripDetails(departurecity, departurevenue, departuredate, departuretime,
		destinationcity, destinationvenue, numadults, numseniors, numstudents, numchildren)

	//var client Client
	//client.ClientID = reservation.ClientID
	//client.ClientID = reservation.ClientDetails.ClientID

	//store.GetClientInfo(&client)

	//SendConfirmationEmail(client.Email, departuredetails, returndetails, reservation.ReservationID, reservation.Price)
	SendConfirmationEmail(reservation.ClientDetails.Email, departuredetails, returndetails, reservation.ReservationID, reservation.Price)

	return insertError
}

//SignInUser - authenticate user
func (store *dbStore) SignInUser(client *Client) error {

	//create a client object
	storedClient := &Client{}

	err := store.db.QueryRow("select password, roleid, clientid from accountdetails where username=?",
		client.Username).Scan(&storedClient.Password, &client.RoleID, &client.ClientID)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Unauthorized")
			//***handle this better, display credential error***
		} else {
			log.Printf("Sign In Error: %s", err.Error())
		}

	} else {
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(storedClient.Password), []byte(client.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			log.Print("Incorrect Password")
			err = errors.New("Incorrect Password")
		} else {
			log.Print("Success!")
			err = nil
		}
	}

	return err
}

//GetVenues - return all venues in database
func (store *dbStore) GetVenues() []Venues {

	//Get venue count
	venueCount := store.GetVenueCount()

	row, err := store.db.Query("SELECT venueid, c.cityid, c.name, v.name, v.extracost, v.active, v.extratime FROM venues v " +
		"INNER JOIN cities c ON v.cityid = c.cityid")

	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		log.Printf("Error retrieving venues: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store venues
	var venueSlice = make([]Venues, venueCount)
	var indx int

	//loop through returned venues
	indx = 0
	for row.Next() {
		err = row.Scan(
			&venueSlice[indx].VenueID, &venueSlice[indx].CityID,
			&venueSlice[indx].CityName, &venueSlice[indx].VenueName,
			&venueSlice[indx].ExtraCost, &venueSlice[indx].Active,
			&venueSlice[indx].ExtraTime,
		)
		if err != nil {
			//if error print to console
			if err == sql.ErrNoRows {
				log.Print("No venue found")
			} else {
				log.Printf("Error retrieving venues: %s", err.Error())
			}
		}

		indx++
	}

	return venueSlice
}

//GetVenueName - return venue name given id
func (store *dbStore) GetVenueName(venueid int) string {

	row := store.db.QueryRow("SELECT name from venues where venueid = ?",
		venueid)

	var venuename string

	err := row.Scan(&venuename)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No vanue matches the id")
		} else {
			log.Printf("Error retrieving venue name from id: %s", err.Error())
		}
	}

	return venuename
}

//GetVenueCount - return number of venues
func (store *dbStore) GetVenueCount() int {
	var venueCount int

	row, err := store.db.Query("select count(venueid) from venues")

	if err != nil {
		log.Printf("Error querying venues: %s", err.Error())
	}

	row.Next()
	err = row.Scan(
		&venueCount,
	)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Print("No venues returned")
		} else {
			log.Printf("Error retrieving venue count: %s", err.Error())
		}
	}

	return venueCount
}

//AddVenue - add venue
func (store *dbStore) AddVenue(venue Venues) error {

	_, err := store.db.Exec("INSERT INTO venues(cityid, name, extracost, active, extratime) "+
		"VALUES (?,?,?,?,?)", venue.CityID, venue.VenueName, venue.ExtraCost, venue.Active, venue.ExtraTime)

	if err != nil {
		log.Printf("Error inserting venue: %s", err.Error())
	}

	return err
}

//UpdateVenue - update venue details
func (store *dbStore) UpdateVenue(venue *Venues) error {
	log.Printf("update %s in database", venue.VenueName)

	_, updateerr := store.db.Exec("UPDATE venues SET name = ?, extracost = ?, active = ?, extratime = ? "+
		" WHERE venueid = ?", venue.VenueName, venue.ExtraCost, venue.Active, venue.ExtraTime, venue.VenueID)

	if updateerr != nil {
		log.Printf("Error updating venue: %s", updateerr.Error())
	} else {
		log.Printf("Update Venue: %d", venue.VenueID)
	}

	return updateerr
}

//DeleteVenue - delete venue record
func (store *dbStore) DeleteVenue(venueid int) error {
	_, updateerr := store.db.Exec("DELETE FROM venues WHERE venueid = ?", venueid)

	if updateerr != nil {
		log.Printf("Error deleting venue: %s", updateerr.Error())
	} else {
		log.Printf("Delete Venue: %d", venueid)
	}

	return updateerr
}

//GetCities - return all cities in database
func (store *dbStore) GetCities() []Cities {

	//get count of all cities
	cityCount := store.GetCityCount()

	row, err := store.db.Query("select c.cityid, name, northoffset, southoffset from cities c " +
		"inner join cityoffsets cos on c.cityid = cos.cityid")

	if err != nil {
		log.Printf("Error querying cities: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all cities
	var citySlice = make([]Cities, cityCount)
	var indx int

	//loop through returned cities
	indx = 0
	for row.Next() {
		err = row.Scan(
			&citySlice[indx].CityID, &citySlice[indx].CityName, &citySlice[indx].NorthOffset,
			&citySlice[indx].SouthOffset,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No cities found")
			} else {
				log.Printf("Error retrieving cities: %s", err.Error())
			}
		}

		indx++
	}

	return citySlice
}

//GetCityID - return city id given the name
func (store *dbStore) GetCityID(cityname string) int {

	row := store.db.QueryRow("SELECT cityid from cities where name = ?",
		cityname)

	var cityid int

	err := row.Scan(&cityid)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No city matches the name")
		} else {
			log.Printf("Error retrieving cityid from name: %s", err.Error())
		}
	}

	return cityid
}

//GetCityName - return city name given the id
func (store *dbStore) GetCityName(cityid int) string {

	row := store.db.QueryRow("SELECT name from cities where cityid = ?",
		cityid)

	var name string

	err := row.Scan(&name)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No city matches the id")
		} else {
			log.Printf("Error retrieving city name from id: %s", err.Error())
		}
	}

	return name
}

//GetCityCount  - return number of cities
func (store *dbStore) GetCityCount() int {
	var cityCount int

	row, err := store.db.Query("select count(cityid) from cities")

	row.Next()
	err = row.Scan(
		&cityCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No cities returned")
		} else {
			log.Printf("Error retrieving city count: %s", err.Error())
		}
	}

	return cityCount
}

//AddCity - add city
func (store *dbStore) AddCity(city Cities) error {

	result, err := store.db.Exec("INSERT INTO cities(name) "+
		"VALUES (?)", city.CityName)

	if err != nil {
		log.Printf("Error inserting city: %s", err.Error())
	}

	//get id from client insertion transaction
	id, _ := result.LastInsertId()

	//create account details record linked to client record
	result, err = store.db.Exec("INSERT INTO cityoffsets(cityid, northoffset, southoffset) "+
		"VALUES (?, ?, ?)", id, city.NorthOffset, city.SouthOffset)

	return err
}

//UpdateCity - update city details
func (store *dbStore) UpdateCity(city *Cities) error {
	log.Printf("update %s in database", city.CityName)

	_, updateerr := store.db.Exec("UPDATE cities SET name = ? WHERE cityid = ?", city.CityName, city.CityID)

	if updateerr != nil {
		log.Printf("Error updating city: %s", updateerr.Error())
	} else {
		log.Printf("Update City: %d", city.CityID)
	}

	_, updateerr = store.db.Exec("UPDATE cityoffsets SET northoffset = ?, southoffset = ? "+
		"WHERE cityid = ?", city.NorthOffset, city.SouthOffset, city.CityID)

	if updateerr != nil {
		log.Printf("Error updating city offset: %s", updateerr.Error())
	} else {
		log.Printf("Update City Offset: %d", city.CityID)
	}

	return updateerr
}

//DeleteCity - delete city record
func (store *dbStore) DeleteCity(cityid int) error {
	_, updateerr := store.db.Exec("DELETE FROM cities WHERE cityid = ?", cityid)

	if updateerr != nil {
		log.Printf("Error deleting city: %s", updateerr.Error())
	} else {
		log.Printf("Delete City: %d", cityid)
	}

	return updateerr
}

//GetDepartureTimes - return all departuretimes in database
func (store *dbStore) GetDepartureTimes() []DepartureTimes {

	//get count for all departure times
	departuretimeCount := store.GetDepartureTimesCount()

	row, err := store.db.Query("select departuretimeid, cityid, departuretime, " +
		"recurring, startdate, enddate from departuretimes")

	if err != nil {
		log.Printf("Error retrieving departure times: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var departureTimesSlice = make([]DepartureTimes, departuretimeCount)

	var startdate mysql.NullTime
	var enddate mysql.NullTime
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&departureTimesSlice[indx].DepartureTimeID, &departureTimesSlice[indx].CityID,
			&departureTimesSlice[indx].DepartureTime, &departureTimesSlice[indx].Recurring,
			&startdate, &enddate,
		)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No times found")
			} else {
				log.Printf("Error retrieving times: %s", err.Error())
			}
		} else {
			//store dates in departure time slice if valid dates, otherwise empty date
			if startdate.Valid {
				departureTimesSlice[indx].StartDate = startdate.Time
			} else {
				departureTimesSlice[indx].StartDate = time.Time{}
			}

			if enddate.Valid {
				departureTimesSlice[indx].EndDate = enddate.Time
			} else {
				departureTimesSlice[indx].EndDate = time.Time{}
			}

			departureTimesSlice[indx].Epoch = time.Time{}
		}

		departureTimesSlice[indx].CityList = store.GetCities()

		indx++
	}

	return departureTimesSlice
}

//GetDepartureTime - return departuretime name given the id
func (store *dbStore) GetDepartureTime(departuretimeid int) int {
	row := store.db.QueryRow("SELECT departuretime from departuretimes where departuretimeid = ?",
		departuretimeid)

	var departuretime int

	err := row.Scan(&departuretime)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No departuretime matches the id")
		} else {
			log.Printf("Error retrieving departuretime from id: %s", err.Error())
		}
	}

	return departuretime
}

//GetDepartureTimesCount - return count of all departure times
func (store *dbStore) GetDepartureTimesCount() int {
	var departuretimeCount int

	row, err := store.db.Query("select count(departuretimeid) from departuretimes")

	row.Next()
	err = row.Scan(
		&departuretimeCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No departtimes returned")
		} else {
			log.Printf("Error retrieving departure times: %s", err.Error())
		}
	}

	return departuretimeCount
}

//UpdateDepartureTime - update departure time
func (store *dbStore) UpdateDepartureTime(departuretime *DepartureTimes) error {

	var epoch = time.Time{}
	var err error

	log.Printf("cityid: %d", departuretime.CityID)
	log.Printf("departuretime: %d", departuretime.DepartureTime)
	log.Printf("recurring: %d", departuretime.Recurring)
	log.Printf("startdate: %s", departuretime.StartDate)
	log.Printf("enddate: %s", departuretime.EndDate)

	if departuretime.StartDate.After(epoch) && departuretime.EndDate.After(epoch) {
		log.Printf("start and end date present")
		_, err = store.db.Exec("UPDATE departuretimes SET cityid = ?, departuretime = ?,"+
			" recurring = ?, startdate = ?, enddate = ? WHERE departuretimeid = ?",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.StartDate, departuretime.EndDate, departuretime.DepartureTimeID)
	} else if departuretime.StartDate.After(epoch) {
		log.Printf("start date present")
		_, err = store.db.Exec("UPDATE departuretimes SET cityid = ?, departuretime = ?,"+
			" recurring = ?, startdate = ? WHERE departuretimeid = ?",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.StartDate, departuretime.DepartureTimeID)
	} else if departuretime.EndDate.After(epoch) {
		log.Printf("end date present")
		_, err = store.db.Exec("UPDATE departuretimes SET cityid = ?, departuretime = ?,"+
			" recurring = ?, enddate = ? WHERE departuretimeid = ?",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.EndDate, departuretime.DepartureTimeID)
	} else {
		log.Printf("neither start or end date present")
		_, err = store.db.Exec("UPDATE departuretimes SET cityid = ?, departuretime = ?,"+
			" recurring = ? WHERE departuretimeid = ?",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.DepartureTimeID)
	}

	if err != nil {
		log.Printf("Error updating departure time: %s", err.Error())
	} else {
		log.Printf("Update Departure Time: %d", departuretime.DepartureTimeID)
	}

	return err
}

//AddDepartureTime - add departure time
func (store *dbStore) AddDepartureTime(departuretime *DepartureTimes) error {

	var epoch = time.Time{}
	var err error
	var result sql.Result

	log.Printf("cityid: %d", departuretime.CityID)
	log.Printf("departuretime: %d", departuretime.DepartureTime)
	log.Printf("recurring: %d", departuretime.Recurring)
	log.Printf("startdate: %s", departuretime.StartDate)
	log.Printf("enddate: %s", departuretime.EndDate)

	if departuretime.StartDate.After(epoch) && departuretime.EndDate.After(epoch) {
		log.Printf("start and end date present")
		result, err = store.db.Exec("INSERT INTO departuretimes("+
			"cityid, departuretime, recurring, startdate, enddate) VALUES (?, ?, ?, ?, ?)",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.StartDate, departuretime.EndDate)
	} else if departuretime.StartDate.After(epoch) {
		log.Printf("start date present")
		result, err = store.db.Exec("INSERT INTO departuretimes("+
			"cityid, departuretime, recurring, enddate) VALUES (?, ?, ?, ?)",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.EndDate)
	} else if departuretime.EndDate.After(epoch) {
		log.Printf("end date present")
		result, err = store.db.Exec("INSERT INTO departuretimes("+
			"cityid, departuretime, recurring, startdate) VALUES (?, ?, ?, ?)",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring,
			departuretime.StartDate)
	} else {
		log.Printf("neither start or end date present")
		result, err = store.db.Exec("INSERT INTO departuretimes("+
			"cityid, departuretime, recurring) VALUES (?, ?, ?)",
			departuretime.CityID, departuretime.DepartureTime, departuretime.Recurring)
	}

	id, _ := result.LastInsertId()
	departuretime.DepartureTimeID = int(id)

	if err != nil {
		log.Printf("Error updating departure time: %s", err.Error())
	} else {
		log.Printf("Update Departure Time: %d", departuretime.DepartureTimeID)
	}

	return err
}

//GetClientInfo - from client username return all client info
func (store *dbStore) GetClientInfo(client *Client) error {

	var row *sql.Rows
	var err error

	if client.ClientID > 0 {
		row, err = store.db.Query(
			"select c.clientid, a.roleid, firstname, lastname, phone, email, streetaddress, "+
				"city, province, postalcode, country from clients c inner join "+
				"accountdetails a on c.clientid = a.clientid "+
				"where c.clientid=?", client.ClientID)
	} else {
		row, err = store.db.Query(
			"select c.clientid, a.roleid, firstname, lastname, phone, email, streetaddress, "+
				"city, province, postalcode, country from clients c inner join "+
				"accountdetails a on c.clientid = a.clientid "+
				"where a.username=?", client.Username)
	}

	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		log.Printf("Error retrieving client 1: %s", err.Error())
		return err
	}
	defer row.Close()

	//store client into into local variables
	row.Next()
	err = row.Scan(
		&client.ClientID, &client.RoleID, &client.Firstname, &client.Lastname, &client.Phone, &client.Email,
		&client.StreetAddress, &client.City, &client.Province, &client.PostalCode, &client.Country,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No client found")
		} else {
			log.Printf("Error retrieving client 2: %s", err.Error())
		}
	}

	return err
}

//GetDriverIDFromClient - given the clientid return driverid
func (store *dbStore) GetDriverIDFromClient(clientid int) (int, error) {

	var row *sql.Rows
	var err error

	if clientid > 0 {
		row, err = store.db.Query(
			//get accountdetailid from via clientid
			"SELECT d.driverid FROM accountdetails a inner join drivers d on "+
				"a.accountdetailid = d.accountdetailid WHERE clientid = ?", clientid)
	} else {
		return 0, nil
	}

	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		log.Printf("Error retrieving client: %s", err.Error())
		return 0, err
	}
	defer row.Close()

	var driverid int

	//store client into into local variables
	row.Next()
	err = row.Scan(
		&driverid,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No driver found")
		} else {
			log.Printf("Error retrieving driver: %s", err.Error())
		}
	}

	return driverid, err
}

//GetOrAddTrip - store appropriate tripid in reservation, if no trip
//				 generate new trip
func (store *dbStore) GetOrAddTrip(reservation *Reservation) error {
	var tripid int
	var returntripid int
	var tripid64 int64
	var returntripid64 int64
	var numtrippassengers int
	var numreturntrippassengers int
	var numpassengers int
	var numreturnpassengers int
	var capacity int
	var returncapacity int
	var omitted int
	var returnomitted int

	numpassengers = reservation.DepartureNumAdults + reservation.DepartureNumChildren + reservation.DepartureNumSeniors + reservation.DepartureNumStudents
	numreturnpassengers = reservation.ReturnNumAdults + reservation.ReturnNumChildren + reservation.ReturnNumSeniors + reservation.ReturnNumStudents

	//err := store.db.QueryRow("select tripid, numpassengers, capacity, omitted from trips where tripid=?",
	//	reservation.TripID).Scan(&tripid, &numtrippassengers, &capacity, &omitted)

	err := store.db.QueryRow("select tripid, numpassengers, capacity, omitted from trips where "+
		"departuredate=? and departuretimeid=?",
		reservation.DepartureDate, reservation.DepartureTimeID).Scan(&tripid, &numtrippassengers,
		&capacity, &omitted)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No trip, create one")

			//check if trip in omitted table

			result, inserterr := store.db.Exec("INSERT INTO trips(departuredate, departuretimeid, numpassengers) "+
				"VALUES (?,?,?)", reservation.DepartureDate, reservation.DepartureTimeID, numpassengers)

			if inserterr != nil {
				log.Printf("Error inserting new trip: %s", inserterr.Error())
			} else {
				tripid64, _ = result.LastInsertId()
				log.Printf("TripID: %d", tripid64)
				tripid = int(tripid64)
				err = nil
			}
		} else {
			log.Printf("Error retrieving trip: %s", err.Error())
		}
	} else {
		//add passengers from current reservation to trip passengers
		vacancies := capacity - numpassengers

		if omitted == 1 {
			errorstring := fmt.Sprintf("This trip is has been cancelled")
			err = errors.New(errorstring)
		} else {
			if vacancies < numpassengers {
				errorstring := fmt.Sprintf("This trip is over capacity, only %d vacancies", vacancies)
				err = errors.New(errorstring)
			} else {
				if float32(capacity)*(0.75) >= float32(vacancies) {

					trip := store.GetTrip(tripid)

					SendCapacityEmail(trip)
				}

				numpassengers += numtrippassengers

				_, err = store.db.Exec("UPDATE trips SET numpassengers "+
					"= ? WHERE tripid = ?", numpassengers, tripid)

				if err != nil {
					log.Printf("Error updating trip passengers: %s", err.Error())
				}

				log.Print("updated passengers num")

			}
		}

		//***if over 75% need notification***
	}

	reservation.TripID = int(tripid)

	if numreturnpassengers > 0 {

		err = store.db.QueryRow("select tripid, numpassengers, capacity, omitted from trips where "+
			"returndate=? and returndeparturetimeid=?",
			reservation.ReturnDate, reservation.ReturnDepartureTimeID).Scan(&returntripid,
			&numreturntrippassengers, &returncapacity, &returnomitted)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No trip, create one")

				//check if trip in omitted table
				result, inserterr := store.db.Exec("INSERT INTO trips(departuredate, departuretimeid, numpassengers) "+
					"VALUES (?,?,?)", reservation.ReturnDate, reservation.ReturnDepartureTimeID, numreturnpassengers)

				if inserterr != nil {
					log.Printf("Error inserting new return trip: %s", inserterr.Error())
				} else {
					returntripid64, _ = result.LastInsertId()
					log.Printf("ReturnTripID: %d", returntripid64)
					returntripid = int(returntripid64)
					err = nil
				}
			} else {
				log.Printf("Error retrieving return trip: %s", err.Error())
			}
		} else {
			//add passengers from current reservation to trip passengers
			vacancies := returncapacity - numreturnpassengers

			if returnomitted == 1 {
				errorstring := fmt.Sprintf("This return trip is has been cancelled")
				err = errors.New(errorstring)
			} else {
				if vacancies < numreturnpassengers {
					errorstring := fmt.Sprintf("This return trip is over capacity, only %d vacancies", vacancies)
					err = errors.New(errorstring)
				} else {
					if float32(returncapacity)*(0.75) >= float32(vacancies) {

						trip := store.GetTrip(returntripid)

						SendCapacityEmail(trip)
					}

					numreturnpassengers += numreturntrippassengers

					_, err = store.db.Exec("UPDATE trips SET numpassengers "+
						"= ? WHERE tripid = ?", numreturnpassengers, returntripid)

					if err != nil {
						log.Printf("Error updating trip passengers: %s", err.Error())
					}

					log.Print("updated passengers num for return trip")

				}
			}

			//***if over 75% need notification***
		}
	}

	return err
}

//GetTrips - return all trips
func (store *dbStore) GetTrips() []Trips {
	var tripCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(tripid) from trips")

	row.Next()
	err = row.Scan(
		&tripCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No trips returned")
		} else {
			log.Printf("Error retrieving trip count: %s", err.Error())
		}
	}

	//problem with more than 20, not sure why
	if tripCount > 20 {
		tripCount = 20
	}

	row, err = store.db.Query("select tripid, departuredate, t.departuretimeid, " +
		"numpassengers, driverid, vehicleid, capacity, departuretime, " +
		"omitted from trips t inner join " +
		"departuretimes dt on t.departuretimeid = dt.departuretimeid " +
		//"where year(departuredate) = year(CURDATE()) AND month(departuredate) = month(CURDATE()) " +
		"order by departuredate desc " +
		"limit 20")

	if err != nil {
		log.Printf("Error retrieving trips: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var tripSlice = make([]Trips, tripCount)

	var departuredate mysql.NullTime
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&tripSlice[indx].TripID, &departuredate,
			&tripSlice[indx].DepartureTimeID, &tripSlice[indx].NumPassengers,
			&tripSlice[indx].DriverID, &tripSlice[indx].VehicleID,
			&tripSlice[indx].Capacity, &tripSlice[indx].DepartureTime,
			&tripSlice[indx].Omitted,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No trips found")
			} else {
				log.Printf("Error retrieving trips: %s", err.Error())
			}
		} else {
			//store dates in departure time slice if valid dates, otherwise empty date
			if departuredate.Valid {
				tripSlice[indx].DepartureDate = departuredate.Time
			} else {
				tripSlice[indx].DepartureDate = time.Time{}
			}

			//populate drivers
			tripSlice[indx].DriverList = make([]Drivers, store.GetDriverCount())
			tripSlice[indx].DriverList = store.GetDrivers(0)

			//populate vehicle
			tripSlice[indx].VehicleList = make([]Vehicles, store.GetVehicleCount())
			tripSlice[indx].VehicleList = store.GetVehicles()
		}

		indx++
	}

	return tripSlice
}

//GetTrip - return one trip based on id
func (store *dbStore) GetTrip(tripid int) Trips {

	row := store.db.QueryRow("select tripid, departuredate, t.departuretimeid, "+
		"numpassengers, driverid, vehicleid, capacity, "+
		"omitted from trips t inner join "+
		"departuretimes dt on t.departuretimeid = dt.departuretimeid "+
		"where tripid = ?", tripid)

	trip := Trips{}
	var departuredate mysql.NullTime

	err := row.Scan(
		&trip.TripID, &departuredate,
		&trip.DepartureTimeID, &trip.NumPassengers,
		&trip.DriverID, &trip.VehicleID,
		&trip.Capacity, &trip.Omitted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No trips found")
		} else {
			log.Printf("Error retrieving trips: %s", err.Error())
		}
	} else {
		//store dates in departure time slice if valid dates, otherwise empty date
		if departuredate.Valid {
			trip.DepartureDate = departuredate.Time
		} else {
			trip.DepartureDate = time.Time{}
		}

		dtrow := store.db.QueryRow("select departuretime "+
			"from departuretimes where departuretimeid = ? ",
			trip.DepartureTimeID)

		dterr := dtrow.Scan(
			&trip.DepartureTime,
		)

		if dterr != nil {
			if dterr == sql.ErrNoRows {
				log.Print("No departuretime found")
			} else {
				log.Printf("Error retrieving departuretime: %s", dterr.Error())
			}
		}
	}

	return trip
}

//UpdateTrip - update driver and vehicle associated with trip
func (store *dbStore) UpdateTrip(trip *Trips) error {

	var capacity int

	if trip.VehicleID > 0 {
		row := store.db.QueryRow("SELECT numseats from vehicles where vehicleid = ?",
			trip.VehicleID)

		err := row.Scan(&capacity)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No capacity matches the vehicle id")
			} else {
				log.Printf("Error retrieving capacity from vehicle id: %s", err.Error())
			}
		}
	} else {
		capacity = 11
	}

	log.Printf("DriverID: %d, VehicleID: %d, Capacity: %d, TripID: %d",
		trip.DriverID, trip.VehicleID, capacity, trip.TripID)

	_, updateerr := store.db.Exec("UPDATE trips SET driverid = ?, "+
		" vehicleid = ?, capacity = ? WHERE tripid = ?", trip.DriverID, trip.VehicleID, capacity, trip.TripID)

	if updateerr != nil {
		log.Printf("Error updating trip: %s", updateerr.Error())
	} else {
		log.Printf("Update Trip: %d", trip.TripID)
	}

	return updateerr
}

//OmitTrip - cancel trip
func (store *dbStore) OmitTrip(trip *Trips) error {

	var tripid int

	row := store.db.QueryRow("SELECT tripid FROM trips WHERE departuredate = ? AND departuretimeid = ?",
		trip.DepartureDate, trip.DepartureTimeID)

	err := row.Scan(
		&tripid,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			_, inserterr := store.db.Exec("INSERT INTO trips (departuredate, departuretimeid, omitted) "+
				"VALUE (?,?,?)", trip.DepartureDate, trip.DepartureTimeID, 1)

			if inserterr != nil {
				log.Printf("Error creating omitted trip: %s", inserterr.Error())
			} else {
				log.Printf("Created Omitted Trip: %d", trip.TripID)
			}

			err = inserterr
		} else {
			log.Printf("Error getting tripid: %s", err.Error())
		}
	} else {
		_, updateerr := store.db.Exec("UPDATE trips SET omitted = ? WHERE tripid = ?",
			1, tripid)

		if updateerr != nil {
			log.Printf("Error omitting trip: %s", updateerr.Error())
		} else {
			log.Printf("Omitted Trip: %d", trip.TripID)
		}

		err = updateerr
	}

	return err
}

//SearchReservations - return all trips (must add parameter to return by date)
func (store *dbStore) SearchReservations(name string, phone int, email string) []SearchReservations {
	var reservationCount int

	var addWhere bool
	var whereClause string

	whereClause = " where "
	addWhere = false

	if len(name) > 0 {
		log.Printf("add name to where")
		whereClause += " c.Firstname like ('" + name + "') or c.LastName like ('" + name + "') "
		addWhere = true
	} else if phone > 0 {
		log.Printf("add phone to where")
		whereClause += " c.Phone = " + strconv.Itoa(phone)
		addWhere = true
	} else if len(email) > 0 {
		log.Printf("add email to where")
		whereClause += " c.Email = '" + email + "' "
		addWhere = true
	}

	var sqlString string

	if addWhere {
		sqlString = "select count(reservationid) " +
			" from reservations r inner join clients c on r.clientid = c.clientid " + whereClause
		log.Printf("sql string created with where, %s", sqlString)
	} else {
		sqlString = "select count(reservationid) from reservations"
		log.Printf("sql string created")
	}

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query(sqlString)

	row.Next()
	err = row.Scan(
		&reservationCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No search reservations returned")
		} else {
			log.Printf("Error retrieving search reservation count: %s", err.Error())
		}
	}

	log.Printf("got the count: %d", reservationCount)

	if reservationCount > 100 {
		reservationCount = 100
	}

	var numadults int
	var numstudents int
	var numseniors int
	var numchildren int

	if addWhere {
		row, err = store.db.Query("select r.reservationid, concat(c.firstname, ' ', c.lastname) as clientname, c.phone, c.email, " +
			"depv.name as departurevenue, desv.name as destinationvenue, r.triptypeid, " +
			"r.departurenumadults, r.departurenumstudents, r.departurenumseniors, r.departurenumchildren, " +
			"r.departuredate, dt.departuretime, r.postponed, r.cancelled " +
			"from reservations r inner join clients c on r.clientid = c.clientid " +
			"inner join venues depv on r.departurevenueid = depv.venueid " +
			"inner join venues desv on r.destinationvenueid = desv.venueid " +
			"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid " + whereClause +
			"order by r.departuredate " +
			"limit 100")
	} else {
		row, err = store.db.Query("select r.reservationid, concat(c.firstname, ' ', c.lastname) as clientname, c.phone, c.email, " +
			"depv.name as departurevenue, desv.name as destinationvenue, r.triptypeid, " +
			"r.departurenumadults, r.departurenumstudents, r.departurenumseniors, r.departurenumchildren, " +
			"r.departuredate, dt.departuretime, r.postponed, r.cancelled " +
			"from reservations r inner join clients c on r.clientid = c.clientid " +
			"inner join venues depv on r.departurevenueid = depv.venueid " +
			"inner join venues desv on r.destinationvenueid = desv.venueid " +
			"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid " +
			"order by r.departuredate " +
			"limit 100")
	}

	if err != nil {
		log.Printf("Error retrieving search reservations: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var searchReservationSlice = make([]SearchReservations, reservationCount)

	log.Printf("retrieved records")

	var departuredate mysql.NullTime
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&searchReservationSlice[indx].ReservationID, &searchReservationSlice[indx].ClientName,
			&searchReservationSlice[indx].Phone, &searchReservationSlice[indx].Email,
			&searchReservationSlice[indx].DepartureVenue, &searchReservationSlice[indx].DestinationVenue,
			&searchReservationSlice[indx].Return, &numadults, &numstudents, &numseniors, &numchildren,
			&departuredate, &searchReservationSlice[indx].DepartureTime,
			&searchReservationSlice[indx].Postponed, &searchReservationSlice[indx].Cancelled,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No search reservations found")
			} else {
				log.Printf("Error retrieving search reservations: %s", err.Error())
			}
		} else {
			//store dates in departure time slice if valid dates, otherwise empty date
			if departuredate.Valid {
				searchReservationSlice[indx].DepartureDate = departuredate.Time
			} else {
				searchReservationSlice[indx].DepartureDate = time.Time{}
			}

			searchReservationSlice[indx].NumPassengers = numadults + numstudents + numseniors + numchildren
		}

		indx++
	}

	return searchReservationSlice
}

//TripReservations - return trips for given date
//reportType 0 = trips, 1 = calendar
func (store *dbStore) SearchTrips(tripdate time.Time, reportType int) []Trips {
	var tripCount int

	var whereClause string

	//populate drivers
	driverlist := make([]Drivers, store.GetDriverCount())
	driverlist = store.GetDrivers(0)
	//populate vehicle
	vehiclelist := make([]Vehicles, store.GetVehicleCount())
	vehiclelist = store.GetVehicles()

	whereClause = " where "

	var epoch = time.Time{}
	var today = time.Now()

	if tripdate.After(epoch) {
		whereClause += " departuredate = '" + tripdate.Format("2006-01-02") + "' "
	} else {
		whereClause += " departuredate = '" + today.Format("2006-01-02") + "' "
	}

	var sqlString string

	sqlString = "select count(departuretimeid) " +
		"from (select departuretimeid " +
		"from trips " + whereClause + " " +
		"group by departuretimeid) as deptimeids"

	log.Printf("%s", sqlString)

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query(sqlString)

	row.Next()
	err = row.Scan(
		&tripCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No trips returned")
		} else {
			log.Printf("Error retrieving trip count: %s", err.Error())
		}
	}
	row.Close()

	log.Printf("got the count: %d", tripCount)

	if tripCount > 50 {
		tripCount = 50
	}

	sqlstring := "select t.tripid, departuredate, t.departuretimeid, departuretime, " +
		" sum(numpassengers) as numpassengers, t.driverid, t.vehicleid, capacity, omitted, " +
		"if(t.driverid > 0, concat(d.firstname, ' ', d.lastname), 'no driver') as drivername, " +
		"if(t.vehicleid > 0, v.licenseplate, 'no vehicle') as vehicle " +
		"from trips t inner join departuretimes dt on t.departuretimeid = dt.departuretimeid " +
		"left join drivers d on d.driverid = t.driverid " +
		"left join vehicles v on v.vehicleid = t.vehicleid " +
		whereClause +
		"group by t.tripid, departuredate, t.departuretimeid, departuretime, " +
		"t.driverid, t.vehicleid, capacity, omitted, " +
		"drivername, vehicle " +
		"order by departuredate desc limit 50"

	row, err = store.db.Query(sqlstring)

	log.Printf(sqlstring)

	if err != nil {
		log.Printf("Error retrieving search trips: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var searchTripSlice = make([]Trips, tripCount)

	log.Printf("SearchTrips: retrieved records")

	var departuredate mysql.NullTime
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&searchTripSlice[indx].TripID, &departuredate, &searchTripSlice[indx].DepartureTimeID,
			&searchTripSlice[indx].DepartureTime,
			&searchTripSlice[indx].NumPassengers, &searchTripSlice[indx].DriverID,
			&searchTripSlice[indx].VehicleID, &searchTripSlice[indx].Capacity,
			&searchTripSlice[indx].Omitted, &searchTripSlice[indx].DriverName,
			&searchTripSlice[indx].LicensePlate,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No search trips found")
			} else {
				log.Printf("Error retrieving search trips: %s", err.Error())
			}
		} else {
			//store dates in departure time slice if valid dates, otherwise empty date
			if departuredate.Valid {
				searchTripSlice[indx].DepartureDate = departuredate.Time
			} else {
				searchTripSlice[indx].DepartureDate = time.Time{}
			}
		}

		searchTripSlice[indx].DriverList = driverlist
		searchTripSlice[indx].VehicleList = vehiclelist

		indx++
	}

	row.Close()

	return searchTripSlice
}

//PostponeReservation - mark reservation as postponed
func (store *dbStore) PostponeReservation(searchreservation *SearchReservations) error {

	log.Printf("Postponing Reservation: %d", searchreservation.ReservationID)

	_, updateerr := store.db.Exec("UPDATE reservations SET postponed = 1 "+
		" WHERE reservationid = ?", searchreservation.ReservationID)

	if updateerr != nil {
		log.Printf("Error postponing reservation: %s", updateerr.Error())
	} else {
		log.Printf("Postpone trip: %d", searchreservation.ReservationID)
	}

	return updateerr
}

//CancelReservation - mark reservation as cancelled
func (store *dbStore) CancelReservation(searchreservation *SearchReservations) error {

	log.Printf("Cancelling Reservation: %d", searchreservation.ReservationID)

	_, updateerr := store.db.Exec("UPDATE reservations SET cancelled = 1 "+
		" WHERE reservationid = ?", searchreservation.ReservationID)

	if updateerr != nil {
		log.Printf("Error postponing reservation: %s", updateerr.Error())
	} else {
		log.Printf("Postpone trip: %d", searchreservation.ReservationID)
	}

	return updateerr
}

//GetVehicleCount - return count of all vehicles
func (store *dbStore) GetVehicleCount() int {
	var vehicleCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(vehicleid) from vehicles")

	row.Next()
	err = row.Scan(
		&vehicleCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No vehicles returned")
		} else {
			log.Printf("Error retrieving vehicle count: %s", err.Error())
		}

		return 0
	}

	row.Close()

	return vehicleCount
}

//GetVehicles - return all vehicles
func (store *dbStore) GetVehicles() []Vehicles {

	row, err := store.db.Query("select vehicleid, licenseplate, numseats, " +
		"make from vehicles")

	if err != nil {
		log.Printf("Error retrieving vehicles times: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var vehicleSlice = make([]Vehicles, store.GetVehicleCount())

	var indx int
	indx = 0
	for row.Next() {
		err = row.Scan(
			&vehicleSlice[indx].VehicleID, &vehicleSlice[indx].LicensePlate,
			&vehicleSlice[indx].NumSeats, &vehicleSlice[indx].Make,
		)

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No vehicles found")
			} else {
				log.Printf("Error retrieving vehicles: %s", err.Error())
			}
		}

		indx++
	}

	row.Close()

	return vehicleSlice
}

//AddVehicle - add vehicle
func (store *dbStore) AddVehicle(vehicle Vehicles) error {

	_, err := store.db.Exec("INSERT INTO vehicles(licenseplate, numseats, make) "+
		"VALUES (?,?,?)", vehicle.LicensePlate, vehicle.NumSeats, vehicle.Make)

	if err != nil {
		log.Printf("Error inserting vehicle: %s", err.Error())
	}

	return err
}

//UpdateVehicle - update vehicle details
func (store *dbStore) UpdateVehicle(vehicle *Vehicles) error {
	log.Printf("update %s in database", vehicle.LicensePlate)

	_, updateerr := store.db.Exec("UPDATE vehicles SET licenseplate = ?, numseats = ?, make = ?"+
		" WHERE vehicleid = ?", vehicle.LicensePlate, vehicle.NumSeats, vehicle.Make, vehicle.VehicleID)

	if updateerr != nil {
		log.Printf("Error updating vehicle: %s", updateerr.Error())
	} else {
		log.Printf("Update Vehicle: %d", vehicle.VehicleID)
	}

	return updateerr
}

//DeleteVenue - delete venue record
func (store *dbStore) DeleteVehicle(vehicleid int) error {
	_, updateerr := store.db.Exec("DELETE FROM vehicles WHERE vehicleid = ?", vehicleid)

	if updateerr != nil {
		log.Printf("Error deleting vehicle: %s", updateerr.Error())
	} else {
		log.Printf("Delete Vehicle: %d", vehicleid)
	}

	return updateerr
}

//GetDriverCount - return count of all drivers
func (store *dbStore) GetDriverCount() int {
	var driverCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(driverid) from drivers")

	row.Next()
	err = row.Scan(
		&driverCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No driver returned")
		} else {
			log.Printf("Error retrieving driver count: %s", err.Error())
		}

		return 0
	}

	row.Close()

	return driverCount
}

//GetDrivers - return all drivers
func (store *dbStore) GetDrivers(driverid int) []Drivers {

	var row *sql.Rows
	var err error
	var length int

	if driverid == 0 {
		row, err = store.db.Query("select driverid, firstname, lastname,  " +
			"concat(firstname, ' ', lastname) as drivername from drivers")

		if err != nil {
			log.Printf("Error retrieving drivers: %s", err.Error())
			return nil
		}
		defer row.Close()

		//create slice to store all departure times
		length = store.GetDriverCount()
	} else {
		row, err = store.db.Query("select driverid, firstname, lastname,  "+
			"concat(firstname, ' ', lastname) as drivername from drivers where driverid = ?",
			driverid)

		if err != nil {
			log.Printf("Error retrieving drivers: %s", err.Error())
			return nil
		}
		defer row.Close()

		//create slice to store all departure times
		length = 1
	}

	driverSlice := make([]Drivers, length)

	var indx int
	indx = 0
	for row.Next() {
		err = row.Scan(
			&driverSlice[indx].DriverID, &driverSlice[indx].FirstName,
			&driverSlice[indx].LastName, &driverSlice[indx].DriverName,
		)

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No drivers found")
			} else {
				log.Printf("Error retrieving drivers: %s", err.Error())
			}
		}

		indx++
	}

	row.Close()

	return driverSlice
}

//AddDriver - add drivers
func (store *dbStore) AddDriver(driver Drivers) error {

	_, err := store.db.Exec("INSERT INTO drivers(firstname, lastname) "+
		"VALUES (?,?)", driver.FirstName, driver.LastName)

	if err != nil {
		log.Printf("Error inserting vehicle: %s", err.Error())
	} else {
		log.Print("Driver added")
	}

	return err
}

//UpdateDriver - update driver details
func (store *dbStore) UpdateDriver(driver *Drivers) error {
	log.Printf("update %d in database", driver.DriverID)

	_, updateerr := store.db.Exec("UPDATE drivers SET firstname = ?, lastname = ?"+
		" WHERE driverid = ?", driver.FirstName, driver.LastName, driver.DriverID)

	if updateerr != nil {
		log.Printf("Error updating driver: %s", updateerr.Error())
	} else {
		log.Printf("Update Driver: %d", driver.DriverID)
	}

	return updateerr
}

//DeleteDriver - delete driver record
func (store *dbStore) DeleteDriver(driverid int) error {
	_, updateerr := store.db.Exec("DELETE FROM drivers WHERE driverid = ?", driverid)

	if updateerr != nil {
		log.Printf("Error deleting driver: %s", updateerr.Error())
	} else {
		log.Printf("Delete Driver: %d", driverid)
	}

	return updateerr
}

//GetAirlineCount - return count of all airlines
func (store *dbStore) GetAirlineCount() int {
	var airlineCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(airlineid) from airlines")

	row.Next()
	err = row.Scan(
		&airlineCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No airlines returned")
		} else {
			log.Printf("Error retrieving airline count: %s", err.Error())
		}

		return 0
	}

	row.Close()

	return airlineCount
}

//GetAirlines - return all airlines
func (store *dbStore) GetAirlines() []Airlines {

	row, err := store.db.Query("select airlineid, name, terminal " +
		"from airlines")

	if err != nil {
		log.Printf("Error retrieving airlines times: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var airlineSlice = make([]Airlines, store.GetAirlineCount())

	var indx int
	indx = 0
	for row.Next() {
		err = row.Scan(
			&airlineSlice[indx].AirlineID, &airlineSlice[indx].Name,
			&airlineSlice[indx].Terminal,
		)

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No airlines found")
			} else {
				log.Printf("Error retrieving airlines: %s", err.Error())
			}
		}

		indx++
	}

	row.Close()

	return airlineSlice
}

//GetPrice - return price for a trip
func (store *dbStore) GetPrice(departurecityid int, destinationcityid int, retdeparturecityid int, retdestinationcityid int, customertypeid int, triptypeid int, discountcode string) float32 {
	var err error

	row := store.db.QueryRow("SELECT price FROM prices "+
		"WHERE departurecityid = ? and destinationcityid = ? and customertypeid = ?",
		departurecityid, destinationcityid, customertypeid)

	if err != nil {
		log.Printf("Error retrieving prices: %s", err.Error())
		return 0
	}

	var price float32

	err = row.Scan(
		&price,
	)

	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Print("No prices found")
		} else {
			log.Printf("Error retrieving prices: %s", err.Error())
		}
	}

	if (departurecityid != retdeparturecityid || destinationcityid != retdestinationcityid) && triptypeid == 2 {
		row := store.db.QueryRow("SELECT price FROM prices "+
			"WHERE departurecityid = ? and destinationcityid = ? and customertypeid = ?",
			retdeparturecityid, retdestinationcityid, customertypeid)

		if err != nil {
			log.Printf("Error retrieving return prices: %s", err.Error())
			return 0
		}

		var retprice float32

		err = row.Scan(
			&retprice,
		)

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No return prices found")
			} else {
				log.Printf("Error retrieving return prices: %s", err.Error())
			}
		}

		price = (price + retprice) * 0.9

	} else if triptypeid == 2 {
		price = (price * 2) * 0.9
	}

	if discountcode != "" {
		discount := DiscountCode{}

		discount.Name = discountcode

		err = store.GetDiscount(&discount)

		if err != nil {
			log.Printf("Error retrieving discount code: %s", err.Error())
		} else {
			if discount.Percentage > 0 {
				price *= (float32)(100-discount.Percentage) / 100.0
				log.Printf("2 Price: %f", price)
			} else if discount.Amount > 0 {
				price -= (float32)(discount.Amount)
			} else {
				log.Print("Invalid discount code, both percentage and amount are zero")
			}
		}
	}

	return price
}

//GetPrices - return all prices
func (store *dbStore) GetPrices() []Prices {
	var priceCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(priceid) from prices")

	row.Next()
	err = row.Scan(
		&priceCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No prices returned")
		} else {
			log.Printf("Error retrieving price count: %s", err.Error())
		}
	}

	row, err = store.db.Query("select p.priceid, p.departurecityid, dep.name, p.destinationcityid, " +
		"dest.name, p.customertypeid, c.name, p.price " +
		"from prices p inner join cities dep on p.departurecityid = dep.cityid " +
		"inner join cities dest on p.destinationcityid = dest.cityid " +
		"inner join customertypes c on p.customertypeid = c.customertypeid ")

	if err != nil {
		log.Printf("Error retrieving prices: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var priceSlice = make([]Prices, priceCount)

	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&priceSlice[indx].PriceID, &priceSlice[indx].DepartureCityID,
			&priceSlice[indx].DepartureCity, &priceSlice[indx].DestinationCityID,
			&priceSlice[indx].DestinationCity, &priceSlice[indx].CustomerTypeID,
			&priceSlice[indx].CustomerType, &priceSlice[indx].Price,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No prices found")
			} else {
				log.Printf("Error retrieving prices: %s", err.Error())
			}
		}

		indx++
	}

	return priceSlice
}

//UpdatePrices - update price for a trip
func (store *dbStore) UpdatePrice(price *Prices) error {
	log.Printf("update price: %d in database", price.PriceID)

	_, updateerr := store.db.Exec("UPDATE prices SET price = ? "+
		"WHERE priceid = ?", price.Price, price.PriceID)

	if updateerr != nil {
		log.Printf("Error updating price: %s", updateerr.Error())
	} else {
		log.Printf("Update Price: %d", price.PriceID)
	}

	return updateerr
}

//GetDiscount - populate discount if applicable discount code found
func (store *dbStore) GetDiscount(discountcode *DiscountCode) error {
	row := store.db.QueryRow("SELECT discountcodeid, percentage, amount, startdate, enddate from discountcodes where name = ?",
		discountcode.Name)

	var startdate mysql.NullTime
	var enddate mysql.NullTime

	err := row.Scan(&discountcode.DiscountCodeID, &discountcode.Percentage, &discountcode.Amount, &startdate, &enddate)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No discount code matches the name")
		} else {
			log.Printf("Error retrieving discount code from name: %s", err.Error())
		}
	} else {

		today := time.Now()

		//ensure discount code is valid, if so get type and value
		if today.After(startdate.Time) && today.Before(enddate.Time) {
			discountcode.StartDate = startdate.Time
			discountcode.EndDate = enddate.Time

			if discountcode.Percentage > 0 {
				//mark type as percentage
				discountcode.Type = 1
			} else if discountcode.Amount > 0 {
				//mark type as amount
				discountcode.Type = 2
			} else {
				log.Print("Both percentage and amount of discount code is zero, no discount provided")
			}

		} else {
			log.Print("Discount code is not currently valid")
		}
	}

	return err
}

//AddVenueFee - return addition venue charge (only for certain venues)
func (store *dbStore) AddVenueFee(departurevenueid int, destinationvenueid int, retdeparturevenueid int, retdestinationvenueid int) float32 {

	var totalcosts float32
	totalcosts = 0.0

	venueids := [4]int{departurevenueid, destinationvenueid, retdeparturevenueid, retdestinationvenueid}

	for i := 0; i < len(venueids); i++ {
		row, err := store.db.Query("SELECT extracost FROM venues "+
			"WHERE venueid = ?", venueids[i])

		if err != nil {
			log.Printf("Error retrieving extra costs: %s", err.Error())
			return 0
		}
		defer row.Close()

		//slices to store extra costs, can only have a max of 4
		var extracost float32

		for row.Next() {
			err = row.Scan(
				&extracost,
			)

			if err != nil {
				// If an entry with the username does not exist, send an "Unauthorized"(401) status
				if err == sql.ErrNoRows {
					log.Print("No extra costs found")
				} else {
					log.Printf("Error retrieving extra costs: %s", err.Error())
				}
			} else {
				totalcosts += extracost

				log.Printf("extra cost %f added totalcost %f", extracost, totalcosts)
			}
		}
	}

	return totalcosts

}

//DriverReservations - return list of reservations for driver
func (store *dbStore) DriverReservations(driverid int, reportdate time.Time) DriverReportForm {
	var reservationCount int

	var whereClause string
	var sqlString string

	drivers := store.GetDrivers(0)

	whereClause = " where r.postponed = 0 and r.cancelled = 0 and t.driverid = " + strconv.Itoa(driverid)

	var nulldate time.Time

	if reportdate != nulldate {
		whereClause += " and r.departuredate = '" + reportdate.Format("2006-01-02") + "'"
	} else {
		whereClause += " and r.departuredate = '" + time.Time{}.Format("2006-01-02") + "'"
	}

	sqlString = "select count(reservationid) " +
		"from reservations r inner join trips t on r.tripid = t.tripid " +
		"inner join venues depv on r.departurevenueid = depv.venueid " +
		"inner join venues desv on r.destinationvenueid = desv.venueid " +
		"inner join cities depc on r.departurecityid = depc.cityid " +
		"inner join cities des on r.destinationcityid = des.cityid " +
		"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid " +
		"" + whereClause
	log.Printf("sql string created with where, %s", sqlString)

	row, err := store.db.Query(sqlString)

	row.Next()
	err = row.Scan(
		&reservationCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No driver reservations returned")
		} else {
			log.Printf("Error retrieving driver reservation count: %s", err.Error())
		}
	}

	log.Printf("reservation count: %d", reservationCount)
	var driverReportSlice = make([]DriverReport, reservationCount)

	if reservationCount > 0 {
		row, err = store.db.Query("select r.reservationid, concat(c.firstname, ' ', c.lastname) as clientname, " +
			"depv.name as departurevenue, desv.name as destinationvenue, " +
			"depc.name as departurecity, des.name as destinationcity, " +
			"sum(r.departurenumadults + r.departurenumstudents + r.departurenumseniors + r.departurenumchildren) as numpassengers, " +
			"dt.departuretime from reservations r inner join clients c on r.clientid = c.clientid " +
			"inner join trips t on r.tripid = t.tripid " +
			"inner join venues depv on r.departurevenueid = depv.venueid " +
			"inner join venues desv on r.destinationvenueid = desv.venueid " +
			"inner join cities depc on r.departurecityid = depc.cityid " +
			"inner join cities des on r.destinationcityid = des.cityid " +
			"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid " + whereClause +
			"group by r.reservationid, clientname, departurevenue, destinationvenue, " +
			"departurecity, destinationcity, dt.departuretime")

		if err != nil {
			log.Printf("Error retrieving driver reservations: %s", err.Error())
			return DriverReportForm{}
		}
		defer row.Close()

		log.Printf("retrieved records")

		var indx int

		indx = 0
		for row.Next() {
			err = row.Scan(
				&driverReportSlice[indx].ReservationID, &driverReportSlice[indx].ClientName,
				&driverReportSlice[indx].DepartureVenue, &driverReportSlice[indx].DestinationVenue,
				&driverReportSlice[indx].DepartureCity, &driverReportSlice[indx].DestinationCity,
				&driverReportSlice[indx].NumPassengers, &driverReportSlice[indx].DepartureTime,
			)

			if err != nil {
				if err == sql.ErrNoRows {
					log.Print("No driver reservations found")
				} else {
					log.Printf("Error retrieving driver reservations: %s", err.Error())
				}
			}

			indx++
		}
	}

	driverReportForm := DriverReportForm{}

	driverReportForm.Drivers = drivers
	driverReportForm.DriverReports = driverReportSlice

	return driverReportForm
}

//GetTravelAgencyCount - return count of all travel agencies
func (store *dbStore) GetTravelAgencyCount() int {
	var travelAgencyCount int

	//need to add where clause to these queries to use today's date
	row, err := store.db.Query("select count(travelagencyid) from travelagencies")

	row.Next()
	err = row.Scan(
		&travelAgencyCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No travel agencies returned")
		} else {
			log.Printf("Error retrieving travel agency count: %s", err.Error())
		}

		return 0
	}

	return travelAgencyCount
}

//GetTravelAgencies - return all travel agencies
func (store *dbStore) GetTravelAgencies() []TravelAgencies {

	row, err := store.db.Query("select travelagencyid, travelagentname, iatanumber" +
		"from travelagencies")

	if err != nil {
		log.Printf("Error retrieving travel agencies: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var travelAgencySlice = make([]TravelAgencies, store.GetTravelAgencyCount())

	var indx int
	indx = 0
	for row.Next() {
		err = row.Scan(
			&travelAgencySlice[indx].TravelAgencyID, &travelAgencySlice[indx].TravelAgencyName,
			&travelAgencySlice[indx].IATANumber,
		)

		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No travel agencies found")
			} else {
				log.Printf("Error retrieving travel agencies: %s", err.Error())
			}
		}

		indx++
	}

	return travelAgencySlice
}

//TravelAgentRecords - return list of reservations for driver
func (store *dbStore) TravelAgencyReports(month int, year int) []TravelAgencyReport {
	var travelAgencyCount int

	var whereClause string
	var sqlString string

	log.Printf("query travel agent report")

	if month > 0 && year > 0 {
		log.Printf("we have date criteria")

		whereClause = " where postponed = 0 and cancelled = 0 and " +
			"month(departuredate) = " + strconv.Itoa(month) + " and year(departuredate) = " + strconv.Itoa(year)

		sqlString = "select count(r.travelagencyid) from reservations r inner join travelagencies ta on r.travelagencyid = ta.travelagencyid" +
			"" + whereClause

		log.Printf("%s", sqlString)

		row, err := store.db.Query(sqlString)

		log.Printf("query executed")

		row.Next()
		err = row.Scan(
			&travelAgencyCount,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No travel agency report returned")
			} else {
				log.Printf("Error retrieving travel agency report count: %s", err.Error())
			}
		}

		log.Printf("travel agency report count: %d", travelAgencyCount)
		var travelAgencyReportSlice = make([]TravelAgencyReport, travelAgencyCount)

		if travelAgencyCount > 0 {
			row, err = store.db.Query("select r.travelagencyid, ta.travelagencyname, count(r.ReservationID) " +
				"reservationcount, sum(r.Price) totalcost, sum(r.Price)*0.10 as commission " +
				"from reservations r inner join travelagencies ta on r.travelagencyid = ta.travelagencyid " +
				"" + whereClause + " group by ta.travelagencyid, ta.travelagencyname")

			if err != nil {
				if err == sql.ErrNoRows {
					log.Print("No travel agency reports found")
				} else {
					log.Printf("Error retrieving travel agency records: %s", err.Error())
				}
				return []TravelAgencyReport{}
			}
			defer row.Close()

			log.Printf("retrieved records")

			var indx int
			indx = 0
			for row.Next() {
				err = row.Scan(
					&travelAgencyReportSlice[indx].TravelAgencyID, &travelAgencyReportSlice[indx].TravelAgencyName,
					&travelAgencyReportSlice[indx].ReservationCount, &travelAgencyReportSlice[indx].TotalCost,
					&travelAgencyReportSlice[indx].Commission,
				)

				if err != nil {
					log.Printf("Error retrieving travel agent reports: %s", err.Error())
				}

				log.Printf("travel agency: %s", travelAgencyReportSlice[indx].TravelAgencyName)

				indx++
			}
		}

		return travelAgencyReportSlice
	}

	log.Printf("no date criteria, return empty report")

	return []TravelAgencyReport{}
}

//GetClientFromReservation - return client info used to create reservation
func (store *dbStore) GetClientFromReservation(reservationid int) Client {

	row := store.db.QueryRow("select firstname, lastname, phone, email, streetaddress, city, "+
		"province, postalcode, country from clients c inner join reservations r "+
		"on c.clientid = r.clientid where r.reservationid = ?", reservationid)

	client := Client{}

	err := row.Scan(
		&client.Firstname, &client.Lastname, &client.Phone,
		&client.Email, &client.StreetAddress, &client.City,
		&client.Province, &client.PostalCode, &client.Country,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No client found")
		} else {
			log.Printf("Error retrieving client: %s", err.Error())
		}
	}

	return client
}

//AGTAQueryReport - Get Data for AGTA Report
func (store *dbStore) AGTAQueryReport(startdate time.Time, enddate time.Time) []AGTAReport {

	AGTACount := 0

	log.Printf("AGTA get query count")

	sqlString := "select count(reservationid) " +
		"FROM northernairport.reservations r JOIN northernairport.trips t ON r.tripid = t.tripid " +
		"LEFT JOIN northernairport.trips rt ON r.returntripid = rt.tripid " +
		"JOIN northernairport.clients c ON r.clientid = c.clientid " +
		"JOIN northernairport.vehicles v ON t.vehicleid = v.vehicleid " +
		"JOIN northernairport.departuretimes dt ON dt.departuretimeid = r.departuretimeid " +
		"LEFT JOIN northernairport.departuretimes rdt ON dt.departuretimeid = r.returndeparturetimeid " +
		"WHERE (cancelled is null or cancelled = 0) AND (departurecityid=2 AND " +
		" r.departuredate >= '" + startdate.Format("2006-01-02") + "' AND " +
		" r.departuredate < '" + enddate.Format("2006-01-02") + "') OR " +
		" (destinationcityid=2 AND r.returndate >= '" + startdate.Format("2006-01-02") + "' AND " +
		" r.returndate < '" + enddate.Format("2006-01-02") + "')"

	row, err := store.db.Query(sqlString)

	log.Printf("AGTA query count received")

	row.Next()
	err = row.Scan(
		&AGTACount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No AGTA report returned")
		} else {
			log.Printf("Error retrieving AGTA report count: %s", err.Error())
		}
	}

	log.Printf("AGTA retrieved records")

	log.Printf("AGTA report count: %d", AGTACount)
	var AGTAReportSlice = make([]AGTAReport, AGTACount)

	if AGTACount > 0 {
		sqlString =
			"SELECT r.reservationid, flighttime, " +
				"IF(departurecityid=2, " +
				"(SELECT name FROM airlines WHERE airlineid=departureairlineid), " +
				"(SELECT name FROM airlines WHERE airlineid=returnairlineid)) AS AirlineName, " +
				"flightnumber, " +
				"IF(departurecityid=2, " +
				"(select terminal FROM airlines WHERE airlineid=departureairlineid), " +
				"(select terminal FROM airlines WHERE airlineid=returnairlineid)) AS TerminalName, " +
				"CONCAT(lastname, ', ', firstname) AS PaxName, r.reservationid AS confirmationnumber, " +
				"departurenumchildren + departurenumstudents + departurenumadults + departurenumseniors as NumPax, " +
				"IF(departurecityid=2, " +
				"IF(destinationvenueid=100, dropoffaddress, " +
				"(SELECT name FROM venues WHERE venueid=destinationvenueid)), " +
				"IF(returndestinationvenueid=100, returndropoffaddress, " +
				"(SELECT name FROM venues WHERE venueid=returndestinationvenueid))) AS DropLocation, " +
				"IF(departurecityid=2, " +
				"(SELECT name FROM cities WHERE cityid=destinationcityid), " +
				"(SELECT name FROM cities WHERE cityid=returndestinationcityid)) AS DropCity, " +
				"internalnotes, drivernotes, " +
				"IF(departurecityid=2, dt.departuretime, rdt.departuretime) as departuretime, " +
				"IF(departurecityid=2, " +
				"((SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE driverid=t.driverid)), " +
				"((SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE driverid=rt.driverid))) AS DriverName, " +
				"IF(departurecityid=2, t.driverid, rt.driverid) As driverid, " +
				"IF(departurecityid=2, " +
				"((SELECT numseats FROM vehicles WHERE vehicleid=t.vehicleid)), " +
				"((SELECT numseats FROM vehicles WHERE vehicleid=rt.vehicleid))) AS VehicleNum, " +
				"(cancelled is null or cancelled = 0) AS IsValid, r.departuredate, " +
				"IF(departurecityid=2, " +
				"IF(departurevenueid=99, '', (SELECT name FROM venues WHERE venueid=departurevenueid)), " +
				"IF(returndeparturevenueid=99, '', (SELECT name FROM venues WHERE venueid=returndeparturevenueid))) AS HotelInfo, " +
				"cancelled " +
				"FROM northernairport.reservations r JOIN northernairport.trips t ON r.tripid = t.tripid " +
				"LEFT JOIN northernairport.trips rt ON r.returntripid = rt.tripid " +
				"JOIN northernairport.clients c ON r.clientid = c.clientid " +
				"JOIN northernairport.vehicles v ON t.vehicleid = v.vehicleid " +
				"JOIN northernairport.departuretimes dt ON dt.departuretimeid = r.departuretimeid " +
				"LEFT JOIN northernairport.departuretimes rdt ON rdt.departuretimeid = r.returndeparturetimeid " +
				"WHERE (cancelled is null or cancelled = 0) AND (departurecityid=2 AND " +
				" r.departuredate >= '" + startdate.Format("2006-01-02") + "' AND " +
				" r.departuredate < '" + enddate.Format("2006-01-02") + "') OR " +
				" (destinationcityid=2 AND r.returndate >= '" + startdate.Format("2006-01-02") + "' AND " +
				" r.returndate < '" + enddate.Format("2006-01-02") + "')"

		row, err = store.db.Query(sqlString)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Print("No AGTA report found")
			} else {
				log.Printf("Error retrieving AGTA records: %s", err.Error())
			}
			return []AGTAReport{}
		}
		defer row.Close()

		var indx int
		var flighttime sql.NullString
		var flightnumber sql.NullString
		var airlinename sql.NullString
		var terminalname sql.NullString
		var internalnotes sql.NullString
		var droplocation sql.NullString
		var drivernotes sql.NullString
		var drivername sql.NullString
		var hotelinfo sql.NullString
		var vehiclenum sql.NullString
		var departuredate mysql.NullTime

		indx = 0
		for row.Next() {
			err = row.Scan(
				&AGTAReportSlice[indx].ReservationID,
				&flighttime, &airlinename, &flightnumber, &terminalname,
				&AGTAReportSlice[indx].PaxName, &AGTAReportSlice[indx].ConfirmationNumber,
				&AGTAReportSlice[indx].NumPax, &droplocation,
				&AGTAReportSlice[indx].DropCity, &internalnotes, &drivernotes,
				&AGTAReportSlice[indx].DepartureTime, &drivername,
				&AGTAReportSlice[indx].DriverID, &vehiclenum,
				&AGTAReportSlice[indx].IsValid, &departuredate,
				&hotelinfo, &AGTAReportSlice[indx].Cancelled,
			)

			if err != nil {
				log.Printf("Error retrieving AGTA reports: %s", err.Error())
			}

			if len(flighttime.String) > 0 {
				AGTAReportSlice[indx].FlightTime = flighttime.String
			}

			if len(flightnumber.String) > 0 {
				AGTAReportSlice[indx].FlightNumber = flightnumber.String
			}

			if len(airlinename.String) > 0 {
				AGTAReportSlice[indx].AirlineName = airlinename.String
			}

			if len(terminalname.String) > 0 {
				AGTAReportSlice[indx].TerminalName = terminalname.String
			}

			if len(internalnotes.String) > 0 {
				AGTAReportSlice[indx].InternalNotes = internalnotes.String
			}

			if len(droplocation.String) > 0 {
				AGTAReportSlice[indx].DropLocation = droplocation.String
			}

			if len(drivernotes.String) > 0 {
				AGTAReportSlice[indx].DriverNotes = drivernotes.String
			}

			if len(drivername.String) > 0 {
				AGTAReportSlice[indx].DriverName = drivername.String
			}

			if len(hotelinfo.String) > 0 {
				AGTAReportSlice[indx].HotelInfo = hotelinfo.String
			}

			if len(vehiclenum.String) > 0 {
				AGTAReportSlice[indx].VehicleNum = vehiclenum.String
			}

			if departuredate.Valid {
				AGTAReportSlice[indx].DepartureDate = departuredate.Time
			} else {
				AGTAReportSlice[indx].DepartureDate = time.Time{}
			}

			indx++
		}
	}

	return AGTAReportSlice
}

//UpdateStatus - Update Reservation Status
func (store *dbStore) UpdateStatus(reservationid int, elavontransactionid string, status string) error {
	log.Printf("update reservation status: %s in database", status)

	if status != "approval" && status != "declined" && status != "error" {
		log.Printf("Invalid status received: %s", status)
	} else if status == "approval" {
		status = "approved"
	}

	_, updateerr := store.db.Exec("UPDATE reservations SET status = ?, elavontransactionid = ? "+
		"WHERE reservationid = ?", status, elavontransactionid, reservationid)

	if updateerr != nil {
		log.Printf("Error updating price: %s", updateerr.Error())
	} else {
		log.Printf("Update Status: %s", status)
	}

	return updateerr
}

//GetReservationAmount - Get Reservation Amount
func (store *dbStore) GetReservationAmount(reservationid int) float32 {
	row := store.db.QueryRow("select price from reservations "+
		"where reservationid = ?", reservationid)

	price := 0.0

	err := row.Scan(
		&price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No reservation found")
		} else {
			log.Printf("Error retrieving reservation: %s", err.Error())
		}
	}

	return float32(price)
}

func (store *dbStore) MigrateDB() {
	log.Printf("migate database")

	log.Printf("delete existing records")

	const layoutISO = "2006-01-02"

	//delete clients, update auto increment
	{
		_, updateerr := store.db.Exec("DELETE FROM clients WHERE ClientID > 7")

		if updateerr != nil {
			log.Printf("Error deleting clients: %s", updateerr.Error())
		} else {
			log.Printf("Clients deleted")
		}

		log.Printf("update auto increment to 8")

		_, updateerr = store.db.Exec("ALTER TABLE clients AUTO_INCREMENT = 8")

		if updateerr != nil {
			log.Printf("Error updating clients auto increment: %s", updateerr.Error())
		} else {
			log.Printf("Clients auto increment updated")
		}
	}

	//delete reservations, update auto increment
	{
		_, updateerr := store.db.Exec("DELETE FROM reservations WHERE ReservationID > 2")

		if updateerr != nil {
			log.Printf("Error deleting reservations: %s", updateerr.Error())
		} else {
			log.Printf("Reservations deleted")
		}

		log.Printf("update auto increment to 3")

		_, updateerr = store.db.Exec("ALTER TABLE reservations AUTO_INCREMENT = 3")

		if updateerr != nil {
			log.Printf("Error updating reservations auto increment: %s", updateerr.Error())
		} else {
			log.Printf("Reservations auto increment updated")
		}
	}

	//add trips for each day starting 01/01/2016
	{
		log.Printf("Creating Trips")

		datecounter, err := time.Parse("2006-01-02", "2020-01-01")

		if err != nil {
			log.Printf("Error creating start date: %s", err.Error())
		}

		today, err := time.Parse("2006-01-02", time.Now().Format(layoutISO))

		if err != nil {
			log.Printf("Error creating current date: %s", err.Error())
		}

		for datecounter.Before(today) {

			_, triperr := store.db.Exec("INSERT INTO trips (departuredate, departuretimeid, "+
				"numpassengers, driverid, vehicleid, capacity, omitted) "+
				"VALUES	(?, 1, 0, 0, 0, 11, 0), "+
				" (?, 2, 0, 0, 0, 11, 0), "+
				" (?, 3, 0, 0, 0, 11, 0), "+
				" (?, 4, 0, 0, 0, 11, 0), "+
				" (?, 5, 0, 0, 0, 11, 0), "+
				" (?, 6, 0, 0, 0, 11, 0), "+
				" (?, 7, 0, 0, 0, 11, 0), "+
				" (?, 8, 0, 0, 0, 11, 0) ",
				datecounter.Format(layoutISO), datecounter.Format(layoutISO), datecounter.Format(layoutISO),
				datecounter.Format(layoutISO), datecounter.Format(layoutISO), datecounter.Format(layoutISO),
				datecounter.Format(layoutISO), datecounter.Format(layoutISO))

			if triperr != nil {
				log.Printf("Error inserting trip for day: %s Error: %s", datecounter.Format(layoutISO), triperr.Error())
				//log.Printf(sqlstmt)
			}

			datecounter = datecounter.Add(24 * time.Hour)
		}

		log.Printf("Trips Created")
	}

	//migrate clients and reservations
	{
		log.Printf("Migrating Clients")

		datecounter, err := time.Parse("2006-01-02", "2020-01-01")

		if err != nil {
			log.Printf("Error creating start date: %s", err.Error())
		}

		today, err := time.Parse("2006-01-02", time.Now().Format(layoutISO))

		if err != nil {
			log.Printf("Error creating current date: %s", err.Error())
		}

		for datecounter.Before(today) {

			plusonemonth := datecounter.AddDate(0, 1, 0)

			result, clienterr := store.db.Exec("INSERT INTO clients (firstname, lastname, phone, "+
				"email, streetaddress, city, province, postalcode, country) "+
				"SELECT firstname, lastname, CASE WHEN REGEXP_REPLACE(phone, '[^0-9]','') = '' THEN 0 "+
				"ELSE LEFT(REGEXP_REPLACE(phone, '[^0-9]',''),10) END AS phone, "+
				"email, billingaddress, city, province, postal, country "+
				"FROM reservations_old "+
				"WHERE departuredate >= ? AND "+
				"departuredate <= ? ", datecounter.Format(layoutISO), plusonemonth.Format(layoutISO))

			if clienterr != nil {
				log.Printf("Error migrating legacy client: %s Error: %s", datecounter.Format(layoutISO), clienterr.Error())
			}

			id, _ := result.LastInsertId()

			_, reservationerr := store.db.Exec("INSERT INTO reservations (clientid, travelagencyid, departurecityid, "+
				"departurevenueid, departuretimeid, destinationcityid, destinationvenueid, returndeparturecityid, "+
				"returndeparturevenueid, returndeparturetimeid, returndestinationcityid, returndestinationvenueid, "+
				"discountcodeid, departureairlineid, returnairlineid, drivernotes, internalnotes, departurenumadults, "+
				"departurenumstudents, departurenumseniors, departurenumchildren, returnnumadults, returnnumstudents, "+
				"returnnumseniors, returnnumchildren, price, customdepartureid, customdestinationid, departuredate, "+
				"returndate, triptypeid, balanceowing, elavontransactionid, postponed, cancelled, tripid) "+
				"SELECT "+strconv.Itoa(int(id))+", 1, departurecity, departurevenue, 1, "+
				"destinationcity, destinationvenue, destinationcity, IF(returndeparturevenue LIKE '', null, returndeparturevenue), REGEXP_REPLACE(returntime, '[^0-9]',null), "+
				"IF(returndestinationcity LIKE '', null, returndestinationcity), IF(returndestinationvenue LIKE '', null, returndestinationvenue), 1, IF(airline LIKE '', null, airline), "+
				"IF(returnairline LIKE '', null, returnairline),	drivernotes, tanotes, departureadults, departurestudents, departureseniors, departurechildren, "+
				"IF(returnadults LIKE '', null, returnadults), IF(returnstudents LIKE '', null, returnstudents), IF(returnseniors LIKE '', null, returnseniors), "+
				"IF(returnchildren LIKE '', null, returnchildren), cost, 0, 0, departuredate, IF(returndate LIKE '', null, returndate), IF(triptype LIKE 'oneway', 1, 2), balanceowing, elavontransactionid, "+
				"0, 0, 1 FROM reservations_old r "+
				"WHERE departuredate >= ? "+
				"AND departuredate <= ?", datecounter.Format(layoutISO), plusonemonth.Format(layoutISO))

			if reservationerr != nil {
				log.Printf("Error migrating legacy reservation: %s Error: %s", datecounter.Format(layoutISO), reservationerr.Error())
			}

			datecounter = plusonemonth
		}

	}
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/*
InitStore We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}
