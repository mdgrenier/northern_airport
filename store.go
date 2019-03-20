package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//Store offers an interface for various db function to the rest of the application
type Store interface {
	CreateUser(clients *Clients) error
	SignInUser(clients *Clients) error
	GetVenues() error
	GetClientInfo(clients *Clients, username string) error
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateUser(clients *Clients) error {
	result, err := store.db.Exec("INSERT INTO clients(firstname, lastname, phone, email, streetaddress, city, province, postalcode, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		clients.Firstname, clients.Lastname, clients.Phone, clients.Email, clients.StreetAddress, clients.City, clients.Province, clients.PostalCode, clients.Country)

	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()

	log.Printf("New Client ID: %d", id)

	result, err = store.db.Exec("INSERT INTO accountdetails(clientid, password, roleid, username) VALUES (?, ?, 2, ?)",
		id, clients.Password, clients.Username)

	id, _ = result.LastInsertId()

	log.Printf("Account Details: %d", id)

	return err
}

func (store *dbStore) SignInUser(clients *Clients) error {
	// Parse and decode the request body into a new `Credentials` instance

	// Query the database for all birds, and return the result to the `rows` object
	//rows, err := store.db.Query("SELECT username, password from users")
	row, err := store.db.Query("select password from accountdetails where username=?", clients.Username)
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return err
	}
	defer row.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	storedClients := &Clients{}

	row.Next()
	err = row.Scan(&storedClients.Password)
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
	if err = bcrypt.CompareHashAndPassword([]byte(storedClients.Password), []byte(clients.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		//w.WriteHeader(http.StatusUnauthorized)
		log.Print("Incorrect Password")
	}

	return err
}

func (store *dbStore) GetVenues() error {
	row, err := store.db.Query("select venueid, name from venues")
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return err
	}
	defer row.Close()

	return err
}

//GetClientInfo takes a username and returns client info
func (store *dbStore) GetClientInfo(clients *Clients, username string) error {

	log.Printf("Username: %s", username)

	// Query the database for all birds, and return the result to the `rows` object
	row, err := store.db.Query(
		"select firstname, lastname, phone, email, streetaddress, "+
			"city, province, postalcode, country from clients c inner join "+
			"accountdetails a on c.clientid = a.clientid "+
			"where a.username=?", username)
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
		clients.Firstname = firstname
		clients.Lastname = lastname
		clients.Phone = phone
		clients.Email = email
		clients.StreetAddress = streetaddress
		clients.City = city
		clients.Province = province
		clients.PostalCode = postalcode
		clients.Country = country
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
