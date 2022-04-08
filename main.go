package main

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"github.com/gorilla/mux"

	"github.com/gorilla/context"

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
}

// VenueWrapper - for templates
type VenueWrapper struct {
	RoleID int
	Venues []Venues
}

// Cities - stores cities
type Cities struct {
	CityID      int    `json:"cityid" db:"cityid"`
	CityName    string `json:"cityname" db:"cityname"`
	NorthOffset int    `json:"northoffset" db:"northoffset"`
	SouthOffset int    `json:"southoffset" db:"southoffset"`
}

// CityWrapper - for templates
type CityWrapper struct {
	RoleID int
	Cities []Cities
}

// ResFormData - store values to populate reservations form
type ResFormData struct {
	Client         Client           `json:"client" db:"client"`
	Venues         []Venues         `json:"venues" db:"venues"`
	VenueCount     int              `json:"venuecount" db:"venuecount"`
	Cities         []Cities         `json:"cities" db:"cities"`
	DepartureTimes []DepartureTimes `json:"departuretimes" db:"departuretimes"`
	Airlines       []Airlines       `json:"airlines" db:"airlines"`
	AirlineCount   int              `json:"airlinecount" db:"airlinecount"`
}

// Reservation - store reservation information
type Reservation struct {
	ClientDetails            Client    `json:"clientdetails" db:"client"`
	ReservationID            int       `json:"reservationid" db:"reservationid"`
	ClientID                 int       `json:"clientid" db:"clientid"`
	DepartureCityID          int       `json:"departurecityid" db:"departurecityid"`
	DepartureVenueID         int       `json:"departurevenueid" db:"departurevenueid"`
	DepartureTimeID          int       `json:"departuretimeid" db:"departuretimeid"`
	DestinationCityID        int       `json:"destinationcityid" db:"destinationcityid"`
	DestinationVenueID       int       `json:"destinationvenueid" db:"destinationvenueid"`
	ReturnDepartureCityID    int       `json:"returndeparturecityid" db:"returndeparturecityid"`
	ReturnDepartureVenueID   int       `json:"returndeparturevenueid" db:"returndeparturevenueid"`
	ReturnDepartureTimeID    int       `json:"returndeparturetimeid" db:"returndeparturetimeid"`
	ReturnDestinationCityID  int       `json:"returndestinationcityid" db:"returndestinationcityid"`
	ReturnDestinationVenueID int       `json:"returndestinationvenueid" db:"returndestinationvenueid"`
	DiscountCodeID           int       `json:"discountcodeid" db:"discountcodeid"`
	DepartureAirlineID       int       `json:"departureairlineid" db:"departureairlineid"`
	ReturnAirlineID          int       `json:"returnairlineid" db:"returnairlineid"`
	DriverNotes              string    `json:"drivernotes" db:"drivernotes"`
	InternalNotes            string    `json:"internalnotes" db:"internalnotes"`
	DepartureNumAdults       int       `json:"departurenumadults" db:"departurenumadults"`
	DepartureNumStudents     int       `json:"departurenumstudents" db:"departurenumstudent"`
	DepartureNumSeniors      int       `json:"departurenumseniors" db:"departurenumseniors"`
	DepartureNumChildren     int       `json:"departurenumchildren" db:"departurenumchildren"`
	ReturnNumAdults          int       `json:"returnnumadults" db:"returnnumadults"`
	ReturnNumStudents        int       `json:"returnnumstudents" db:"returnnumstudent"`
	ReturnNumSeniors         int       `json:"returnnumseniors" db:"returnnumseniors"`
	ReturnNumChildren        int       `json:"returnnumchildren" db:"returnnumchildren"`
	Price                    float32   `json:"price" db:"price"`
	Status                   string    `json:"status" db:"status"`
	Hash                     string    `json:"hash" db:"hash"`
	CustomDepartureID        int       `json:"customdepartureid" db:"customdepartureid"`
	CustomDestinationID      int       `json:"customdestinationid" db:"customdestinationid"`
	DepartureDate            time.Time `json:"departuredate" db:"departuredate"`
	ReturnDate               time.Time `json:"returndate" db:"returndate"`
	TripTypeID               int       `json:"triptypeid" db:"triptypeid"`
	TripID                   int       `json:"tripid" db:"tripid"`
	ReturnTripID             int       `json:"returntripid" db:"returntripid"`
	BalanceOwing             float32   `json:"balanceowing" db:"balanceowing"`
	ElavonTransactionID      int       `json:"elavontranscationid" db:"elavontransactionid"`
	FlightNumber             int       `json:"flightnumber" db:"flightnumber"`
	FlightTime               int       `json:"flighttime" db:"flighttime"`
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
	DepartureTimeID int `json:"departuretimeid" db:"departuretimeid"`
	CityID          int `json:"cityid" db:"cityid"`
	CityList        []Cities
	DepartureTime   int       `json:"departuretime" db:"departuretime"`
	Recurring       int       `json:"recurring" db:"recurring"`
	StartDate       time.Time `json:"startdate" db:"startdate"`
	EndDate         time.Time `json:"enddate" db:"enddate"`
	Epoch           time.Time
}

// DepartureTimeWrapper - for templates
type DepartureTimeWrapper struct {
	RoleID         int
	DepartureTimes []DepartureTimes
}

// Trips - store trip data
type Trips struct {
	TripID          int       `json:"tripid" db:"tripid"`
	DepartureDate   time.Time `json:"departuredate" db:"departuredate"`
	DepartureTimeID int       `json:"departuretimeid" db:"departuretimeid"`
	DepartureTime   int
	NumPassengers   int    `json:"numpassengers" db:"numpassengers"`
	DriverID        int    `json:"driverid" db:"driverid"`
	DriverName      string `json:"drivername" db:"drivername"`
	DriverList      []Drivers
	VehicleID       int    `json:"vehicleid" db:"vehicleid"`
	LicensePlate    string `json:"licenseplate" db:"licenseplate"`
	VehicleList     []Vehicles
	Capacity        int       `json:"capacity" db:"capacity"`
	Omitted         int       `json:"omitted" db:"omitted"`
	RescheduleDate  time.Time `json:"rescheduledate" db:"rescheduledate"`
	RescheduleTime  time.Time `json:"rescheduletime" db:"rescheduletime"`
}

// TripWrapper - for templates
type TripWrapper struct {
	RoleID int
	Trips  []Trips
}

// OmitTrip - store omit trip data
type OmitTrip struct {
	TripID          int              `json:"tripid" db:"tripid"`
	DepartureDate   time.Time        `json:"departuredate" db:"departuredate"`
	DepartureTimeID int              `json:"departuretimeid" db:"departuretimeid"`
	DepartureTimes  []DepartureTimes `json:"departuretimes" db:"departuretimes"`
	RoleID          int
}

//SearchReservations - reservation display structure
type SearchReservations struct {
	ReservationID    int       `json:"reservationid" db:"reservationid"`
	ClientName       string    `json:"clientname" db:"clientname"`
	Phone            int       `json:"phone" db:"phone"`
	Email            string    `json:"email" db:"email"`
	DepartureVenue   string    `json:"departurevenue" db:"departurevenue"`
	DestinationVenue string    `json:"destinationvenue" db:"destinationvenue"`
	Return           int       `json:"return" db:"return"`
	NumPassengers    int       `json:"numpassengers" db:"numpassengers"`
	DepartureDate    time.Time `json:"departuredate" db:"departuredate"`
	DepartureTime    int       `json:"departuretime" db:"departuretime"`
	Postponed        bool      `json:"postponed" db:"postponed"`
	Cancelled        bool      `json:"cancelled" db:"cancelled"`
}

// SearchReservationWrapper - for templates
type SearchReservationWrapper struct {
	RoleID             int
	SearchReservations []SearchReservations
}

//DriverReport - reservation display structure by driver
type DriverReport struct {
	ReservationID    int    `json:"reservationid" db:"reservationid"`
	DriverID         int    `json:"driverid" db:"driverid"`
	DepartureCity    string `json:"departurecity" db:"departurecity"`
	DepartureVenue   string `json:"departurevenue" db:"departurevenue"`
	DestinationCity  string `json:"destinationcity" db:"destinationcity"`
	DestinationVenue string `json:"destinationvenue" db:"destinationvenue"`
	DepartureTime    int    `json:"departuretime" db:"departuretime"`
	ClientName       string `json:"passengername" db:"passengername"`
	NumPassengers    int    `json:"numpassengers" db:"numpassengers"`
}

// Drivers - store driver data
type Drivers struct {
	DriverID   int    `json:"driverid" db:"driverid"`
	FirstName  string `json:"firstname" db:"firstname"`
	LastName   string `json:"lastname" db:"lastname"`
	DriverName string `json:"drivername" db:"drivername"`
}

// DriverWrapper - for templates
type DriverWrapper struct {
	RoleID   int
	DriverID int
	Drivers  []Drivers
}

// DriverReportForm - driver and reservation display structure
type DriverReportForm struct {
	Drivers       []Drivers
	DriverReports []DriverReport
	RoleID        int
}

// Vehicles - store vehicle data
type Vehicles struct {
	VehicleID    int    `json:"vehicleid" db:"vehicleid"`
	LicensePlate string `json:"licenseplate" db:"licenseplate"`
	NumSeats     int    `json:"numseats" db:"numseats"`
	Make         string `json:"make" db:"make"`
}

// VehicleWrapper - for templates
type VehicleWrapper struct {
	RoleID   int
	Vehicles []Vehicles
}

//Airlines - store airline data
type Airlines struct {
	AirlineID int    `json:"airlineid" db:"airlineid"`
	Name      string `json:"name" db:"name"`
	Terminal  int    `json:"terminal" db:"terminal"`
}

// Prices - store price data
type Prices struct {
	PriceID           int     `json:"priceid" db:"priceid"`
	CustomerTypeID    int     `json:"customertypeid" db:"customertypeid"`
	CustomerType      string  `json:"customertype" db:"customertype"`
	Price             float32 `json:"price" db:"price"`
	ReturnPrice       float32 `json:"returnprice" db:"returnprice"`
	DepartureCity     string  `json:"departurecity" db:"departurecity"`
	DepartureCityID   int     `json:"departurecityid" db:"departurecityid"`
	DestinationCity   string  `json:"destinationcity" db:"destinationcity"`
	DestinationCityID int     `json:"destinationcityid" db:"destinationcityid"`
}

// PriceWrapper - for templates
type PriceWrapper struct {
	RoleID int
	Prices []Prices
}

//TravelAgencies - store travel agency data
type TravelAgencies struct {
	TravelAgencyID   int    `json:"travelagencyid" db:"travelagencyid"`
	TravelAgencyName string `json:"travelagencyname" db:"travelagencyname"`
	IATANumber       string `json:"iatanumber" db:"iatanumber"`
}

//TravelAgencyReport - store travel report data
type TravelAgencyReport struct {
	TravelAgencyID   int     `json:"travelagencyid" db:"travelagencyid"`
	TravelAgencyName string  `json:"travelagencyname" db:"travelagencyname"`
	ReservationCount int     `json:"reservationcount" db:"reservationcount"`
	TotalCost        float32 `json:"totalcount" db:"totalcount"`
	Commission       float32 `json:"commission" db:"commission"`
	ReportMonth      int     `json:"reportmonth" db:"reportmonth"`
	ReportYear       int     `json:"reportyear" db:"reportyear"`
}

//TravelAgencyForm - store travel agency data for report form
type TravelAgencyForm struct {
	TravelAgencyReports []TravelAgencyReport
	RoleID              int
}

//CalendarDays - store days and trips for calendar report
type CalendarDays struct {
	DayNum        int
	CalendarTrips []Trips
}

//CalendarReport - store calendar report info
type CalendarReport struct {
	CurrentDate time.Time
	Days        []CalendarDays
	RoleID      int
}

//AGTAReport - store values for AGTA spreadsheet
type AGTAReport struct {
	ReservationID      int
	FlightTime         string
	AirlineName        string
	FlightNumber       string
	FlightCity         string
	TerminalName       string
	ConfirmationNumber int
	PaxName            string
	NumPax             int
	DropLocation       string
	DropCity           string
	InternalNotes      string
	DriverNotes        string
	DepartureTime      int
	DriverName         string
	DriverID           int
	VehicleNum         string
	IsValid            int
	DepartureDate      time.Time
	HotelInfo          string
	Cancelled          int
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
	//get method only
	r.HandleFunc("/trips", TripHandler).Methods("GET")
	r.HandleFunc("/omittrip", OmitTripFormHandler).Methods("GET")
	r.HandleFunc("/vehicles", VehicleHandler).Methods("GET")
	r.HandleFunc("/venues", VenueHandler).Methods("GET")
	r.HandleFunc("/drivers", DriverHandler).Methods("GET")
	r.HandleFunc("/cities", CityHandler).Methods("GET")
	r.HandleFunc("/price", GetPriceHandler).Methods("GET")
	r.HandleFunc("/prices", PriceHandler).Methods("GET")
	r.HandleFunc("/times", DepartureTimeHandler).Methods("GET")
	r.HandleFunc("/badsignin", BadSignInHandler).Methods("GET")
	r.HandleFunc("/search", SearchHandler).Methods("GET")
	r.HandleFunc("/driverreport", DriverReportHandler).Methods("GET")
	r.HandleFunc("/travelagencyreport", TravelAgencyReportHandler).Methods("GET")
	r.HandleFunc("/calendarreport", CalendarReportHandler).Methods("GET")
	r.HandleFunc("/import", ImportHandler).Methods("GET")
	r.HandleFunc("/elavon", ElavonHandler).Methods("GET")
	r.HandleFunc("/migrate", MigrateHandler).Methods("GET")
	//post method only
	r.HandleFunc("/signin", SigninHandler).Methods("POST")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/vehicles", AddVehicleHandler).Methods("POST")
	r.HandleFunc("/venues", AddVenueHandler).Methods("POST")
	r.HandleFunc("/drivers", AddDriverHandler).Methods("POST")
	r.HandleFunc("/cities", AddCityHandler).Methods("POST")
	r.HandleFunc("/times", AddDepartureTimeHandler).Methods("POST")
	//r.HandleFunc("/elavon", ElavonHandler).Methods("POST")
	r.HandleFunc("/transactionstatus", TransactionStatusHandler).Methods("POST")
	//put method only
	r.HandleFunc("/trips", UpdateTripHandler).Methods("PUT")
	r.HandleFunc("/cities", UpdateCityHandler).Methods("PUT")
	r.HandleFunc("/prices", UpdatePriceHandler).Methods("PUT")
	r.HandleFunc("/venues", UpdateVenueHandler).Methods("PUT")
	r.HandleFunc("/vehicles", UpdateVehicleHandler).Methods("PUT")
	r.HandleFunc("/drivers", UpdateDriverHandler).Methods("PUT")
	r.HandleFunc("/times", UpdateDepartureTimeHandler).Methods("PUT")
	r.HandleFunc("/postpone", PostponeHandler).Methods("PUT")
	r.HandleFunc("/cancel", CancelHandler).Methods("PUT")
	r.HandleFunc("/omittrip", OmitTripHandler).Methods("PUT")
	//delete method only
	r.HandleFunc("/cities", DeleteCityHandler).Methods("DELETE")
	r.HandleFunc("/venues", DeleteVenueHandler).Methods("DELETE")
	r.HandleFunc("/vehicles", DeleteVehicleHandler).Methods("DELETE")
	r.HandleFunc("/drivers", DeleteDriverHandler).Methods("DELETE")

	return r
}

func main() {
	r := newRouter()
	http.Handle("/", r)
	r.Use(SetCSPHeaders)

	connString := "root:ah83is82js95pq@tcp(db:3306)/northernairport"
	//connString := "root:@tcp(db:3306)/northernairport"
	db, err := sql.Open("mysql", connString)

	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})

	InitSession()

	//http.ListenAndServe(":80", r)
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

func SetCSPHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  w.Header().Set("Content-Type", "application/json; charset=utf-8")
		//w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'self';style-src 'self' img-src 'self'; media-src 'self'; frame-ancestors 'self'; frame-src 'self'; connect-src 'self'")

		header := w.Header()
		csp := []string{"default-src 'self' ",
			"style-src 'self' https://cdnjs.cloudflare.com",
			"font-src 'self' https://cdnjs.cloudflare.com",
			//"script-src 'self' https://cdnjs.cloudflare.com 'sha256-Pj4adVzOLitxyl2o99sCaYR4mXGXeF6whHx5TnSPjZo=' 'unsafe-hashes'",
			"script-src 'self' https://cdnjs.cloudflare.com https://cdn.mxpnl.com https://api.demo.convergepay.com https://cdn.appdynamics.com 'unsafe-inline'",
			"frame-src 'none'"}

		//header.Set("Content-Type", "text/html; charset=UTF-8")
		header.Set("Content-Security-Policy", strings.Join(csp, "; "))

		next.ServeHTTP(w, r)
	})
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
		MaxAge:   60 * 60,
		HttpOnly: true,
	}

	//register client so it is stored in session
	gob.Register(Client{})

	//map template directory
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
