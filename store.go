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
	CreateUser(client *Client) error
	SignInUser(client *Client) error
	CreateReservation(reservation *Reservation) error
	GetClientInfo(client *Client) error
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
	PostponeReservation(searchreservation *SearchReservations) error
	CancelReservation(searchreservation *SearchReservations) error
	GetDrivers() []Drivers
	AddDriver(driver Drivers) error
	UpdateDriver(driver *Drivers) error
	DeleteDriver(driver int) error
	GetVehicles() []Vehicles
	AddVehicle(vehicle Vehicles) error
	UpdateVehicle(vehicle *Vehicles) error
	DeleteVehicle(vehicle int) error
	GetPrice(departurecityid int, destinationcityid int, retdeparturecityid int, retdestinationcityid int, customertypeid int, reservationtypeid int, discountcode string) float32
	GetPrices() []Prices
	UpdatePrice(price *Prices) error
	GetDiscount(discountcode *DiscountCode) error
	AddVenueFee(departurevenueid int, destinationvenueid int, retdeparturevenueid int, retdestinationvenueid int) float32
}

//The `dbStore` struct will implement the `Store` interface it also takes the sql
//DB connection object, which represents the database connection.
type dbStore struct {
	db *sql.DB
}

//CreateUser - store new client in database
func (store *dbStore) CreateUser(client *Client) error {

	result, err := store.db.Exec("INSERT INTO clients(firstname, lastname, phone, email, "+
		"streetaddress, city, province, postalcode, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		client.Firstname, client.Lastname, client.Phone, client.Email, client.StreetAddress,
		client.City, client.Province, client.PostalCode, client.Country)

	if err != nil {
		return err
	}

	//get id from client insertion transaction
	id, _ := result.LastInsertId()

	//create account details record linked to client record
	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) "+
		"VALUES (?, ?, 2, ?)",
		id, client.Password, client.Username)

	return err
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
			"departuredate, returndate, triptypeid, tripid, balanceowing, elavontransactionid) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?)",
			reservation.ClientID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.ReturnDepartureCityID, reservation.ReturnDepartureVenueID, reservation.ReturnDepartureTimeID,
			reservation.ReturnDestinationCityID, reservation.ReturnDestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID, reservation.ReturnAirlineID,
			reservation.DriverNotes, reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors,
			reservation.DepartureNumChildren, reservation.ReturnNumAdults, reservation.ReturnNumStudents, reservation.ReturnNumSeniors, reservation.ReturnNumChildren,
			reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID, reservation.CustomDepartureID,
			reservation.DepartureDate, reservation.ReturnDate, reservation.TripTypeID, reservation.TripID, reservation.BalanceOwing, reservation.ElavonTransactionID)

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
			"departuredate, triptypeid, tripid, balanceowing, elavontransactionid) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?)",
			reservation.ClientID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID, reservation.DriverNotes,
			reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors, reservation.DepartureNumChildren,
			reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID, reservation.CustomDepartureID,
			reservation.DepartureDate, reservation.TripTypeID, reservation.TripID, reservation.BalanceOwing, reservation.ElavonTransactionID)

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
	destinationvenue := store.GetCityName(reservation.DestinationVenueID)
	numadults := strconv.Itoa(reservation.DepartureNumAdults)
	numseniors := strconv.Itoa(reservation.DepartureNumSeniors)
	numstudents := strconv.Itoa(reservation.DepartureNumStudents)
	numchildren := strconv.Itoa(reservation.DepartureNumChildren)

	log.Printf("Total Reservation Cost: $%f", reservation.Price)

	departuredetails = FormatTripDetails(departurecity, departurevenue, departuredate, departuretime,
		destinationcity, destinationvenue, numadults, numseniors, numstudents, numchildren)

	var client Client
	client.ClientID = reservation.ClientID

	store.GetClientInfo(&client)

	SendConfirmationEmail(client.Email, departuredetails, returndetails, reservation.ReservationID, reservation.Price)

	return insertError
}

//SignInUser - authenticate user
func (store *dbStore) SignInUser(client *Client) error {

	//create a client object
	storedClient := &Client{}

	err := store.db.QueryRow("select password, roleid from accountdetails where username=?",
		client.Username).Scan(&storedClient.Password, &client.RoleID)

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
		log.Printf("Error retrieving client: %s", err.Error())
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
			log.Printf("Error retrieving client: %s", err.Error())
		}
	}

	return err
}

//GetOrAddTrip - store appropriate tripid in reservation, if no trip
//				 generate new trip
func (store *dbStore) GetOrAddTrip(reservation *Reservation) error {
	var tripid int64
	var numtrippassengers int
	var numpassengers int
	var capacity int
	var omitted int

	numpassengers = reservation.DepartureNumAdults + reservation.DepartureNumChildren + reservation.DepartureNumSeniors + reservation.DepartureNumStudents

	err := store.db.QueryRow("select tripid, numpassengers, capacity, omitted, from trips where tripid=?",
		reservation.TripID).Scan(&tripid, &numtrippassengers, &capacity, &omitted)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No trip, create one")

			//check if trip in omitted table

			result, inserterr := store.db.Exec("INSERT INTO trips(departuredate, departuretimeid, numpassengers) "+
				"VALUES (?,?,?)", reservation.DepartureDate, reservation.DepartureTimeID, numpassengers)

			if inserterr != nil {
				log.Printf("Error inserting new trip: %s", inserterr.Error())
			} else {
				tripid, _ = result.LastInsertId()
				log.Printf("TripID: %d", tripid)
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
				numpassengers += numtrippassengers

				_, err = store.db.Exec("UPDATE trips(numpassengers)"+
					"VALUES (?) WHERE tripid = ?", numpassengers, tripid)

				if err != nil {
					log.Printf("Error updating trip passengers: %s", err.Error())
				}

				log.Print("updated passengers num")

			}
		}

		//***if over 75% need notification***
	}

	reservation.TripID = int(tripid)

	return err
}

//GetTrips - return all trips (must add parameter to return by date)
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

	row, err = store.db.Query("select tripid, departuredate, t.departuretimeid, " +
		"numpassengers, driverid, vehicleid, capacity, " +
		"omitted from trips t inner join " +
		"departuretimes dt on t.departuretimeid = dt.departuretimeid ")

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
			&tripSlice[indx].Capacity, &tripSlice[indx].Omitted,
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

			//populate departuretime
			{

				dtrow, dterr := store.db.Query("select departuretime "+
					"from departuretimes where departuretimeid = ? ",
					tripSlice[indx].DepartureTimeID)

				if dterr != nil {
					log.Printf("Error retrieving departuretime: %s", dterr.Error())
					return nil
				}

				dtrow.Next()
				dterr = dtrow.Scan(
					&tripSlice[indx].DepartureTime,
				)

				if dterr != nil {
					if dterr == sql.ErrNoRows {
						log.Print("No departuretime found")
					} else {
						log.Printf("Error retrieving departuretime: %s", dterr.Error())
					}
				}
				dtrow.Close()
			}

			//populate drivers
			tripSlice[indx].DriverList = make([]Drivers, store.GetDriverCount())
			tripSlice[indx].DriverList = store.GetDrivers()

			//populate vehicle
			tripSlice[indx].VehicleList = make([]Vehicles, store.GetVehicleCount())
			tripSlice[indx].VehicleList = store.GetVehicles()

		}

		indx++
	}

	return tripSlice
}

//UpdateTrip - update driver and vehicle associated with trip
func (store *dbStore) UpdateTrip(trip *Trips) error {

	row := store.db.QueryRow("SELECT numseats from vehicles where vehicleid = ?",
		trip.VehicleID)

	var capacity int

	err := row.Scan(&capacity)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No capacity matches the vehicle id")
		} else {
			log.Printf("Error retrieving capacity from vehicle id: %s", err.Error())
		}
	}
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
	}

	if phone > 0 {
		log.Printf("add phone to where")
		whereClause += " c.Phone = " + strconv.Itoa(phone)
		addWhere = true
	}

	if len(email) > 0 {
		log.Printf("add email to where")
		whereClause += " c.Email = '" + email + "'"
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
			"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid " + whereClause)
	} else {
		row, err = store.db.Query("select r.reservationid, concat(c.firstname, ' ', c.lastname) as clientname, c.phone, c.email, " +
			"depv.name as departurevenue, desv.name as destinationvenue, r.triptypeid, " +
			"r.departurenumadults, r.departurenumstudents, r.departurenumseniors, r.departurenumchildren, " +
			"r.departuredate, dt.departuretime, r.postponed, r.cancelled " +
			"from reservations r inner join clients c on r.clientid = c.clientid " +
			"inner join venues depv on r.departurevenueid = depv.venueid " +
			"inner join venues desv on r.destinationvenueid = desv.venueid " +
			"inner join departuretimes dt on r.departuretimeid = dt.departuretimeid ")
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
			log.Printf("setting date for record: %d", indx)
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

	return driverCount
}

//GetDrivers - return all drivers
func (store *dbStore) GetDrivers() []Drivers {

	row, err := store.db.Query("select driverid, firstname, lastname " +
		"from drivers")

	if err != nil {
		log.Printf("Error retrieving drivers: %s", err.Error())
		return nil
	}
	defer row.Close()

	//create slice to store all departure times
	var driverSlice = make([]Drivers, store.GetDriverCount())

	var indx int
	indx = 0
	for row.Next() {
		err = row.Scan(
			&driverSlice[indx].DriverID, &driverSlice[indx].FirstName,
			&driverSlice[indx].LastName,
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

	return airlineSlice
}

//GetPrice - return price for a trip
func (store *dbStore) GetPrice(departurecityid int, destinationcityid int, retdeparturecityid int, retdestinationcityid int, customertypeid int, reservationtypeid int, discountcode string) float32 {
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

	if (departurecityid != retdeparturecityid || destinationcityid != retdestinationcityid) && reservationtypeid == 2 {
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

	} else if reservationtypeid == 2 {
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

//GetReports - return list of reports
func (store *dbStore) GetReports() {

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
