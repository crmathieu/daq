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
 	PVELOCITY = 0 
    PCOORDINATES = 1
    PTURBOPUMP = 2
    PENGINEPRE = 3
)   

// each dp is 14 bytes long
type Pvelocity struct {
	Id   uint16
	Velx float32
    Vely float32
    Velz float32
}

type Pcoordinates struct {
	Id    uint16
	Coorx float32
    Coory float32
    Coorz float32
}

type PturboPumpRPM struct {
	Id   uint16
	Rpm  int32
	Reserved [8]byte
}

type PenginePressure struct {
	Id   	 uint16
	Pressure float32
	Reserved [8]byte
}

type Pempty []byte
