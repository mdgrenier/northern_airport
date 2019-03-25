package main

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// Client - data structure to hold client information and account details
type Client struct {
	Authenticated bool
	Password      string `json:"password" db:"password"`
	Username      string `json:"username" db:"username"`
	Firstname     string `json:"firstname" db:"firstname"`
	Lastname      string `json:"lastname" db:"lastname"`
	Phone         string `json:"phone" db:"phone"`
	Email         string `json:"email" db:"email"`
	StreetAddress string `json:"streetaddress" db:"streetaddress"`
	City          string `json:"city" db:"city"`
	Province      string `json:"province" db:"province"`
	PostalCode    string `json:"postalcode" db:"postalcode"`
	Country       string `json:"country" db:"country"`
}

// Venues - stores venues
type Venues struct {
	VenueID   int    `json:"venueid" db:"venueid"`
	CityID    int    `json:"cityid" db:"cityid"`
	VenueName string `json:"venuename" db:"venuename"`
}

// Cities - stores cities
type Cities struct {
	CityID   int    `json:"cityid" db:"cityid"`
	CityName string `json:"cityname" db:"cityname"`
}

// Reservation - store values to populate reservations form
type Reservation struct {
	Client     Client   `json:"client" db:"client"`
	Venues     []Venues `json:"venues" db:"venues"`
	VenueCount int      `json:"venuecount" db:"venuecount"`
	Cities     []Cities `json:"cities" db:"cities"`
}

var sessionStore *sessions.CookieStore

// tpl holds all parsed templates
var tpl *template.Template

func newRouter() *mux.Router {
	r := mux.NewRouter()
	// Declare the static file directory and point it to the directory we just made
	staticFileDirectory := http.Dir("./assets/")
	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")

	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/signup", signupHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/reservation", reservationHandler)

	r.HandleFunc("/signin", signinHandler).Methods("POST")
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/createreservation", createreservationHandler)

	return r
}

func main() {
	r := newRouter()

	connString := "root:test@tcp(db:3306)/northernairport"
	db, err := sql.Open("mysql", connString)

	if err != nil {
		panic(err)
	}
	//err = db.Ping()
	//
	//if err != nil {
	//	panic(err)
	//}

	InitStore(&dbStore{db: db})

	InitSession()

	http.ListenAndServe(":8080", r)
}

// InitSession initializes a user session
func InitSession() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	sessionStore = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	sessionStore.Options = &sessions.Options{
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(Client{})

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
