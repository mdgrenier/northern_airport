package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
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

	reservation := Reservation{}

	log.Printf("Retrieve form data")

	//get reservation data from form
	reservation = GetReservationFormValues(r, true)

	log.Printf("Form information retrieved")

	//check if trip exists, if not create one
	err := store.GetOrAddTrip(&reservation)

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

		http.Redirect(w, r, "/reservationcreated", http.StatusFound)
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
		drivers := store.GetDrivers()

		drivers[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "driver.gohtml", drivers)
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

		drivers := store.GetDrivers()

		drivers[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "driver.gohtml", drivers)
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

		tpl.ExecuteTemplate(w, "driver.gohtml", driver)
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

		driver := store.GetDrivers()

		if err != nil {
			log.Printf("Error converting driverid: %s", err.Error())
		}

		tpl.ExecuteTemplate(w, "driver.gohtml", driver)
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
		vehicles := store.GetVehicles()

		vehicles[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehicles)
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

		vehicles := store.GetVehicles()

		vehicles[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "vehicle.gohtml", vehicles)
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
		trips := store.GetTrips()

		trips[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "trip.gohtml", trips)
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

		tpl.ExecuteTemplate(w, "trip.gohtml", trip)
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
		//get data need to populate dropdowns in reservation form
		searchreservations := store.SearchReservations(name, phone, email)

		if len(searchreservations) > 0 {
			searchreservations[0].RoleID = client.RoleID

			log.Printf("search: we've got some reservations!")
		} else {
			log.Printf("search: no reservations returned!")
		}

		if err := tpl.ExecuteTemplate(w, "searchreservations.gohtml", searchreservations); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//Search2Handler - display reservations
func Search2Handler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("refresh table")

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
		//get data need to populate dropdowns in reservation form
		log.Printf("search reservations")

		searchreservations := store.SearchReservations(name, phone, email)

		if len(searchreservations) > 0 {
			searchreservations[0].RoleID = client.RoleID

			log.Printf("searchreservations: we've got some reservations!")
		} else {
			log.Printf("searchreservations: no reservations returned!")
		}

		//var buf bytes.Buffer
		//if err := tpl.ExecuteTemplate(&buf, "searchreservations.gohtml", searchreservations); err != nil {
		//	log.Printf("Error executing HTML template: %s", err.Error())
		//	http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		//	return
		//}

		searchtemp, err := template.ParseFiles("./templates/layout.gohtml")
		if err != nil {
			log.Printf("Error parsing searchreservations: %s", err.Error())
		} else {
			err = searchtemp.ExecuteTemplate(w, "search", searchreservations)
		}

		if err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
		} else {
			for i, res := range searchreservations {
				log.Printf("ReservationID: %d: %d", i, res.ReservationID)
			}
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

		newsearchreservations := store.SearchReservations("", 0, "")

		newsearchreservations[0].RoleID = client.RoleID

		for i, res := range newsearchreservations {
			log.Printf("#%d ReservationID Returned: %d - Postponed: %t", i, res.ReservationID, res.Postponed)
		}

		if err := tpl.ExecuteTemplate(w, "searchreservations.gohtml", newsearchreservations); err != nil {
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

		newsearchreservations := store.SearchReservations("", 0, "")

		newsearchreservations[0].RoleID = client.RoleID

		for i, res := range newsearchreservations {
			log.Printf("#%d ReservationID Returned: %d - Cancelled: %t", i, res.ReservationID, res.Postponed)
		}

		var buf bytes.Buffer
		//tpl.ExecuteTemplate(&buf, "searchreservations.gohtml", searchreservations)
		if err := tpl.ExecuteTemplate(&buf, "searchreservations.gohtml", newsearchreservations); err != nil {
			log.Printf("Error executing HTML template: %s", err.Error())
			http.Error(w, "Error executing HTML template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		io.Copy(w, &buf)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}

/*
//OmitTripHandler - omit trip
func OmitTripHandler(w http.ResponseWriter, r *http.Request) {
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

		var err error
		trip := Trips{}

		trip.TripID, err = strconv.Atoi(values["tripid"][0])

		if err != nil {
			log.Printf("Error converting tripid: %s", err.Error())
		}

		if values["cancelled"][0] == "1" {
			trip.Cancelled = true
		} else {
			trip.Cancelled = false
		}

		if err != nil {
			log.Printf("Error converting cancelled: %s", err.Error())
		}

		store.CancelTrip(&trip)

		tpl.ExecuteTemplate(w, "trip.gohtml", trip)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}

}
*/

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
		venues := store.GetVenues()

		venues[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "venue.gohtml", venues)
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

		venues := store.GetVenues()

		venues[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "venue.gohtml", venues)
	} else {
		tpl.ExecuteTemplate(w, "accessdenied.gohtml", r)
	}
}

//UpdateVenueHandler - update city in database
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
		cities := store.GetCities()

		cities[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "city.gohtml", cities)
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

		cities := store.GetCities()

		cities[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "city.gohtml", cities)
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

//PriceHandler - return price given the departure, destination and customer type
func PriceHandler(w http.ResponseWriter, r *http.Request) {

	values := r.URL.Query()

	reservationtypeid, err := strconv.Atoi(values["reservationtypeid"][0])
	departurecityid, err := strconv.Atoi(values["departurecityid"][0])
	destinationcityid, err := strconv.Atoi(values["destinationcityid"][0])
	retdeparturecityid, err := strconv.Atoi(values["retdeparturecityid"][0])
	retdestinationcityid, err := strconv.Atoi(values["retdestinationcityid"][0])
	numpassengers, err := strconv.Atoi(values["numpassengers"][0])
	customertypeid, err := strconv.Atoi(values["customertypeid"][0])
	discountcode := values["discount"][0]

	price := store.GetPrice(departurecityid, destinationcityid, retdeparturecityid, retdestinationcityid, customertypeid, reservationtypeid, discountcode)

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
	log.Printf("execute update departure time handler")

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
	log.Printf("execute add departure time handler")

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

		tpl.ExecuteTemplate(w, "departuretime.gohtml", departuretime)
	} else {
		tpl.ExecuteTemplate(w, "departuretime.gohtml", r)
	}

}

//ReportHandler - display reports admin page
func ReportHandler(w http.ResponseWriter, r *http.Request) {

	//session, err := sessionStore.Get(r, "northern-airport")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//get client data from session cookie
	//client := GetClient(session)

	tpl.ExecuteTemplate(w, "reports.gohtml", r)
}
