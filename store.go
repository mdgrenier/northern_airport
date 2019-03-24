package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//Store offers an interface for various db function to the rest of the application
type Store interface {
	CreateUser(client *Client) error
	SignInUser(client *Client) error
	GetVenues() []Venues
	GetVenueCount() int
	GetCities() []Cities
	GetCityCount() int
	GetClientInfo(client *Client) error
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateUser(client *Client) error {
	result, err := store.db.Exec("INSERT INTO clients(firstname, lastname, phone, email, "+
		"streetaddress, city, province, postalcode, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		client.Firstname, client.Lastname, client.Phone, client.Email, client.StreetAddress,
		client.City, client.Province, client.PostalCode, client.Country)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	log.Printf("New Client ID: %d", id)

	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) "+
		"VALUES (?, ?, 2, ?)",
		id, client.Password, client.Username)

	id, _ = result.LastInsertId()

	log.Printf("Account Details: %d", id)

	return err
}

func (store *dbStore) CreateReservation(client *Client) error {
	result, err := store.db.Exec("INSERT INTO clients(firstname, lastname, phone, email, streetaddress, "+
		"city, province, postalcode, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		client.Firstname, client.Lastname, client.Phone, client.Email,
		client.StreetAddress, client.City, client.Province, client.PostalCode, client.Country)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	log.Printf("New Client ID: %d", id)

	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) "+
		"VALUES (?, ?, 2, ?)",
		id, client.Password, client.Username)

	id, _ = result.LastInsertId()

	log.Printf("Account Details: %d", id)

	return err
}

func (store *dbStore) SignInUser(client *Client) error {
	// Parse and decode the request body into a new `Credentials` instance

	// Query the database for all birds, and return the result to the `rows` object
	//rows, err := store.db.Query("SELECT username, password from users")
	row, err := store.db.Query("select password from accountdetails where username=?", client.Username)
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return err
	}
	defer row.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	storedClient := &Client{}

	row.Next()
	err = row.Scan(&storedClient.Password)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Print("Unauthorized")
			//w.WriteHeader(http.StatusUnauthorized)
			//return
		} else {
			log.Printf("Sign In Error: %s", err.Error())
			//log.Fatal("Error: ", err)
		}
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedClient.Password), []byte(client.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		//w.WriteHeader(http.StatusUnauthorized)
		log.Print("Incorrect Password")
	}

	return err
}

func (store *dbStore) GetVenues() []Venues {

	venueCount := store.GetVenueCount()

	row, err := store.db.Query("select venueid, name from venues")
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return nil
	}
	defer row.Close()

	var venueSlice = make([]Venues, venueCount)

	var venueid int
	var name string
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&venueid, &name,
		)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No venue found")
			} else {
				log.Printf("Error retrieving venues: %s", err.Error())
			}
		} else {
			venueSlice[indx].VenueID = venueid
			venueSlice[indx].VenueName = name
		}

		indx++
	}

	return venueSlice
}

func (store *dbStore) GetVenueCount() int {
	var venueCount int

	row, err := store.db.Query("select count(venueid) from venues")

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

func (store *dbStore) GetCities() []Cities {
	cityCount := store.GetCityCount()

	row, err := store.db.Query("select cityid, name from cities")
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return nil
	}
	defer row.Close()

	var citySlice = make([]Cities, cityCount)

	var cityid int
	var name string
	var indx int

	indx = 0
	for row.Next() {
		err = row.Scan(
			&cityid, &name,
		)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				log.Print("No cities found")
			} else {
				log.Printf("Error retrieving cities: %s", err.Error())
			}
		} else {
			citySlice[indx].CityID = cityid
			citySlice[indx].CityName = name
		}

		indx++
	}

	return citySlice
}

func (store *dbStore) GetCityCount() int {
	var cityCount int

	row, err := store.db.Query("select count(cityid) from cities")

	row.Next()
	err = row.Scan(
		&cityCount,
	)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Print("No venues returned")
		} else {
			log.Printf("Error retrieving venues: %s", err.Error())
		}
	}

	return cityCount
}

//GetClientInfo takes a client and username
func (store *dbStore) GetClientInfo(client *Client) error {

	// Query the database for all birds, and return the result to the `rows` object
	row, err := store.db.Query(
		"select firstname, lastname, phone, email, streetaddress, "+
			"city, province, postalcode, country from clients c inner join "+
			"accountdetails a on c.clientid = a.clientid "+
			"where a.username=?", client.Username)
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		log.Printf("Error retrieving client: %s", err.Error())
		return err
	}
	defer row.Close()

	var firstname string
	var lastname string
	var phone string
	var email string
	var streetaddress string
	var city string
	var province string
	var postalcode string
	var country string

	row.Next()
	err = row.Scan(
		&firstname, &lastname, &phone, &email, &streetaddress,
		&city, &province, &postalcode, &country,
	)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			log.Print("No client found")
		} else {
			log.Printf("Error saving client: %s", err.Error())
		}
	} else {
		//clients.Username =
		//clients.Password = firstname
		client.Firstname = firstname
		client.Lastname = lastname
		client.Phone = phone
		client.Email = email
		client.StreetAddress = streetaddress
		client.City = city
		client.Province = province
		client.PostalCode = postalcode
		client.Country = country
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
