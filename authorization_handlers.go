package main

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Credentials creates a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}

// User stores authentication information
type User struct {
	Username      string
	Authenticated bool
}

var sessionStore *sessions.CookieStore

func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}

	err := r.ParseForm()

	if err != nil {
		log.Fatal("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	creds.Username = r.Form.Get("username")
	creds.Password = r.Form.Get("password")

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		log.Fatal("Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	creds.Password = string(hashedPassword)

	err = store.CreateUser(creds)
	if err != nil {
		log.Fatal("Error: ", err)
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

	creds := &Credentials{}

	err = r.ParseForm()

	if err != nil {
		log.Fatal("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	creds.Username = r.Form.Get("username")
	creds.Password = r.Form.Get("password")

	err = store.SignInUser(creds)
	if err != nil {
		log.Print("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	user := &User{
		Username:      creds.Username,
		Authenticated: true,
	}

	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	log.Printf("%s is signed in", session.Values["user"])

	http.Redirect(w, r, "/reservation", http.StatusFound)
}

// logout revokes authentication for a user
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = User{}
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

	user := getUser(session)

	log.Printf("username: %s", user.Username)

	tpl.ExecuteTemplate(w, "index.gohtml", user)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated {
		session.Values["user"] = User{}
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(w, "signup.gohtml", user)
}

func reservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated {
		session.Values["user"] = User{}
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(w, "reservation.gohtml", user)
}

func createreservationHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, "northern-airport")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := getUser(session)

	if user.Authenticated {
		session.Values["user"] = User{}
		session.Options.MaxAge = -1

		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//toast notification that reservations has been created
	tpl.ExecuteTemplate(w, "created.gohtml", user)
}
