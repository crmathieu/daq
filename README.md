# daq
Daq is an attempt to simulate the flow of data between a moving vehicle (it could be a rocket, a car or anything else) and a ground station. The ground station stream the data to clients that connect to it using an authentication token.


### Vehicle
In this simulation, we use the dynamics of a rocket launch to generate the data. This data is sent to the ground station that has any number of client connected to it. Clients may use a portion (or the whole set) of the data.


### Ground station
The ground station has 3 essential functions:
- receive data from the vehicle and compare CRC32 calculated and CRC32 value provided in packet header to make sure no error were introduced in the transmission. 
- Place the set of datapoints received to its streaming queue.
- Accept connection requests from clients to allow them to access streaming data


# Simulate the whole thing


### First, start the ground station
Go to daq/cmd/groundstation and type:
```bash
> go build
```
and then
```bash
> ./groundstation
```


### Second, run a client against the ground station
An example of client is provided as an html file. From the address bar of your favorite browser, type:
```
localhost:1969/stream/123
```
- *1969* is the default web port from which you can request a connection to the ground station process. 
- *123* is the authentication token you must provide to access the service.
You should get a page showing you not much, as the vehicle hasn't been launched yet.

### Third, launch the rocket
Go to daq/cmd/downlink/launch and type:
```bash
> go build
```
and then
```bash
> ./launch
```

From now on, you should see data coming to your client.


### data point
Each data point consists of a 16 bytes buffer containing a datapoint Id and, depending on the datapoint, a combination of values held on 4 bytes (int32, uint32, float32).  

### data packet
A data packet is a set of datapoints put together and send to the ground station as a whole. A packet has a header that contains:

```text
offset  0:  a start marker
offset  2:  the 32bits CRC calculated on the payload only
offset  6:  the number of datapoints in this packet
offset  7:  the timestamp on 64bits
offset 15:  1 reserved bytes
```
