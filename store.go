package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"log"
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
	GetVenueCount() int
	GetCities() []Cities
	GetCityCount() int
	GetDepartureTimesCount() int
	GetDepartureTimes() []DepartureTimes
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

	//log.Printf("New Client ID: %d", id)

	//create account details record linked to client record
	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) "+
		"VALUES (?, ?, 2, ?)",
		id, client.Password, client.Username)

	//id, _ = result.LastInsertId()

	//log.Printf("Account Details: %d", id)

	return err
}

//CreateReservation - store new reservation in database
func (store *dbStore) CreateReservation(reservation *Reservation) error {

	//temporarily hardcoding some values until they are either implemented or removed
	reservation.DiscountCodeID = 1
	reservation.DepartureAirlineID = 1
	reservation.ReturnAirlineID = 1
	reservation.Status = ""
	reservation.Hash = ""
	reservation.ElavonTransactionID = 0

	log.Print("Inserting reservation into DB")

	var insertError error
	var id int64

	if reservation.TripTypeID == 2 {
		result, err := store.db.Exec("INSERT INTO reservations("+
			"clientid, reservationtypeid, departurecityid, departurevenueid, departuretimeid, "+
			"destinationcityid, destinationvenueid, returndeparturecityid, returndeparturevenueid, returndeparturetimeid, "+
			"returndestinationcityid, returndestinationvenueid, discountcodeid, departureairlineid, returnairlineid, "+
			"drivernotes, internalnotes, departurenumadults, departurenumstudents, departurenumseniors, "+
			"departurenumchildren, returnnumadults, returnnumstudents, returnnumseniors, returnnumchildren, "+
			"price, status, hash, customdepartureid, customdestinationid, "+
			"departuredate, returndate, triptypeid, balanceowing, elavontransactionid) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?)",
			reservation.ClientID, reservation.ReservationTypeID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.ReturnDepartureCityID, reservation.ReturnDepartureVenueID, reservation.ReturnDepartureTimeID,
			reservation.ReturnDestinationCityID, reservation.ReturnDestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID, reservation.ReturnAirlineID,
			reservation.DriverNotes, reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors,
			reservation.DepartureNumChildren, reservation.ReturnNumAdults, reservation.ReturnNumStudents, reservation.ReturnNumSeniors, reservation.ReturnNumChildren,
			reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID, reservation.CustomDepartureID,
			reservation.DepartureDate, reservation.ReturnDate, reservation.TripTypeID, reservation.BalanceOwing, reservation.ElavonTransactionID)

		insertError = err

		id, _ = result.LastInsertId()
	} else {
		result, err := store.db.Exec("INSERT INTO reservations("+
			"clientid, reservationtypeid, departurecityid, departurevenueid, departuretimeid, "+
			"destinationcityid, destinationvenueid, discountcodeid, departureairlineid, "+
			"drivernotes, internalnotes, departurenumadults, departurenumstudents, departurenumseniors, "+
			"departurenumchildren, price, status, hash, customdepartureid, "+
			"customdestinationid, departuredate, triptypeid, balanceowing, elavontransactionid) VALUES "+
			"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?, ?, ?, ?, ?, ?, ?, "+
			"?, ?, ?, ?)",
			reservation.ClientID, reservation.ReservationTypeID, reservation.DepartureCityID, reservation.DepartureVenueID, reservation.DepartureTimeID,
			reservation.DestinationCityID, reservation.DestinationVenueID, reservation.DiscountCodeID, reservation.DepartureAirlineID,
			reservation.DriverNotes, reservation.InternalNotes, reservation.DepartureNumAdults, reservation.DepartureNumStudents, reservation.DepartureNumSeniors,
			reservation.DepartureNumChildren, reservation.Price, reservation.Status, reservation.Hash, reservation.CustomDepartureID,
			reservation.CustomDepartureID, reservation.DepartureDate, reservation.TripTypeID, reservation.BalanceOwing, reservation.ElavonTransactionID)

		insertError = err

		id, _ = result.LastInsertId()
	}

	log.Printf("New Reservation ID: %d", id)
	return insertError
}

//SignInUser - authenticate user
func (store *dbStore) SignInUser(client *Client) error {

	row, err := store.db.Query("select password from accountdetails where username=?", client.Username)

	if err != nil {
		return err
	}
	defer row.Close()

	//create a client object
	storedClient := &Client{}

	//get returned client record if one exists
	row.Next()
	err = row.Scan(&storedClient.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Unauthorized")
			//***handle this better, display credential error***
		} else {
			log.Printf("Sign In Error: %s", err.Error())
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedClient.Password), []byte(client.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		log.Print("Incorrect Password")
	}

	return err
}

//GetVenues - return all venues in database
func (store *dbStore) GetVenues() []Venues {

	//Get venue count
	venueCount := store.GetVenueCount()

	row, err := store.db.Query("select venueid, cityid, name from venues")
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
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
			&venueSlice[indx].VenueID, &venueSlice[indx].CityID, &venueSlice[indx].VenueName,
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
			log.Printf("Error retrieving venues: %s", err.Error())
		}
	}

	return venueCount
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
			log.Print("No venues returned")
		} else {
			log.Printf("Error retrieving venues: %s", err.Error())
		}
	}

	return cityCount
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
		}

		indx++
	}

	return departureTimesSlice
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
			log.Print("No venues returned")
		} else {
			log.Printf("Error retrieving venues: %s", err.Error())
		}
	}

	return departuretimeCount
}

//GetClientInfo - from client username return all client info
func (store *dbStore) GetClientInfo(client *Client) error {

	row, err := store.db.Query(
		"select c.clientid, firstname, lastname, phone, email, streetaddress, "+
			"city, province, postalcode, country from clients c inner join "+
			"accountdetails a on c.clientid = a.clientid "+
			"where a.username=?", client.Username)
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		log.Printf("Error retrieving client: %s", err.Error())
		return err
	}
	defer row.Close()

	//store client into into local variables
	row.Next()
	err = row.Scan(
		&client.ClientID, &client.Firstname, &client.Lastname, &client.Phone, &client.Email,
		&client.StreetAddress, &client.City, &client.Province, &client.PostalCode, &client.Country,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("No client found")
		} else {
			log.Printf("Error saving client: %s", err.Error())
		}
	}

	return err
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
