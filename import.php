<?php

include "../private/Library/config.php";

if( !isset($_GET['key']) OR $_GET['key'] != $apiKey )
{
	die('Permission Denied');
}

$conn = mysql_connect($dbhost, $dbuser, $dbpass);
if(! $conn )
{
  die('Could not connect: ' . mysql_error());
}

$sql = "SELECT r.id, flighttime, (select name FROM airlines WHERE id=airline) AS AirlineName, flightnumber, '' as FlightCity, " .
"(select terminal FROM airlines WHERE id=airline) AS TerminalName, confirmationnumber, CONCAT(lastname, ', ', firstname) AS PaxName, departurechildren + departurestudents + departureadults + departureseniors as NumPax, " .
"IF (destinationvenue=100, destinationdropoffaddress, (SELECT name FROM venues WHERE id=destinationvenue)) AS DropLocation, (select name FROM cities WHERE ID = destinationcity) AS DropCity, clientnotes, drivernotes,  departuretime,  (SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE id=tv.driverid) AS DriverName , tv.driverid, " .
"(SELECT licenseplate FROM vehicles WHERE id=tv.vehicleid) AS VehicleNum,(datecancelled is null) AS IsValid, departuredate, IF(departurevenue=99, '', (SELECT name FROM venues WHERE id=departurevenue)) AS HotelInfo " .
"FROM 201655_a.reservations r JOIN 201655_a.trips t " .
"ON r.id = t.reservationid " .
"AND r.departuredate = t.tripdate " .
"JOIN 201655_a.tripvehicles tv " .
"ON t.triptypeid = tv.triptypeid and t.tripdate = tv.tripdate " .
"WHERE departurecity=2 and datecancelled is null ";

if(isset($_GET['startDate']))
{
	$sql .= " AND r.departuredate >= '" . $_GET['startDate'] . "' ";
}

if(isset($_GET['endDate']))
{
	$sql .= " AND r.departuredate < '" . $_GET['endDate'] . "' ";
}


mysql_select_db('201655_a');
$retval = mysql_query( $sql, $conn );
if(! $retval )
{
  die('Could not get data: ' . mysql_error());
}



/** Error reporting */
error_reporting(E_ALL);

ini_set('display_startup_errors',1);
ini_set('display_errors',1);
error_reporting(-1);

/** PHPExcel */
include '/private/PHPExcel.php';

/** PHPExcel_Writer_Excel2007 */
include '/private/PHPExcel/Writer/Excel2007.php';

$objPHPExcel = new PHPExcel();
$objPHPExcel->setActiveSheetIndex(0);
$objPHPExcel->getActiveSheet()->setTitle('Simple');

$index = 1;
$objPHPExcel->getActiveSheet()->SetCellValue('A' . $index, "ID");
$objPHPExcel->getActiveSheet()->SetCellValue('B' . $index, "Flight Arr");
$objPHPExcel->getActiveSheet()->SetCellValue('C' . $index, "Airline");
$objPHPExcel->getActiveSheet()->SetCellValue('D' . $index, "Flight No");
$objPHPExcel->getActiveSheet()->SetCellValue('E' . $index, "Flight City");
$objPHPExcel->getActiveSheet()->SetCellValue('F' . $index, "Terminal");
$objPHPExcel->getActiveSheet()->SetCellValue('G' . $index, "Confirmation/Docket.");
$objPHPExcel->getActiveSheet()->SetCellValue('H' . $index, "Ref .");
$objPHPExcel->getActiveSheet()->SetCellValue('I' . $index, "Pax Name");
$objPHPExcel->getActiveSheet()->SetCellValue('J' . $index, "Pax");
$objPHPExcel->getActiveSheet()->SetCellValue('K' . $index, "Drop Location");
$objPHPExcel->getActiveSheet()->SetCellValue('L' . $index, "Drop City");
$objPHPExcel->getActiveSheet()->SetCellValue('M' . $index, "Zone");
$objPHPExcel->getActiveSheet()->SetCellValue('N' . $index, "Service Type");
$objPHPExcel->getActiveSheet()->SetCellValue('O' . $index, "Additional Notes 1");
$objPHPExcel->getActiveSheet()->SetCellValue('P' . $index, "Additional Notes 2");
$objPHPExcel->getActiveSheet()->SetCellValue('Q' . $index, "EDT");
$objPHPExcel->getActiveSheet()->SetCellValue('R' . $index, "Driver Name");
$objPHPExcel->getActiveSheet()->SetCellValue('S' . $index, "Driver .");
$objPHPExcel->getActiveSheet()->SetCellValue('T' . $index, "Vehicle .");
$objPHPExcel->getActiveSheet()->SetCellValue('U' . $index, "At Booth");
$objPHPExcel->getActiveSheet()->SetCellValue('V' . $index, "Check In");
$objPHPExcel->getActiveSheet()->SetCellValue('W' . $index, "Depart Time");
$objPHPExcel->getActiveSheet()->SetCellValue('X' . $index, "Driver Departed");
$objPHPExcel->getActiveSheet()->SetCellValue('Y' . $index, "ReservationDate");
$objPHPExcel->getActiveSheet()->SetCellValue('Z' . $index, "HotelInfo");


$index = 2;
while($row = mysql_fetch_array($retval, MYSQL_ASSOC))
{
    $objPHPExcel->getActiveSheet()->SetCellValue('A' . $index, $row["id"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('B' . $index, $row["flighttime"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('C' . $index, $row["AirlineName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('D' . $index, $row["flightnumber"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('E' . $index, $row["FlightCity"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('F' . $index, $row["TerminalName"]);
	$objPHPExcel->getActiveSheet()->setCellValueExplicit('G' . $index,$row["confirmationnumber"],PHPExcel_Cell_DataType::TYPE_STRING);
	$objPHPExcel->getActiveSheet()->SetCellValue('H' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('I' . $index, $row["PaxName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('J' . $index, $row["NumPax"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('K' . $index, $row["DropLocation"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('L' . $index, $row["DropCity"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('M' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('N' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('O' . $index, $row["drivernotes"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('P' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('Q' . $index, $row["departuretime"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('R' . $index, $row["DriverName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('S' . $index, $row["driverid"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('T' . $index, $row["VehicleNum"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('U' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('V' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('W' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('X' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('Y' . $index, $row["departuredate"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('Z' . $index, $row["HotelInfo"]);
	$index++;
} 

//RETURNS

$sql = "SELECT r.id, flighttime, (select name FROM airlines WHERE id=returnairline) AS AirlineName, flightnumber, '' as FlightCity, " .
"(select terminal FROM airlines WHERE id=returnairline) AS TerminalName, confirmationnumber, CONCAT(lastname, ', ', firstname) AS PaxName, returnchildren + returnstudents + returnadults + returnseniors as NumPax, " .
"IF (returndestinationvenue=100, returndestinationdropoffaddress, (SELECT name FROM venues WHERE id=returndestinationvenue)) AS DropLocation, (select name FROM cities WHERE ID = returndestinationcity) AS DropCity, clientnotes, drivernotes,  returntime,  (SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE id=tv.driverid) AS DriverName , tv.driverid, " .
"(SELECT licenseplate FROM vehicles WHERE id=tv.vehicleid) AS VehicleNum,(datecancelled is null) AS IsValid, returndate, IF(returndeparturevenue=99, '', (SELECT name FROM venues WHERE id=returndeparturevenue)) AS HotelInfo " .
"FROM 201655_a.reservations r JOIN 201655_a.trips t " .
"ON r.id = t.reservationid " .
"AND r.returndate = t.tripdate " .
"JOIN 201655_a.tripvehicles tv " .
"ON t.triptypeid = tv.triptypeid and t.tripdate = tv.tripdate " .
"JOIN 201655_a.triptypes tt ON tt.id = t.triptypeid " .
"WHERE destinationcity=2 and tt.direction='N' and datecancelled is null ";

if(isset($_GET['startDate']))
{
	$sql .= " AND r.returndate >= '" . $_GET['startDate'] . "' ";
}

if(isset($_GET['endDate']))
{
	$sql .= " AND r.returndate < '" . $_GET['endDate'] . "' ";
}


mysql_select_db('201655_a');
$retval = mysql_query( $sql, $conn );
if(! $retval )
{
  die('Could not get data: ' . mysql_error());
}
while($row = mysql_fetch_array($retval, MYSQL_ASSOC))
{
    //echo $row["PaxName"] . "<br />";
	$objPHPExcel->getActiveSheet()->SetCellValue('A' . $index, $row["id"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('B' . $index, $row["flighttime"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('C' . $index, $row["AirlineName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('D' . $index, $row["flightnumber"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('E' . $index, $row["FlightCity"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('F' . $index, $row["TerminalName"]);
	$objPHPExcel->getActiveSheet()->setCellValueExplicit('G' . $index,$row["confirmationnumber"],PHPExcel_Cell_DataType::TYPE_STRING);
	$objPHPExcel->getActiveSheet()->SetCellValue('H' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('I' . $index, $row["PaxName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('J' . $index, $row["NumPax"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('K' . $index, $row["DropLocation"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('L' . $index, $row["DropCity"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('M' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('N' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('O' . $index, $row["drivernotes"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('P' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('Q' . $index, $row["returntime"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('R' . $index, $row["DriverName"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('S' . $index, $row["driverid"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('T' . $index, $row["VehicleNum"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('U' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('V' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('W' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('X' . $index, "");
	$objPHPExcel->getActiveSheet()->SetCellValue('Y' . $index, $row["returndate"]);
	$objPHPExcel->getActiveSheet()->SetCellValue('Z' . $index, $row["HotelInfo"]);
	$index++;
} 


mysql_close($conn);

header("Pragma: public");
header("Expires: 0");
header("Cache-Control: must-revalidate, post-check=0, pre-check=0"); 
header("Content-Type: application/force-download");
header("Content-Type: application/octet-stream");
header("Content-Type: application/download");;
header("Content-Disposition: attachment;filename=test.xlsx"); 
header("Content-Transfer-Encoding: binary ");
$objWriter = new PHPExcel_Writer_Excel2007($objPHPExcel);
$objWriter->save('php://output');
?>