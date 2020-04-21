package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//RegisterHandler - parse data from client registration from and store in db
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	client := &Client{}

	//get form data
	err := r.ParseForm()

	if err != nil {
		log.Fatal("Parse Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//store form data in client data structure
	client.Username = r.Form.Get("username")
	client.Password = r.Form.Get("password")
	client.Firstname = r.Form.Get("firstname")
	client.Lastname = r.Form.Get("lastname")
	client.Phone = r.Form.Get("phone")
	client.Email = r.Form.Get("email")
	client.StreetAddress = r.Form.Get("streetaddress")
	client.City = r.Form.Get("city")
	client.Province = r.Form.Get("provstate")
	client.PostalCode = r.Form.Get("postalzip")
	client.Country = r.Form.Get("country")

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(client.Password), 8)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Fatal("Hashing Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client.Password = string(hashedPassword)

	err = store.CreateUser(client)
	if err != nil {
		log.Fatal("Error Creating User: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print("user created")

	http.Redirect(w, r, "/", http.StatusFound)
}

//SigninHandler - compare credentials entered vs credentials in db, if match found store
//				  the username in client structure and create authentication cookie to
//				  allow for authentication persistence site wide
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//create a client object
	client := &Client{}

	//get credentials from form
	err = r.ParseForm()

	if err != nil {
		log.Fatal("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client.Username = r.Form.Get("username")
	client.Password = r.Form.Get("password")

	//validate user credentials
	err = store.SignInUser(client)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Error: ", err)
			//w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/badsignin", http.StatusFound)
		} else if err.Error() == "Incorrect Password" {
			log.Print("Error: ", err)
			//w.WriteHeader(http.StatusUnauthorized)
			http.Redirect(w, r, "/badsignin", http.StatusFound)
		} else {
			//this should go somewhere else, not just a bad sign in attempt if this happens
			log.Print("Error: ", err)
			//w.WriteHeader(http.StatusInternalServerError)
			http.Redirect(w, r, "/badsignin", http.StatusFound)
		}
	}

	client.Authenticated = true

	//update authentication cookie
	session.Values["client"] = client

	log.Printf("creating session for %s", client.Username)

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s is signed in", session.Values["client.Username"])

	http.Redirect(w, r, "/reservation", http.StatusFound)
}

//LogoutHandler - set authentication cookie to be expired
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//set authentication cookie to be expired
	session.Values["client"] = Client{}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//BadSignInHandler - display homepage
func BadSignInHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "badsignin.gohtml", nil)
}

//IndexHandler - display homepage
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	log.Printf("username: %s", client.Username)

	tpl.ExecuteTemplate(w, "index.gohtml", client)
}

//SignupHandler - display signup page
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	tpl.ExecuteTemplate(w, "signup.gohtml", client)
}

//ReservationHandler - display reservation page, populate client data if user is authenticated
func ReservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated {
		err := store.GetClientInfo(&client)

		if err != nil {
			log.Panicf("Error getting client info: %s", err.Error())
		}
	}

	//get data need to populate dropdowns in reservation form
	venues := store.GetVenues()
	cities := store.GetCities()
	departuretimes := store.GetDepartureTimes()
	airlines := store.GetAirlines()

	//store retrieved data in reservation structure and pass to template
	resform := ResFormData{}

	resform.Client = client
	resform.Venues = venues
	resform.VenueCount = len(venues)
	resform.Cities = cities
	resform.DepartureTimes = departuretimes
	resform.Airlines = airlines

	tpl.ExecuteTemplate(w, "reservation.gohtml", resform)
}

//CreateReservationHandler - store reservation in database
func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	reservation := Reservation{}

	log.Printf("Retrieve form data")

	//get reservation data from form
	reservation = GetReservationFormValues(r, true)

	log.Printf("Form information retrieved")

	//check if trip exists, if not create one
	err = store.GetOrAddTrip(&reservation)

	if err != nil {
		log.Printf("Error creating reservation: %s", err.Error())
		//http.Redirect(w, r, "/reservationnotcreated", http.StatusFound)
		tpl.ExecuteTemplate(w, "reservationnotcreated.gohtml", err.Error())
	} else {
		log.Printf("Trip retrieved/added")
		err = store.CreateReservation(&reservation)

		if err != nil {
			log.Fatal("Error creating reservation: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Print("reservation created")
		}

		//http.Redirect(w, r, "/reservationcreated", http.StatusFound)
		tpl.ExecuteTemplate(w, "created.gohtml", client)
	}

}

//ReservationCreatedHandler - redirect to created page
func ReservationCreatedHandler(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "created.gohtml", r)
}

//DriverHandler - display driver admin page
func DriverHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		driverwrapper := DriverWrapper{}
		driverwrapper.Drivers = store.GetDrivers()

		driverwrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "driver.gohtml", driverwrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//AddDriverHandler - add driver to database
func AddDriverHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {

		err := r.ParseForm()

		driver := Drivers{}

		driver.FirstName = r.FormValue("firstname")
		driver.LastName = r.FormValue("lastname")

		err = store.AddDriver(driver)

		if err != nil {
			log.Printf("Error adding driver: %s", err.Error())
		} else {
			log.Print("Drivers added")
		}

		driverwrapper := DriverWrapper{}
		driverwrapper.Drivers = store.GetDrivers()

		driverwrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "driver.gohtml", driverwrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateDriverHandler - update driver in database
func UpdateDriverHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("execute update driver handler")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		driver := Drivers{}
		var err error

		driver.DriverID, err = strconv.Atoi(values["driverid"][0])

		if err != nil {
			log.Printf("Error converting driverid: %s", err.Error())
		}

		driver.FirstName = values["firstname"][0]
		driver.LastName = values["lastname"][0]

		store.UpdateDriver(&driver)

		driverwrapper := DriverWrapper{}
		driverwrapper.Drivers[0] = driver

		driverwrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "driver.gohtml", driverwrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//DeleteDriverHandler - delete driver from database
func DeleteDriverHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		driverid, err := strconv.Atoi(values["driverid"][0])

		store.DeleteDriver(driverid)

		driverwrapper := DriverWrapper{}
		driverwrapper.Drivers = store.GetDrivers()

		driverwrapper.RoleID = client.RoleID

		if err != nil {
			log.Printf("Error converting driverid: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "driver.gohtml", driverwrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//VehicleHandler - display van admin page
func VehicleHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		//get data need to populate dropdowns in reservation form
		vehiclewrapper := VehicleWrapper{}
		vehiclewrapper.Vehicles = store.GetVehicles()

		vehiclewrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehiclewrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//AddVehicleHandler - add vehicle to database
func AddVehicleHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {

		err := r.ParseForm()

		vehicle := Vehicles{}

		vehicle.LicensePlate = r.FormValue("license-plate")
		vehicle.NumSeats, err = strconv.Atoi(r.FormValue("num-seats"))
		vehicle.Make = r.FormValue("make")

		err = store.AddVehicle(vehicle)

		if err != nil {
			log.Printf("Error adding vehicle: %s", err.Error())
		} else {
			log.Print("Vehicles added")
		}

		vehiclewrapper := VehicleWrapper{}
		vehiclewrapper.Vehicles = store.GetVehicles()

		vehiclewrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehiclewrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateVehicleHandler - update vehicle in database
func UpdateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("execute update vehicle handler")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		vehicle := Vehicles{}
		var err error

		vehicle.VehicleID, err = strconv.Atoi(values["vehicleid"][0])

		if err != nil {
			log.Printf("Error converting vehicleid: %s", err.Error())
		}

		vehicle.LicensePlate = values["licenseplate"][0]
		vehicle.NumSeats, err = strconv.Atoi(values["numseats"][0])

		if err != nil {
			log.Printf("Error converting numseats: %s", err.Error())
		}

		vehicle.Make = values["make"][0]

		store.UpdateVehicle(&vehicle)

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehicle)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//DeleteVehicleHandler - delete vehicle from database
func DeleteVehicleHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		vehicleid, err := strconv.Atoi(values["vehicleid"][0])

		store.DeleteVehicle(vehicleid)

		//get data need to populate dropdowns in reservation form
		vehicle := store.GetVehicles()

		if err != nil {
			log.Printf("Error converting vehicleid: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehicle)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//CreateUserHandler - display create user page
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "createuser.gohtml", r)
}

//TripHandler - display trip admin page
func TripHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	var tripdate time.Time

	if values["searchdate"] != nil {
		tripdate, err = time.Parse("2006-01-02", values["searchdate"][0])
	}

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		searchtripwrapper := TripWrapper{}
		searchtripwrapper.Trips = store.SearchTrips(tripdate, 0)

		if len(searchtripwrapper.Trips) > 0 {
			searchtripwrapper.RoleID = client.RoleID

			log.Printf("trips: we've got %d trips!", len(searchtripwrapper.Trips))
		} else {
			log.Printf("trips: no trips returned!")
		}

		if err := tpl.ExecuteTemplate(w, "trip.gohtml", searchtripwrapper); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateTripHandler - update trip
func UpdateTripHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		values := r.URL.Query()

		var err error
		trip := Trips{}

		trip.TripID, err = strconv.Atoi(values["tripid"][0])

		if err != nil {
			log.Printf("Error converting tripid: %s", err.Error())
		}

		trip.DriverID, err = strconv.Atoi(values["driverid"][0])

		if err != nil {
			log.Printf("Error converting driverid: %s", err.Error())
		}

		trip.VehicleID, err = strconv.Atoi(values["vehicleid"][0])

		if err != nil {
			log.Printf("Error converting vehicleid: %s", err.Error())
		}

		store.UpdateTrip(&trip)

		tripwrapper := TripWrapper{}
		tripwrapper.Trips = store.GetTrips()

		tripwrapper.RoleID = client.RoleID

		//tpl.ExecuteTemplate(w, "trip.gohtml", trips)
		http.Redirect(w, r, "/trips", http.StatusFound)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//OmitTripFormHandler - omit trip
func OmitTripFormHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		omittrip := OmitTrip{}

		omittrip.DepartureTimes = store.GetDepartureTimes()

		omittrip.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "omittrip.gohtml", omittrip)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//OmitTripHandler - omit trip
func OmitTripHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("omitting trip")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("got session")

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		trip := Trips{}

		log.Printf("populate trip data")

		trip.DepartureDate, err = time.Parse("2006-01-02", values["departuredate"][0])

		trip.DepartureTimeID, err = strconv.Atoi(values["departuretimeid"][0])

		log.Printf("call omit function")

		err := store.OmitTrip(&trip)

		if err != nil {
			log.Printf("Error omitting trip: %s", err.Error())
		} else {
			log.Printf("Trip has been omitted!")
		}

		omittrip := OmitTrip{}

		omittrip.DepartureTimes = store.GetDepartureTimes()

		tpl.ExecuteTemplate(w, "omittrip.gohtml", omittrip)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//SearchHandler - display reservations
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var name string
	var phone int
	var email string

	if values["searchname"] != nil {
		name = values["searchname"][0]
	} else {
		name = ""
	}

	if values["searchphone"] != nil {
		phone, err = strconv.Atoi(values["searchphone"][0])
	} else {
		phone = 0
	}

	if values["searchemail"] != nil {
		email = values["searchemail"][0]
	} else {
		email = ""
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		searchreservationwrapper := SearchReservationWrapper{}

		//get data need to populate dropdowns in reservation form
		searchreservationwrapper.SearchReservations = store.SearchReservations(name, phone, email)

		if len(searchreservationwrapper.SearchReservations) > 0 {
			searchreservationwrapper.RoleID = client.RoleID

			log.Printf("search: we've got some reservations!")
		} else {
			log.Printf("search: no reservations returned!")
		}

		if err := tpl.ExecuteTemplate(w, "searchreservations.gohtml", searchreservationwrapper); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//PostponeHandler - postpone reservation
func PostponeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		values := r.URL.Query()

		var err error
		searchreservations := SearchReservations{}

		searchreservations.ReservationID, err = strconv.Atoi(values["reservationid"][0])

		if err != nil {
			log.Printf("Error converting reservationid: %s", err.Error())
		}

		store.PostponeReservation(&searchreservations)

		searchreservationwrapper := SearchReservationWrapper{}

		searchreservationwrapper.SearchReservations = store.SearchReservations("", 0, "")

		searchreservationwrapper.RoleID = client.RoleID

		for i, res := range searchreservationwrapper.SearchReservations {
			log.Printf("#%d ReservationID Returned: %d - Postponed: %t", i, res.ReservationID, res.Postponed)
		}

		if err := tpl.ExecuteTemplate(w, "searchreservations.gohtml", searchreservationwrapper); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//CancelHandler - cancel reservation
func CancelHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		values := r.URL.Query()

		var err error
		searchreservations := SearchReservations{}

		searchreservations.ReservationID, err = strconv.Atoi(values["reservationid"][0])

		if err != nil {
			log.Printf("Error converting reservationid: %s", err.Error())
		}

		store.CancelReservation(&searchreservations)

		searchreservationwrapper := SearchReservationWrapper{}

		searchreservationwrapper.SearchReservations = store.SearchReservations("", 0, "")

		searchreservationwrapper.RoleID = client.RoleID

		for i, res := range searchreservationwrapper.SearchReservations {
			log.Printf("#%d ReservationID Returned: %d - Cancelled: %t", i, res.ReservationID, res.Postponed)
		}

		var buf bytes.Buffer
		//tpl.ExecuteTemplate(&buf, "searchreservations.gohtml", searchreservations)
		if err := tpl.ExecuteTemplate(&buf, "searchreservations.gohtml", searchreservationwrapper); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		io.Copy(w, &buf)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//VenueHandler - display venue admin page
func VenueHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		venuewrapper := VenueWrapper{}

		venuewrapper.Venues = store.GetVenues()

		venuewrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "venue.gohtml", venuewrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//AddVenueHandler - add venue to database
func AddVenueHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {

		err := r.ParseForm()

		venue := Venues{}

		cityname := r.FormValue("cityname")

		venue.CityID = store.GetCityID(cityname)
		venue.VenueName = r.FormValue("venuename")
		venue.Active, err = strconv.Atoi(r.FormValue("active"))

		err = store.AddVenue(venue)

		if err != nil {
			log.Printf("Error adding venue: %s", err.Error())
		} else {
			log.Print("Venues added")
		}

		venuewrapper := VenueWrapper{}

		venuewrapper.Venues = store.GetVenues()

		venuewrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "venue.gohtml", venuewrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateVenueHandler - update city in1 database
func UpdateVenueHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("execute update venue handler")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		venue := Venues{}
		var err error

		venue.VenueID, err = strconv.Atoi(values["venueid"][0])

		if err != nil {
			log.Printf("Error converting venueid: %s", err.Error())
		}

		venue.VenueName = values["venuename"][0]
		venue.ExtraCost, err = strconv.Atoi(values["extracost"][0])

		if err != nil {
			log.Printf("Error converting extracost: %s", err.Error())
		}

		venue.Active, err = strconv.Atoi(values["active"][0])

		if err != nil {
			log.Printf("Error converting active: %s", err.Error())
		}

		venue.ExtraTime, err = strconv.Atoi(values["extratime"][0])

		if err != nil {
			log.Printf("Error converting extratime: %s", err.Error())
		}

		store.UpdateVenue(&venue)

		tpl.ExecuteTemplate(w, "venue.gohtml", venue)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//DeleteVenueHandler - delete venue from database
func DeleteVenueHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		venueid, err := strconv.Atoi(values["venueid"][0])

		store.DeleteVenue(venueid)

		//get data need to populate dropdowns in reservation form
		venues := store.GetVenues()

		if err != nil {
			log.Printf("Error converting northoffset or southoffset: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "venue.gohtml", venues)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//CityHandler - display city admin page
func CityHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		citywrapper := CityWrapper{}

		citywrapper.Cities = store.GetCities()

		citywrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "city.gohtml", citywrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//AddCityHandler - add city to database
func AddCityHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {

		err := r.ParseForm()

		city := Cities{}

		city.CityName = r.FormValue("cityname")
		city.NorthOffset, err = strconv.Atoi(r.FormValue("northoffset"))
		city.SouthOffset, err = strconv.Atoi(r.FormValue("southoffset"))

		err = store.AddCity(city)

		if err != nil {
			log.Printf("Error adding city: %s", err.Error())
		} else {
			log.Print("Cities added")
		}

		citywrapper := CityWrapper{}

		citywrapper.Cities = store.GetCities()

		citywrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "city.gohtml", citywrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateCityHandler - update city in database
func UpdateCityHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("execute update city handler")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		city := Cities{}
		var err error

		city.CityID, err = strconv.Atoi(values["cityid"][0])
		city.CityName = values["cityname"][0]
		city.NorthOffset, err = strconv.Atoi(values["northoffset"][0])
		city.SouthOffset, err = strconv.Atoi(values["southoffset"][0])

		store.UpdateCity(&city)

		if err != nil {
			log.Printf("Error converting northoffset or southoffset: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "city.gohtml", city)
	} else {
		tpl.ExecuteTemplate(w, "city.gohtml", r)
	}

}

//PriceHandler - display price admin page
func PriceHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		pricewrapper := PriceWrapper{}

		pricewrapper.Prices = store.GetPrices()

		pricewrapper.RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "prices.gohtml", pricewrapper)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdatePriceHandler - update price in database
func UpdatePriceHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("execute update price handler")

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		price := Prices{}
		var err error

		log.Printf("get values from form")

		price.PriceID, err = strconv.Atoi(values["priceid"][0])
		var tempPrice float64
		tempPrice, err = strconv.ParseFloat(values["price"][0], 32)
		price.Price = float32(tempPrice)

		log.Printf("saves form values")

		store.UpdatePrice(&price)

		log.Printf("updated ")

		pricewrapper := PriceWrapper{}

		pricewrapper.Prices = store.GetPrices()

		pricewrapper.RoleID = client.RoleID

		if err != nil {
			log.Printf("Error converting northoffset or southoffset: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "prices.gohtml", pricewrapper)
	} else {
		tpl.ExecuteTemplate(w, "city.gohtml", r)
	}

}

//DeleteCityHandler - delete city from database
func DeleteCityHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		cityid, err := strconv.Atoi(values["cityid"][0])

		store.DeleteCity(cityid)

		if err != nil {
			log.Printf("Error converting northoffset or southoffset: %s", err.Error())
		}

		//get data need to populate dropdowns in reservation form
		cities := store.GetCities()

		tpl.ExecuteTemplate(w, "city.gohtml", cities)
	} else {
		tpl.ExecuteTemplate(w, "city.gohtml", r)
	}

}

//GetPriceHandler - return price given the departure, destination and customer type
func GetPriceHandler(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()

	triptypeid, err := strconv.Atoi(values["triptypeid"][0])
	departurecityid, err := strconv.Atoi(values["departurecityid"][0])
	destinationcityid, err := strconv.Atoi(values["destinationcityid"][0])
	retdeparturecityid, err := strconv.Atoi(values["retdeparturecityid"][0])
	retdestinationcityid, err := strconv.Atoi(values["retdestinationcityid"][0])
	numpassengers, err := strconv.Atoi(values["numpassengers"][0])
	customertypeid, err := strconv.Atoi(values["customertypeid"][0])
	discountcode := values["discount"][0]

	price := store.GetPrice(departurecityid, destinationcityid, retdeparturecityid, retdestinationcityid, customertypeid, triptypeid, discountcode)

	departurevenueid, err := strconv.Atoi(values["departurevenueid"][0])
	destinationvenueid, err := strconv.Atoi(values["destinationvenueid"][0])
	retdeparturevenueid, err := strconv.Atoi(values["retdeparturevenueid"][0])
	retdestinationvenueid, err := strconv.Atoi(values["retdestinationvenueid"][0])

	extracost := store.AddVenueFee(departurevenueid, destinationvenueid, retdeparturevenueid, retdestinationvenueid)

	totalprice := price*float32(numpassengers) + extracost

	log.Printf("Total Price in handler: %f", totalprice)

	if err != nil {
		log.Printf("Error converting northoffset or southoffset: %s", err.Error())
	}

	//tpl.ExecuteTemplate(w, "city.gohtml", city)
	fmt.Fprintf(w, "%f", totalprice)
}

//DepartureTimeHandler - display time add/update/change page
func DepartureTimeHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		//get data need to populate dropdowns in reservation form
		times := store.GetDepartureTimes()

		//trips[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "departuretime.gohtml", times)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateDepartureTimeHandler - update departuretime in database
func UpdateDepartureTimeHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {
		values := r.URL.Query()

		departuretime := DepartureTimes{}
		var err error

		departuretime.DepartureTimeID, err = strconv.Atoi(values["departuretimeid"][0])
		departuretime.CityID, err = strconv.Atoi(values["cityid"][0])
		departuretime.DepartureTime, err = strconv.Atoi(values["departuretime"][0])
		departuretime.Recurring, err = strconv.Atoi(values["recurring"][0])

		if values["startdate"][0] == "" {
			departuretime.StartDate = time.Time{}

			log.Printf("StartDate: %s", departuretime.StartDate)
		} else {
			departuretime.StartDate, err = time.Parse("2006-01-02", values["startdate"][0])

			if err != nil {
				log.Printf("Error converting StartDate to time.Time: %s", err.Error())
			}
		}

		if values["enddate"][0] == "" {
			departuretime.EndDate = time.Time{}

			log.Printf("EndDate: %s", departuretime.EndDate)
		} else {
			departuretime.EndDate, err = time.Parse("2006-01-02", values["enddate"][0])

			if err != nil {
				log.Printf("Error converting EndDate to time.Time: %s", err.Error())
			}
		}

		store.UpdateDepartureTime(&departuretime)

		tpl.ExecuteTemplate(w, "departuretime.gohtml", departuretime)
	} else {
		tpl.ExecuteTemplate(w, "departuretime.gohtml", r)
	}

}

//AddDepartureTimeHandler - update departuretime in database
func AddDepartureTimeHandler(w http.ResponseWriter, r *http.Request) {

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 4) {

		values := r.URL.Query()

		departuretimewrapper := DepartureTimeWrapper{}
		departuretime := DepartureTimes{}
		var err error

		departuretimewrapper.RoleID = client.RoleID

		departuretime.CityID, err = strconv.Atoi(values["cityid"][0])
		departuretime.DepartureTime, err = strconv.Atoi(values["departuretime"][0])
		departuretime.Recurring, err = strconv.Atoi(values["recurring"][0])

		if values["startdate"][0] == "" {
			departuretime.StartDate = time.Time{}

			log.Printf("StartDate: %s", departuretime.StartDate)
		} else {
			departuretime.StartDate, err = time.Parse("2006-01-02", values["startdate"][0])

			if err != nil {
				log.Printf("Error converting StartDate to time.Time: %s", err.Error())
			}
		}

		if values["enddate"][0] == "" {
			departuretime.EndDate = time.Time{}

			log.Printf("EndDate: %s", departuretime.EndDate)
		} else {
			departuretime.EndDate, err = time.Parse("2006-01-02", values["enddate"][0])

			if err != nil {
				log.Printf("Error converting EndDate to time.Time: %s", err.Error())
			}
		}

		store.AddDepartureTime(&departuretime)

		departuretimewrapper.DepartureTimes[0] = departuretime

		tpl.ExecuteTemplate(w, "departuretime.gohtml", departuretimewrapper)
	} else {
		tpl.ExecuteTemplate(w, "departuretime.gohtml", r)
	}

}

//DriverReportHandler - display driver report page
func DriverReportHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("get the filter criteria")

	var driverid int
	var date time.Time

	if values["driverid"] != nil {
		driverid, err = strconv.Atoi(values["driverid"][0])

		if err != nil {
			log.Printf("Error getting driverid: %s", err.Error())
		}
	} else {
		driverid = 0
	}

	if values["departuredate"] != nil {
		date, err = time.Parse("2006-01-02", values["departuredate"][0])

		if err != nil {
			log.Printf("Error getting date: %s", err.Error())
		}
	} else {
		date = time.Time{}
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4 || client.RoleID == 5) {
		//get data need to populate dropdowns in reservation form
		log.Printf("get driver reservations")
		driverreservations := store.DriverReservations(driverid, date)
		log.Printf("reservations retrieved")

		if len(driverreservations.DriverReports) > 0 {
			log.Printf("driver report: we've got some reservations!")
		} else {
			log.Printf("driver report: no reservations returned!")
		}

		driverreservations.RoleID = client.RoleID

		if err := tpl.ExecuteTemplate(w, "driverreport.gohtml", driverreservations); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//TravelAgencyReportHandler - display travel agent report page
func TravelAgencyReportHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("get travel agent criteria")

	var month int
	var year int

	if values["month"] != nil {
		month, err = strconv.Atoi(values["month"][0])

		if err != nil {
			log.Printf("Error getting month: %s", err.Error())
		}
	} else {
		month = 0
	}

	if values["year"] != nil {
		year, err = strconv.Atoi(values["year"][0])

		if err != nil {
			log.Printf("Error getting year: %s", err.Error())
		}
	} else {
		year = 0
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		//get data need to populate dropdowns in reservation form
		log.Printf("get travel agent records")
		travelagentreports := store.TravelAgencyReports(month, year)
		log.Printf("travel agent reports retrieved")

		if len(travelagentreports) > 0 {
			log.Printf("travel agent report: we've got some records!")
			log.Printf("handler: travel agent report: %s", travelagentreports[0].TravelAgencyName)
		} else {
			log.Printf("travel agent report: no records returned!")
		}

		travelagencyform := TravelAgencyForm{}
		travelagencyform.TravelAgencyReports = travelagentreports
		travelagencyform.RoleID = client.RoleID

		if err := tpl.ExecuteTemplate(w, "travelagentreport.gohtml", travelagencyform); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//CalendarReportHandler - display reports admin page
func CalendarReportHandler(w http.ResponseWriter, r *http.Request) {
	//values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//get client data from session cookie
	client := GetClient(session)

	//if authenticated get all client info
	if client.Authenticated && (client.RoleID == 3 || client.RoleID == 4) {
		//get days in month
		currentTime := time.Now()
		month := currentTime.Month()
		year := currentTime.Year()

		numdays := NumDaysLookup(int(month), year)

		calReport := CalendarReport{}
		calDays := make([]CalendarDays, numdays)

		calReport.RoleID = client.RoleID

		for indx := 0; indx < numdays; indx++ {
			calDays[indx].DayNum = indx + 1

			tripday := time.Date(year, month, indx, 0, 0, 0, 0, time.UTC)
			tripsbyday := store.SearchTrips(tripday, 1)
			calDays[indx].CalendarTrips = tripsbyday
		}

		calReport.Days = calDays
		calReport.CurrentDate = currentTime

		tpl.ExecuteTemplate(w, "calendar.gohtml", calReport)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//ImportHandler - run import script
func ImportHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("create AGTA script")

	apikey := "jgoiwjerfgi8432u"
	parameterfailure := 0

	values := r.URL.Query()

	var key string
	var startdate time.Time
	var enddate time.Time
	var err error

	//check if api key received
	if values["key"] != nil {
		key = values["key"][0]
	} else {
		log.Fatal("No api key found")
		parameterfailure = 1
	}

	//check if valid start date
	if values["startDate"] != nil {
		startdate, err = time.Parse("2006-01-02", values["startDate"][0])

		if err != nil {
			log.Fatal("Problem parsing startdate")
		}
	} else {
		log.Fatal("Start date not found")
		parameterfailure = 1
	}

	//check if valid end date
	if values["endDate"] != nil {
		enddate, err = time.Parse("2006-01-02", values["endDate"][0])

		if err != nil {
			log.Fatal("Problem parsing enddate")
		}
	} else {
		log.Fatal("End date not found")
		parameterfailure = 1
	}

	if (key == apikey || startdate.IsZero() || enddate.IsZero()) && parameterfailure == 0 {
		//agtadata := store.AGTAQueryReport(time.Time{}, time.Time{})
		agtadata := store.AGTAQueryReport(startdate, enddate)

		for indx := 0; indx < len(agtadata); indx++ {
			log.Printf("Reservation ID: %d", agtadata[indx].ReservationID)
		}

		CreateExcelFile(agtadata)

		http.ServeFile(w, r, "./import.xlsx")
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

//ElavonHandler - sent to Elavon for payment
func ElavonHandler(w http.ResponseWriter, r *http.Request) {

	// Provide Converge Credentials
	//Converge 6-Digit Account ID *Not the 10-Digit Elavon Merchant ID*
	//const MERCHANTID = 631103
	const MERCHANTID = "011427"
	//Converge User ID *MUST FLAG AS HOSTED API USER IN CONVERGE UI*
	const USERID = "webpage"
	//Converge PIN (64 CHAR A/N)
	//const PIN = "80KYG17V8IBW89MTJYJZIQ3C31DCCG9BJRYP9IYZ4D83ZGQEHCUQDVZB2YBSIG7S"
	const PIN = "WVYWOH"
	const CVVINDICATOR = '1' //means "present"
	//demo url
	const PAYSUCCESSURL = "/paysuccess.php"
	const PAYFAILURL = "/payerror.php"
	const CARDTYPE = "CREDITCARD"

	// URL to Converge demo session token server
	const tokenurl = "https://api.demo.convergepay.com/hosted-payments/transaction_token"
	// URL to Converge production session token server
	//const tokenurl := "https://api.convergepay.com/hosted-payments/transaction_token"

	// URL to the demo Hosted Payments Page
	hppurl := "https://api.demo.convergepay.com/hosted-payments"
	// URL to the production Hosted Payments Page
	//hppurl := "https://api.convergepay.com/hosted-payments"

	//err := r.ParseForm()

	//Follow the above pattern to add additional fields to be sent in curl request below
	requestBody, err := json.Marshal(map[string]string{
		"ssl_merchant_id":      MERCHANTID,
		"ssl_user_id":          USERID,
		"ssl_pin":              PIN,
		"ssl_transaction_type": "CCSALE",
		"ssl_amount":           "1.00",
	})

	if err != nil {
		log.Printf("Error marshaling json %s", err)
	}

	log.Printf("Request Body: %s", string(requestBody))

	resp, err := http.Post(tokenurl, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Printf("Error creating post request to %s: %s", tokenurl, err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("Error reading post response: %s", err)
	}

	log.Printf("Response Body: %s", string(respBody))

	hppurl = hppurl + "?ssl_txn_auth_token=" + string(respBody)

	http.Redirect(w, r, hppurl, http.StatusFound)

}

//ApprovedHander - approved confirmation from Elavon
func ApprovedHandler(w http.ResponseWriter, r *http.Request) {

}

//DeclinedHander - decline confirmation from Elavon
func DeclinedHandler(w http.ResponseWriter, r *http.Request) {

}

//ErrorHander - error confirmation from Elavon
func ErrorHandler(w http.ResponseWriter, r *http.Request) {

}
