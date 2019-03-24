package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	client := &Client{}

	err := r.ParseForm()

	if err != nil {
		log.Fatal("Parse Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	log.Print("user created")

	http.Redirect(w, r, "/", http.StatusFound)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := &Client{}

	err = r.ParseForm()

	if err != nil {
		log.Fatal("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client.Username = r.Form.Get("username")
	client.Password = r.Form.Get("password")

	err = store.SignInUser(client)
	if err != nil {
		log.Print("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	client.Authenticated = true

	session.Values["client"] = client

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	log.Printf("%s is signed in", session.Values["client"])

	http.Redirect(w, r, "/reservation", http.StatusFound)
}

// logout revokes authentication for a user
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["client"] = Client{}
	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	log.Printf("username: %s", client.Username)

	tpl.ExecuteTemplate(w, "index.gohtml", client)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	tpl.ExecuteTemplate(w, "signup.gohtml", client)
}

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	if client.Authenticated {
		err := store.GetClientInfo(&client)

		if err != nil {
			log.Panicf("Error getting client info: %s", err.Error())
		}
	}

	venues := store.GetVenues()

	cities := store.GetCities()

	reservation := Reservation{}

	reservation.Client = client
	reservation.Venues = venues
	reservation.Cities = cities

	//for i, venue := range reservation.venues {
	//	log.Printf("Venue Name %d: %s", i, venue.VenueName)
	//}
	//
	//log.Printf("Firstname: %s", reservation.client.Firstname

	tpl.ExecuteTemplate(w, "reservation.gohtml", reservation)
}

func createreservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := GetClient(session)

	tpl.ExecuteTemplate(w, "created.gohtml", client)
}
