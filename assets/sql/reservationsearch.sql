use northernairport;

select departuredate, t.departuretimeid, departuretime,
	sum(numpassengers), t.driverid, t.vehicleid, capacity, omitted,
    d.firstname, d.lastname, v.licenseplate
from trips t inner join departuretimes dt on t.departuretimeid = dt.departuretimeid
	inner join drivers d on d.driverid = t.driverid
    inner join vehicles v on v.vehicleid = t.vehicleid
group by departuredate, t.departuretimeid, departuretime,
	t.driverid, t.vehicleid, capacity, omitted,
    d.firstname, d.lastname, v.licenseplate
order by departuredate desc;
