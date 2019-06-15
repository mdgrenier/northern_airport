package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
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

//SendEmail - given the recipient and email body, send a email
func SendEmail(to string, departuredetails string, returndetails string, reservationid int) {
	log.Print("Attempt to send email")

	from := "mdgrenier@gmail.com"
	pass := "Dh76nm6m*"
	//to := "matt@mgrenier.ca"
	body := fmt.Sprintf("Northern Airport Passenger Service Confirmation\n\n"+
		"Your Confirmation ID is #%d\n\nDeparture information:\n%s\n\n", reservationid, departuredetails)

	if returndetails != "" {
		body += "Return information:\n" + returndetails + "\n\n\n"
	}

	body += "For any concerns please contact us at (705) 474-7942\n\n" +
		"Thank you for booking with Northern Airport Passenger Service."

	msg := "From: " + from + "\n" +
		"To: " + to + "; mdgrenier@gmail.com\n" +
		"Subject: Northern Airport Reservation Confirmation\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit http://reserve.northernairport.ca")
}

//FormatTripDetails - given trip details return a string for use in confirmation emails
func FormatTripDetails(departurecity string, departurevenue string, departuredate string, departuretime string,
	destinationcity string, destinationvenue string, numadults string, numseniors string, numstudents string,
	numchildren string) string {

	departureinfo := "Departing from " + departurevenue + " in " + departurecity + "\n" +
		"Departing on " + departuredate + " at " + departuretime + "\n"

	destinationinfo := "Arriving at " + destinationvenue + " in " + destinationcity + "\n"

	passengers := "Trip will include:\n" + numadults + " Adults, " + numseniors + " Seniors, " +
		numstudents + " Students, and " + numchildren + " Children"

	return departureinfo + destinationinfo + passengers
}
