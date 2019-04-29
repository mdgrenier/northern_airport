package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
)

//GetClient - return the signed in user stored in a session cookie
func GetClient(s *sessions.Session) Client {
	val := s.Values["client"]
	var client = Client{}
	client, ok := val.(Client)

	if !ok {
		return Client{Authenticated: false}
	}
	return client
}

//GetReservationFormValues - get the field values from the reservation form
func GetReservationFormValues(r *http.Request, gettripdata bool) Reservation {

	//resform := ResFormData{}
	reservation := Reservation{}

	//populate Form structure of the http request
	err := r.ParseForm()

	if err != nil {
		log.Fatal("Parse Error: ", err)
		return reservation
	}

	//store client info in resformdata
	//resform.Client.ClientID, err = strconv.Atoi(r.Form.Get("clientid"))
	if err != nil {
		log.Panicf("Error converting ClientID to string: %s", err.Error())
	}

	/*********************************************************************
	*
	*	compare existing client data to form data, if different prompt user
	*	to save new info or just use for this reservation
	*
	**********************************************************************/

	//resform.Client.Username = r.Form.Get("username")
	//resform.Client.Password = r.Form.Get("password")
	//resform.Client.Firstname = r.Form.Get("firstname")
	//resform.Client.Lastname = r.Form.Get("lastname")
	//resform.Client.Phone = r.Form.Get("phone")
	//resform.Client.Email = r.Form.Get("email")
	//resform.Client.StreetAddress = r.Form.Get("streetaddress")
	//resform.Client.City = r.Form.Get("city")
	//resform.Client.Province = r.Form.Get("provstate")
	//resform.Client.PostalCode = r.Form.Get("postalzip")
	//resform.Client.Country = r.Form.Get("country")

	reservation.ClientID, err = strconv.Atoi(r.FormValue("clientid"))

	//store trip info in resformdata
	if gettripdata {
		reservation.TripTypeID, err = strconv.Atoi(r.FormValue("triptype"))

		reservation.DepartureCityID, err = strconv.Atoi(r.FormValue("departurecity"))
		reservation.DepartureVenueID, err = strconv.Atoi(r.FormValue("departurevenue"))
		reservation.DestinationCityID, err = strconv.Atoi(r.FormValue("destinationcity"))
		reservation.DestinationVenueID, err = strconv.Atoi(r.FormValue("destinationvenue"))
		reservation.DepartureDate, err = time.Parse("2006-01-02", r.FormValue("departuredate"))
		reservation.DepartureTimeID, err = strconv.Atoi(r.FormValue("departuretime"))

		reservation.DepartureNumAdults, err = strconv.Atoi(r.FormValue("departurenumadults"))
		reservation.DepartureNumSeniors, err = strconv.Atoi(r.FormValue("departurenumseniors"))
		reservation.DepartureNumStudents, err = strconv.Atoi(r.FormValue("departurenumstudents"))
		reservation.DepartureNumChildren, err = strconv.Atoi(r.FormValue("departurenumchildren"))

		if reservation.TripTypeID == 2 {
			reservation.ReturnDepartureCityID, err = strconv.Atoi(r.FormValue("returndeparturecity"))
			reservation.ReturnDepartureVenueID, err = strconv.Atoi(r.FormValue("returndeparturevenue"))
			reservation.ReturnDestinationCityID, err = strconv.Atoi(r.FormValue("returndestinationcity"))
			reservation.ReturnDestinationVenueID, err = strconv.Atoi(r.FormValue("returndestinationvenue"))
			reservation.ReturnDate, err = time.Parse("2006-01-02", r.FormValue("returndeparturedate"))
			reservation.ReturnDepartureTimeID, err = strconv.Atoi(r.FormValue("returndeparturetime"))

			reservation.ReturnNumAdults, err = strconv.Atoi(r.FormValue("returnnumadults"))
			reservation.ReturnNumSeniors, err = strconv.Atoi(r.FormValue("returnnumseniors"))
			reservation.ReturnNumStudents, err = strconv.Atoi(r.FormValue("returnnumstudents"))
			reservation.ReturnNumChildren, err = strconv.Atoi(r.FormValue("returnnumchildren"))
		} else {
			reservation.ReturnDepartureCityID = 0
			reservation.ReturnDepartureVenueID = 0
			reservation.ReturnDestinationCityID = 0
			reservation.ReturnDestinationVenueID = 0
			reservation.ReturnDate, err = time.Parse("2006-01-02", "0000-00-00")
			reservation.ReturnDepartureTimeID = 0

			reservation.ReturnNumAdults = 0
			reservation.ReturnNumSeniors = 0
			reservation.ReturnNumStudents = 0
			reservation.ReturnNumChildren = 0
		}

		reservation.InternalNotes = r.Form.Get("internalnotes")
		reservation.DriverNotes = r.Form.Get("drivernotes")

		//must map discount code to discountcodeid
		reservation.DiscountCodeID, err = strconv.Atoi(r.FormValue("discountcode"))
		var price float64
		price, err = strconv.ParseFloat(r.FormValue("price"), 32)
		reservation.Price = float32(price)

		log.Printf("Reservation values retrieved from form")

		/*
			Status
			Hash
			CustomDepartureID
			CustomDestinationID
			TripTypeID
			BalanceOwing
			ElavonTransactionID
		*/
	}
	return reservation
}
