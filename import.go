package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/tealeg/xlsx"
)

// CreateExcelFile - test file
func CreateExcelFile(agtareport []AGTAReport) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var err error

	log.Printf("Creating log file")

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Simple")
	if err != nil {
		fmt.Printf(err.Error())
	}

	log.Printf("File created")

	CreateHeaderRow(sheet)

	PopulateData(sheet, agtareport)

	err = file.Save("test.xlsx")

	if err != nil {
		fmt.Printf(err.Error())
	}

	log.Printf("File Saved")
}

//CreateHeaderRow - create the header row
func CreateHeaderRow(sheet *xlsx.Sheet) {
	row1 := sheet.AddRow()
	cell1 := row1.AddCell()
	cell1.Value = "ID"

	cell2 := row1.AddCell()
	cell2.Value = "Flight Arr"

	cell3 := row1.AddCell()
	cell3.Value = "Airline"

	cell4 := row1.AddCell()
	cell4.Value = "Flight No"

	cell5 := row1.AddCell()
	cell5.Value = "Flight City"

	cell6 := row1.AddCell()
	cell6.Value = "Terminal"

	cell7 := row1.AddCell()
	cell7.Value = "Confirmation/Docket."

	cell8 := row1.AddCell()
	cell8.Value = "Ref ."

	cell9 := row1.AddCell()
	cell9.Value = "Pax Name"

	cell10 := row1.AddCell()
	cell10.Value = "Pax"

	cell11 := row1.AddCell()
	cell11.Value = "Drop Location"

	cell12 := row1.AddCell()
	cell12.Value = "Drop City"

	cell13 := row1.AddCell()
	cell13.Value = "Zone"

	cell14 := row1.AddCell()
	cell14.Value = "Service Type"

	cell15 := row1.AddCell()
	cell15.Value = "Additional Notes 1"

	cell16 := row1.AddCell()
	cell16.Value = "Additional Notes 2"

	cell17 := row1.AddCell()
	cell17.Value = "EDT"

	cell18 := row1.AddCell()
	cell18.Value = "Driver Name"

	cell19 := row1.AddCell()
	cell19.Value = "Driver ."

	cell20 := row1.AddCell()
	cell20.Value = "Vehicle ."

	cell21 := row1.AddCell()
	cell21.Value = "At Booth"

	cell22 := row1.AddCell()
	cell22.Value = "Check In"

	cell23 := row1.AddCell()
	cell23.Value = "Depart Time"

	cell24 := row1.AddCell()
	cell24.Value = "Driver Departed"

	cell25 := row1.AddCell()
	cell25.Value = "ReservationDate"

	cell26 := row1.AddCell()
	cell26.Value = "HotelInfo"
}

//PopulateData - populate data in excel file
func PopulateData(sheet *xlsx.Sheet, agtareport []AGTAReport) {

	for indx := 0; indx < len(agtareport); indx++ {
		row1 := sheet.AddRow()
		cell1 := row1.AddCell()
		cell1.Value = strconv.Itoa(agtareport[indx].ReservationID)

		cell2 := row1.AddCell()
		cell2.Value = agtareport[indx].FlightTime

		cell3 := row1.AddCell()
		cell3.Value = agtareport[indx].AirlineName

		cell4 := row1.AddCell()
		cell4.Value = agtareport[indx].FlightNumber

		cell5 := row1.AddCell()
		cell5.Value = agtareport[indx].FlightCity

		cell6 := row1.AddCell()
		cell6.Value = agtareport[indx].TerminalName

		cell7 := row1.AddCell()
		cell7.Value = strconv.Itoa(agtareport[indx].ConfirmationNumber)

		cell8 := row1.AddCell()
		cell8.Value = ""

		cell9 := row1.AddCell()
		cell9.Value = agtareport[indx].PaxName

		cell10 := row1.AddCell()
		cell10.Value = strconv.Itoa(agtareport[indx].NumPax)

		cell11 := row1.AddCell()
		cell11.Value = agtareport[indx].DropLocation

		cell12 := row1.AddCell()
		cell12.Value = agtareport[indx].DropCity

		//Zone
		cell13 := row1.AddCell()
		cell13.Value = ""

		//Service Type
		cell14 := row1.AddCell()
		cell14.Value = ""

		cell15 := row1.AddCell()
		cell15.Value = agtareport[indx].DriverNotes

		cell16 := row1.AddCell()
		cell16.Value = agtareport[indx].InternalNotes

		cell17 := row1.AddCell()
		cell17.Value = strconv.Itoa(agtareport[indx].DepartureTime)

		cell18 := row1.AddCell()
		cell18.Value = agtareport[indx].DriverName

		cell19 := row1.AddCell()
		cell19.Value = strconv.Itoa(agtareport[indx].DriverID)

		cell20 := row1.AddCell()
		cell20.Value = strconv.Itoa(agtareport[indx].DepartureTime) + " seats"

		cell21 := row1.AddCell()
		cell21.Value = ""

		cell22 := row1.AddCell()
		cell22.Value = ""

		cell23 := row1.AddCell()
		cell23.Value = ""

		cell24 := row1.AddCell()
		cell24.Value = ""

		cell25 := row1.AddCell()
		cell25.Value = agtareport[indx].DepartureDate.Format("2006-01-02")

		cell26 := row1.AddCell()
		cell26.Value = agtareport[indx].HotelInfo
	}
}

//DownloadFile - prompt download to being
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
