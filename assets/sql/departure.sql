USE northernairport;

SELECT r.reservationid, flighttime, (select name FROM airlines WHERE airlineid=departureairlineid) AS AirlineName, 
	flightnumber, '' as FlightCity, (select terminal FROM airlines WHERE airlineid=departureairlineid) AS TerminalName, 
    r.reservationid AS confirmationnumber,
    departurenumchildren + departurenumstudents + departurenumadults + departurenumseniors as NumPax, 
	IF (destinationvenueid=100, dropoffaddress, 
	(SELECT name FROM venues WHERE venueid=destinationvenueid)) AS DropLocation, 
	(select name FROM cities WHERE cityid=destinationcityid) AS DropCity, internalnotes, drivernotes, dt.departuretime,  
	(SELECT CONCAT(lastname, ', ' , firstname) FROM drivers WHERE driverid=t.driverid) AS DriverName , t.driverid,
	(SELECT licenseplate FROM vehicles WHERE vehicleid=t.vehicleid) AS VehicleNum, (cancelled is null or cancelled = 0) AS IsValid, 
	r.departuredate, IF(departurevenueid=99, '', (SELECT name FROM venues WHERE venueid=departurevenueid)) AS HotelInfo,
    cancelled
FROM northernairport.reservations r JOIN northernairport.trips t ON r.tripid = t.tripid
	JOIN northernairport.vehicles v	ON t.vehicleid = v.vehicleid
    JOIN northernairport.departuretimes dt ON dt.departuretimeid = r.departuretimeid
WHERE departurecityid=2 and (cancelled is null or cancelled = 0);