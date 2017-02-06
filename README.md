# SNMP Mock
This is simple server that responds snmp queries.

## Configuration
### Mock configuration
There are a few things that can be changed with environment variables:
* `SNMP_PORT`: default is 161.
* `SNMP_COMMUNITY`: default is public
* `BASE_OID`: the base oid the mock server is going to handle. i.e.: `1.3.6.100.1.2.3.4`. All oid responses must be within this base oid.
* `APP_PORT`: default is 8080. In the future there will be an endpoint to add oid responses.

### OID responses
The mock server reads all csv files found in `/app/oids`. The csv file has the format `<oid>;<type>;<value>`.

Example:
```
1.3.6.100.1.2.3.5.1.1.0;OctetString;Sample2
1.3.6.100.1.2.3.5.1.2.0;Integer;42
1.3.6.100.1.2.3.5.1.3.0;IPAddress;10.20.30.40
1.3.6.100.1.2.3.5.1.4.0;Counter32;11223344
1.3.6.100.1.2.3.5.1.5.0;Gauge32;88866
1.3.6.100.1.2.3.5.1.6.0;TimeTicks;1876300
1.3.6.100.1.2.3.5.1.7.0;Counter64;111
1.3.6.100.1.2.3.5.1.8.0;ObjectIdentifier;1.3.6.8.9.999.100
1.3.6.100.1.2.3.5.1.9.0;ObjectIdentifier;"Multi
line
response"
```

Available types:
* OctetString
* Integer: 32 bit signed integer
* IPAddress: v4 IP address
* Counter32: 32 bit unsigned integer
* Gauge32: 32 bit unsigned integer
* TimeTickes: 32 bit unsigned integer in hundredths of a second
* Counter64: 64 bit unsigned int
* ObjectIdentifier

## Example
```
$ ls sample_oids/
one.csv  two.csv

$ cat sample_oids/one.csv
1.3.6.100.1.2.3.4.1.1.0;OctetString;Sample 1
1.3.6.100.1.2.3.4.1.3.1.20;OctetString;"Hello
world!
Sample 1"
1.3.6.100.1.2.3.4.1.3.1.21;OctetString;"Bye Sample 1"

$ cat sample_oids/two.csv
1.3.6.100.1.2.3.5.1.1.0;OctetString;Sample2
1.3.6.100.1.2.3.5.1.2.0;Integer;42
1.3.6.100.1.2.3.5.1.3.0;IPAddress;10.20.30.40
1.3.6.100.1.2.3.5.1.4.0;Counter32;11223344
1.3.6.100.1.2.3.5.1.5.0;Gauge32;88866
1.3.6.100.1.2.3.5.1.6.0;TimeTicks;1876300
1.3.6.100.1.2.3.5.1.7.0;Counter64;111
1.3.6.100.1.2.3.5.1.8.0;ObjectIdentifier;1.3.6.8.9.999.100

$ sudo docker run -t --rm -e BASE_OID=1.3.6.100.1.2.3 -e SNMP_COMMUNITY=my_comm -v $(pwd)/sample_oids:/app/oids elpadrinoiv/snmp_mock
Replacing snmp community with my_comm
Replacing base_oid with 1.3.6.100.1.2.3
snmpd[9]: Created directory: /var/lib/net-snmp/mib_indexes
snmpd[9]: Turning on AgentX master support.
snmpd[11]: NET-SNMP version 5.7.3
2017-06-02 02:24:33.444 INFO Running snmp mock
2017-06-02 02:24:33.462 INFO Loading oids...
2017-06-02 02:24:33.462 INFO Loading oids from file /app/oids/_two.csv
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.1.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.2.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.3.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.4.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.5.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.6.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.7.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.5.1.8.0
2017-06-02 02:24:33.462 INFO Loading oids from file /app/oids/_one.csv
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.4.1.1.0
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.4.1.3.1.20
2017-06-02 02:24:33.462 INFO Adding oid 1.3.6.100.1.2.3.4.1.3.1.21
```

Once the mock server is running, you can try snmpget/walk:
```
$ snmpget -v2c -c my_comm 172.17.0.4 1.3.6.100.1.2.3.5.1.2.0
iso.3.6.100.1.2.3.5.1.2.0 = INTEGER: 42

$ snmpwalk -v2c -c my_comm 172.17.0.4 1.3.6.100.1.2.3.5
iso.3.6.100.1.2.3.5.1.1.0 = STRING: "Sample2"
iso.3.6.100.1.2.3.5.1.2.0 = INTEGER: 42
iso.3.6.100.1.2.3.5.1.3.0 = IpAddress: 10.20.30.40
iso.3.6.100.1.2.3.5.1.4.0 = Counter32: 11223344
iso.3.6.100.1.2.3.5.1.5.0 = Gauge32: 88866
iso.3.6.100.1.2.3.5.1.6.0 = Timeticks: (1876300) 5:12:43.00
iso.3.6.100.1.2.3.5.1.7.0 = Counter64: 111
iso.3.6.100.1.2.3.5.1.8.0 = OID: iso.3.6.8.9.999.100
iso.3.6.100.1.2.3.5.1.8.0 = No more variables left in this MIB View (It is past the end of the MIB tree)
```

## docker-compose
```
version: 2

service:
    snmp_mock:
        image: elpadrinoiv/snmp_mock:latest
        environment:
            SNMP_PORT=161
            SNMP_COMMUNITY=public
            BASE_OID=1.3.6
            APP_PORT=8080

        volumes:
            ./my_oid_responses:/app/oids
```

