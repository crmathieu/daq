package data

/*
Packet format:
0 ---------------> STARTMARKER 1 byte (0xAf)
1 ---------------> CRC 4 bytes (int32)
6 ---------------> NDP 1 byte - Number of DataPoints in this packet
7 ---------------> Time Stamp 8 bytes (int64)
15 --------------> reserved 1 byte
16 --------------> Start of payload

- Data between offset 0 to 15 pertains to the packet header
- Each datapoint measurment consist of a 14 bytes buffer
*/

// header is 16 bytes 
const PACKET_PAYLOAD_LENGTH = 32 //255
const PACKET_LENGTH = PACKET_PAYLOAD_LENGTH + PACKET_PAYLOAD_OFFSET
const PACKET_START1 = byte(0xFF)
const PACKET_START2 = byte(0xAA)

// packet offsets
const PACKET_START_OFFSET = 0		// start marker
const PACKET_CRC_OFFSET = 2			// CRC is 32bits and calculated on payload only
const PACKET_NDP_OFFSET = 6			// number of datapoints in this packet
const PACKET_TT_OFFSET = 7			// timestamp is on 64bits
const PACKET_RES_OFFSET = 15		// 1 reserved bytes
const PACKET_PAYLOAD_OFFSET = 16	// payload starts here (15 dp per payload)

const PACKET_HEADER = 16
const DATAPOINT_SIZE = 16

/*
0	255 
1	170 	STARTER
	---
2	143 
3	96 		CRC
4	145 
5	196
	--- 
6	3 		NDP
	---
7	64 
8	106 
9	21 
10	96 		TST
11	156 
12	99 
13	11 
14	22 
	---
15	0 		reserved
	---
16	... the datapoints ...
*/