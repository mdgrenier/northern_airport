package main

import (
	"database/sql"
	"encoding/gob"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/gorilla/securecookie"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

var cache redis.Conn

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

	//connString := "dbname=mydb user=mdgre_000 password=password port=5432 sslmode=disable"
	//db, err := sql.Open("postgres", connString)
	//local development connection
	//connString := "root:test@tcp(127.0.0.1:3306)/northernairport"
	//docker container connection
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

	gob.Register(User{})

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}
