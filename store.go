package main

// The sql go library is needed to interact with the database
import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// Store will have two methods, to add a new bird, and to get all existing birds
// Each method returns an error, in case something goes wrong
type Store interface {
	CreateUser(creds *Credentials) error
	SignInUser(creds *Credentials) error
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateUser(creds *Credentials) error {
	_, err := store.db.Query("INSERT INTO accountdetails(clientid, password, roleid, username) VALUES (3, ?, 1, ?)", creds.Password, creds.Username)
	return err
}

func (store *dbStore) SignInUser(creds *Credentials) error {
	// Parse and decode the request body into a new `Credentials` instance

	// Query the database for all birds, and return the result to the `rows` object
	//rows, err := store.db.Query("SELECT username, password from users")
	row, err := store.db.Query("select password from accountdetails where username=?", creds.Username)
	// We return in case of an error, and defer the closing of the row structure
	if err != nil {
		return err
	}
	defer row.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	storedCreds := &Credentials{}

	row.Next()
	err = row.Scan(&storedCreds.Password)
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
	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		//w.WriteHeader(http.StatusUnauthorized)
		log.Print("Incorrect Password")
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
