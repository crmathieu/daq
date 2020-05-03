package data

// header is 16 bytes 
const PACKET_PAYLOAD_LENGTH = 32 //255
const PACKET_LENGTH = PACKET_PAYLOAD_LENGTH + PACKET_PAYLOAD_OFFSET
const PACKET_START = byte(0xff)

// packet offsets
const PACKET_START_OFFSET = 0		// start marker
const PACKET_CRC_OFFSET = 2			// CRC is 32bits and calculated on payload only
const PACKET_NDP_OFFSET = 6			// number of datapoints in this packet
const PACKET_TT_OFFSET = 7			// timestamp is on 64bits
const PACKET_RES_OFFSET = 15		// 1 reserved bytes
const PACKET_PAYLOAD_OFFSET = 16	// payload starts here (15 dp per payload)

const PACKET_HEADER = 16
const DATAPOINT_SIZE = 14

/*type GSbuf struct {
	Marker byte
	Index byte
	Buffer [PACKET_PAYLOAD_LENGTH]byte
	Crc uint32
	//	Ready bool
}*/

//type PL_DYN struct {
//    CoorX, CoorY, CoorZ float64
//    VelX, VelY, VelZ    float64/
//}

//type GSbuf []byte

const(
	// instruments offset in sensors map
 	SVELOCITY = 0 
    SPOSITION = 1
    STURBOPUMP = 2
	SENGINEPRE = 3
	SMASSPROPELLANT = 4

	INSTRUMENTS_COUNT = SMASSPROPELLANT + 1

	// Rocket constants
	//DRYWEIGHT = 15000
	//MAXVOL_OXYDIZER = 100000
	//MAXVOL_PROPELLANT = 200000
	
	DOWNLINK_PORT = ":2000"
)   


// each dp is 14 bytes long
type SENSvelocity struct {
	Id    	 		uint16
	Velocity 		float32
	Acceleration 	float32
	Reserved 		[4]byte
}

type SENSposition struct {
	Id    		uint16
	Range 		float32
    Inclinaison float32
    Altitude 	float32
}

type SENSturboPump struct {
	Id   uint16
	Rpm  int32
	Reserved [8]byte
}

type SENSenginePressure struct {
	Id   	 uint16
	Pressure float32
	Reserved [8]byte
}

type SENSpropellantMass struct {
	Id   	 uint16
	Mass   	 float32
	Reserved [8]byte
}

type Pempty []byte
