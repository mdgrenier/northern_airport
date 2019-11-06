SELECT r.reservationid, flighttime, (select name FROM airlines WHERE airlineid=returnairlineid) AS AirlineName, 
	flightnumber, '' as FlightCity, (select terminal FROM airlines WHERE airlineid=returnairlineid) AS TerminalName, 
    r.reservationid AS confirmationnumber,  
    returnnumchildren + returnnumstudents + returnnumadults + returnnumseniors as NumPax, 
	IF (returndestinationvenueid=100, dropoffaddress, 
	(SELECT name FROM venues WHERE venueid=returndestinationvenueid)) AS DropLocation, 
	(select name FROM cities WHERE cityid= returndestinationcityid) AS DropCity, internalnotes, drivernotes, dt.departuretime,  
	(SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE driverid=t.driverid) AS DriverName , t.driverid,
	(SELECT licenseplate FROM vehicles WHERE vehicleid=t.vehicleid) AS VehicleNum,(cancelled is null or cancelled = 0) AS IsValid, 
	returndate, IF(returndeparturevenueid=99, '', (SELECT name FROM venues WHERE venueid=returndeparturevenueid)) AS HotelInfo
FROM northernairport.reservations r JOIN northernairport.trips t
	ON r.tripid = t.tripid
	JOIN northernairport.vehicles v	ON t.vehicleid = v.vehicleid
    JOIN northernairport.departuretimes dt ON dt.departuretimeid = r.returndeparturetimeid
WHERE destinationcityid=2 and (cancelled is null or cancelled = 0);