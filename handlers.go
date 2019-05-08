package main

import (
	"log"
	"net/http"
	"strconv"

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
		log.Print("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
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

	//store retrieved data in reservation structure and pass to template
	resform := ResFormData{}

	resform.Client = client
	resform.Venues = venues
	resform.VenueCount = len(venues)
	resform.Cities = cities
	resform.DepartureTimes = departuretimes

	tpl.ExecuteTemplate(w, "reservation.gohtml", resform)
}

//CreateReservationHandler - store reservation in database
func CreateReservationHandler(w http.ResponseWriter, r *http.Request) {
	//session, err := sessionStore.Get(r, "northern-airport")
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//resformdata := ResFormData{}
	reservation := Reservation{}

	//get reservation data from form
	reservation = GetReservationFormValues(r, true)

	//check if trip exists, if not create one
	err := store.GetOrAddTrip(&reservation)

	err = store.CreateReservation(&reservation)

	if err != nil {
		log.Fatal("Error Creating User: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Print("reservation created")

	http.Redirect(w, r, "/reservationcreated", http.StatusFound)
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
		tpl.ExecuteTemplate(w, "index.gohtml", r)
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
		tpl.ExecuteTemplate(w, "index.gohtml", r)
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
		tpl.ExecuteTemplate(w, "index.gohtml", r)
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

		//get data need to populate dropdowns in reservation form
		trips := store.GetTrips()

		for i := 0; i < len(trips); i++ {
			tripid, err := strconv.Atoi(values["tripid"][0])

			if err != nil {
				log.Printf("Error converting tripid: %s", err.Error())
			}

			if trips[i].TripID == tripid {
				log.Printf("Matched Trip ID: %d", trips[i].TripID)
				trips[i].DriverID, err = strconv.Atoi(values["driverid"][0])
				trips[i].VehicleID, err = strconv.Atoi(values["vehicleid"][0])

				store.UpdateTrip(&trips[i])
			}

			if err != nil {
				log.Printf("Error converting driverid or vehilceid: %s", err.Error())
			}

		}

		tpl.ExecuteTemplate(w, "trip.gohtml", trips)
	} else {
		tpl.ExecuteTemplate(w, "index.gohtml", r)
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
		venues := store.GetVenues()

		venues[0].RoleID = client.RoleID

		tpl.ExecuteTemplate(w, "venue.gohtml", venues)
	} else {
		tpl.ExecuteTemplate(w, "index.gohtml", r)
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
