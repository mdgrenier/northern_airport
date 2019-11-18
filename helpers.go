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

		departureairline, _ := strconv.Atoi(r.FormValue("departureairline"))
		if departureairline > 0 {
			reservation.DepartureAirlineID = departureairline
		}

		destinationairline, _ := strconv.Atoi(r.FormValue("destinationairline"))
		if departureairline > 0 {
			reservation.DepartureAirlineID = destinationairline
		}

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

			returndepartureairline, _ := strconv.Atoi(r.FormValue("returndepartureairline"))
			if returndepartureairline > 0 {
				reservation.ReturnAirlineID = returndepartureairline
			}

			returndestinationairline, _ := strconv.Atoi(r.FormValue("returndestinationairline"))
			if returndepartureairline > 0 {
				reservation.ReturnAirlineID = returndestinationairline
			}

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

		reservation.FlightNumber, err = strconv.Atoi(r.Form.Get("flightnumber"))
		reservation.FlightTime, err = strconv.Atoi(r.Form.Get("flighttime"))

		reservation.InternalNotes = r.Form.Get("internalnotes")
		reservation.DriverNotes = r.Form.Get("drivernotes")

		//must map discount code to discountcodeid
		//reservation.DiscountCodeID, err = strconv.Atoi(r.FormValue("discountcode"))

		discount := DiscountCode{}

		discount.Name = r.Form.Get("promocode")

		err := store.GetDiscount(&discount)

		if err != nil {
			log.Printf("Error retrieving discount code: %s", err.Error())
		} else {
			reservation.DiscountCodeID = discount.DiscountCodeID
		}

		var tripprice float64

		log.Printf("Trip price string value: %s", r.Form.Get("tripprice"))

		tripprice, err = strconv.ParseFloat(r.Form.Get("tripprice"), 32)

		if err != nil {
			log.Printf("Error converting trip price to float: %s", err.Error())
		}

		reservation.Price = float32(tripprice)

		log.Printf("reservation price: %f", reservation.Price)

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

//SendConfirmationEmail - given the recipient and email body, send a email
func SendConfirmationEmail(to string, departuredetails string, returndetails string, reservationid int, price float32) {

	from := "mdgrenier@gmail.com"
	pass := "Dh76nm6m*"
	//to := "matt@mgrenier.ca"
	body := fmt.Sprintf("Northern Airport Passenger Service Confirmation\n\n"+
		"Your Confirmation ID is #%d\n\nDeparture information:\n%s\n\n", reservationid, departuredetails)

	if returndetails != "" {
		body += "Return information:\n" + returndetails + "\n\n"
	}

	body += "Total Cost of the trip is $" + fmt.Sprintf("%f", price) + "\n\n\n"

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

//SendCapacityEmail - given trip create capacity email
func SendCapacityEmail(trip Trips) {
	from := "mdgrenier@gmail.com"
	pass := "Dh76nm6m*"
	to := "matt@mgrenier.ca"

	body := fmt.Sprintf("Northern Airport Passenger Service Capacity Alert\n\n"+
		"TripID #%d has reached 75%% capacity.\n\n\n", trip.TripID)

	var drivername string

	for _, driver := range trip.DriverList {
		if driver.DriverID == trip.DriverID {
			drivername = driver.FirstName + " " + driver.LastName
		}
	}

	var vehicleplate string

	for _, vehicle := range trip.VehicleList {
		if vehicle.VehicleID == trip.VehicleID {
			vehicleplate = vehicle.LicensePlate
		}
	}

	body += fmt.Sprintf("Trip Details:\nDeparture Date: %s\nDeparture Time: %d\n"+
		"Number of Passengers: %d\nDriver: %s\nVehicle: %s\nCapacity: %d\n",
		trip.DepartureDate, trip.DepartureTime, trip.NumPassengers, drivername, vehicleplate, trip.Capacity)

	msg := "From: " + from + "\n" +
		"To: mdgrenier@gmail.com\n" +
		"Subject: Northern Airport Reservation Confirmation\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

//NumDaysLookup - given month and year return number of days in the month
func NumDaysLookup(month int, year int) int {
	switch month {
	case 1:
		//January
		return 31
	case 2:
		//February
		//leap year calculation
		if year%400 == 0 {
			return 29
		} else if year%100 == 0 {
			return 28
		} else if year%4 == 0 {
			return 29
		} else {
			return 28
		}
	case 3:
		//March
		return 31
	case 4:
		//April
		return 30
	case 5:
		//May
		return 31
	case 6:
		//June
		return 30
	case 7:
		//July
		return 31
	case 8:
		return 30
	case 9:
		return 31
	case 10:
		return 31
	case 11:
		return 30
	case 12:
		return 31
	default:
		return 0
	}
}
