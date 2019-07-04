package main

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

// Client - data structure to hold client information and account details
type Client struct {
	Authenticated bool
	RoleID        int    `json:"roleid" db:"roleid"`
	ClientID      int    `json:"clientid" db:"clientid"`
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
	VenueID   int `json:"venueid" db:"venueid"`
	CityID    int `json:"cityid" db:"cityid"`
	CityName  string
	VenueName string `json:"venuename" db:"venuename"`
	ExtraCost int    `json:"extracost" db:"extracost"`
	Active    int    `json:"active" db:"active"`
	ExtraTime int    `json:"extratime" db:"extratime"`
	RoleID    int
}

// Cities - stores cities
type Cities struct {
	CityID      int    `json:"cityid" db:"cityid"`
	CityName    string `json:"cityname" db:"cityname"`
	NorthOffset int    `json:"northoffset" db:"northoffset"`
	SouthOffset int    `json:"southoffset" db:"southoffset"`
	RoleID      int
}

// ResFormData - store values to populate reservations form
type ResFormData struct {
	Client         Client           `json:"client" db:"client"`
	Venues         []Venues         `json:"venues" db:"venues"`
	VenueCount     int              `json:"venuecount" db:"venuecount"`
	Cities         []Cities         `json:"cities" db:"cities"`
	DepartureTimes []DepartureTimes `json:"departuretimes" db:"departuretimes"`
}

// Reservation - store reservation information
type Reservation struct {
	ReservationID            int `json:"reservationid" db:"reservationtid"`
	ClientID                 int `json:"clientid" db:"clientid"`
	DepartureCityID          int `json:"departurecityid" db:"departurecityid"`
	DepartureVenueID         int `json:"departurevenueid" db:"departurevenueid"`
	DepartureTimeID          int `json:"departuretimeid" db:"departuretimeid"`
	DestinationCityID        int `json:"destinationcityid" db:"destinationcityid"`
	DestinationVenueID       int `json:"destinationvenueid" db:"destinationvenueid"`
	ReturnDepartureCityID    int `json:"returndeparturecityid" db:"returndeparturecityid"`
	ReturnDepartureVenueID   int `json:"returndeparturevenueid" db:"returndeparturevenueid"`
	ReturnDepartureTimeID    int `json:"returndeparturetimeid" db:"returndeparturetimeid"`
	ReturnDestinationCityID  int `json:"returndestinationcityid" db:"returndestinationcityid"`
	ReturnDestinationVenueID int `json:"returndestinationvenueid" db:"returndestinationvenueid"`
	DiscountCodeID           int       `json:"discountcodeid" db:"discountcodeid"`
	//DiscountCode         DiscountCode
	DepartureAirlineID   int       `json:"departureairlineid" db:"departureairlineid"`
	ReturnAirlineID      int       `json:"returnairlineid" db:"returnairlineid"`
	DriverNotes          string    `json:"drivernotes" db:"drivernotes"`
	InternalNotes        string    `json:"internalnotes" db:"internalnotes"`
	DepartureNumAdults   int       `json:"departurenumadults" db:"departurenumadults"`
	DepartureNumStudents int       `json:"departurenumstudents" db:"departurenumstudent"`
	DepartureNumSeniors  int       `json:"departurenumseniors" db:"departurenumseniors"`
	DepartureNumChildren int       `json:"departurenumchildren" db:"departurenumchildren"`
	ReturnNumAdults      int       `json:"returnnumadults" db:"returnnumadults"`
	ReturnNumStudents    int       `json:"returnnumstudents" db:"returnnumstudent"`
	ReturnNumSeniors     int       `json:"returnnumseniors" db:"returnnumseniors"`
	ReturnNumChildren    int       `json:"returnnumchildren" db:"returnnumchildren"`
	Price                float32   `json:"price" db:"price"`
	Status               string    `json:"status" db:"status"`
	Hash                 string    `json:"hash" db:"hash"`
	CustomDepartureID    int       `json:"customdepartureid" db:"customdepartureid"`
	CustomDestinationID  int       `json:"customdestinationid" db:"customdestinationid"`
	DepartureDate        time.Time `json:"departuredate" db:"departuredate"`
	ReturnDate           time.Time `json:"returndate" db:"returndate"`
	TripTypeID           int       `json:"triptypeid" db:"triptypeid"`
	TripID               int       `json:"tripid" db:"tripid"`
	BalanceOwing         float32   `json:"balanceowing" db:"balanceowing"`
	ElavonTransactionID  int       `json:"elavontranscationid" db:"elavontransactionid"`
}

// DiscountCode - store discount code
type DiscountCode struct {
	DiscountCodeID int       `json:"discountcodeid" db:"discountcodeid"`
	Name           string    `json:"name" db:"name"`
	Percentage     int       `json:"percentage" db:"percentage"`
	Amount         int       `json:"amount" db:"amount"`
	StartDate      time.Time `json:"startdate" db:"startdate"`
	EndDate        time.Time `json:"enddate" db:"enddate"`
	//Whether percentage or amount, 1 = percentage, 2 = amount
	Type int
}

// DepartureTimes - store departure times
type DepartureTimes struct {
	DepartureTimeID int       `json:"departuretimeid" db:"departuretimeid"`
	CityID          int       `json:"cityid" db:"cityid"`
	DepartureTime   int       `json:"departuretime" db:"departuretime"`
	Recurring       int       `json:"recurring" db:"recurring"`
	StartDate       time.Time `json:"startdate" db:"startdate"`
	EndDate         time.Time `json:"enddate" db:"enddate"`
	RoleID          int
}

// Trips - store trip data
type Trips struct {
	TripID          int       `json:"tripid" db:"tripid"`
	DepartureDate   time.Time `json:"departuredate" db:"departuredate"`
	DepartureTimeID int       `json:"departuretimeid" db:"departuretimeid"`
	DepartureTime   int
	NumPassengers   int `json:"numpassengers" db:"numpassengers"`
	DriverID        int `json:"driverid" db:"driverid"`
	DriverList      []Drivers
	VehicleID       int `json:"vehicleid" db:"vehicleid"`
	VehicleList     []Vehicles
	Capacity        int       `json:"capacity" db:"capacity"`
	OmitTrip        bool      `json:"omittrip" db:"omittrip"`
	Postpone        bool      `json:"postpone" db:"postpone"`
	RescheduleDate  time.Time `json:"rescheduledate" db:"rescheduledate"`
	RescheduleTime  time.Time `json:"rescheduletime" db:"rescheduletime"`
	Cancelled       bool      `json:"cancelled" db:"cancelled"`
	RoleID          int
}

//Drivers - store driver data
type Drivers struct {
	DriverID  int    `json:"driverid" db:"driverid"`
	FirstName string `json:"firstname" db:"firstname"`
	LastName  string `json:"lastname" db:"lastname"`
	RoleID    int
}

//Vehicles - store vehicle data
type Vehicles struct {
	VehicleID    int    `json:"vehicleid" db:"vehicleid"`
	LicensePlate string `json:"licenseplate" db:"licenseplate"`
	NumSeats     int    `json:"numseats" db:"numseats"`
	Make         string `json:"make" db:"make"`
	RoleID       int
}

var dateLayout string
var sessionStore *sessions.CookieStore

// tpl holds all parsed templates
var tpl *template.Template

func newRouter() *mux.Router {
	r := mux.NewRouter()
	// Declare the static file directory and point it to the directory
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

	//map urls to handler functions any method
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/signup", SignupHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/reservation", ReservationHandler)
	r.HandleFunc("/createreservation", CreateReservationHandler)
	r.HandleFunc("/reservationcreated", ReservationCreatedHandler)
	r.HandleFunc("/createuser", CreateUserHandler)
	r.HandleFunc("/reports", ReportHandler)
	//get method only
	r.HandleFunc("/trips", TripHandler).Methods("GET")
	r.HandleFunc("/vehicles", VehicleHandler).Methods("GET")
	r.HandleFunc("/venues", VenueHandler).Methods("GET")
	r.HandleFunc("/drivers", DriverHandler).Methods("GET")
	r.HandleFunc("/cities", CityHandler).Methods("GET")
	r.HandleFunc("/price", PriceHandler).Methods("GET")
	//post method only
	r.HandleFunc("/signin", SigninHandler).Methods("POST")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	//r.HandleFunc("/trips", UpdateTripHandler).Methods("POST")
	r.HandleFunc("/vehicles", AddVehicleHandler).Methods("POST")
	r.HandleFunc("/venues", AddVenueHandler).Methods("POST")
	r.HandleFunc("/drivers", AddDriverHandler).Methods("POST")
	r.HandleFunc("/cities", AddCityHandler).Methods("POST")
	//put method only
	r.HandleFunc("/trips", UpdateTripHandler).Methods("PUT")
	r.HandleFunc("/cities", UpdateCityHandler).Methods("PUT")
	r.HandleFunc("/venues", UpdateVenueHandler).Methods("PUT")
	r.HandleFunc("/vehicles", UpdateVehicleHandler).Methods("PUT")
	r.HandleFunc("/drivers", UpdateDriverHandler).Methods("PUT")
	//delete method only
	r.HandleFunc("/cities", DeleteCityHandler).Methods("DELETE")
	r.HandleFunc("/venues", DeleteVenueHandler).Methods("DELETE")
	r.HandleFunc("/vehicles", DeleteVehicleHandler).Methods("DELETE")
	r.HandleFunc("/drivers", DeleteDriverHandler).Methods("DELETE")

	return r
}

func main() {
	r := newRouter()

	connString := "root:test@tcp(db:3306)/northernairport"
	db, err := sql.Open("mysql", connString)

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	InitSession()

	http.ListenAndServe(":8080", r)
}

// InitSession initializes a user session
func InitSession() {
	//generate random cookie authentication and encryption keys
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	//initialize new session store
	sessionStore = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	//must update, currently setting maxage to 2 minutes and http only (will need secure when live)
	sessionStore.Options = &sessions.Options{
		MaxAge:   60 * 2,
		HttpOnly: true,
	}

	//register client so it is stored in session
	gob.Register(Client{})

	//map template directory
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
