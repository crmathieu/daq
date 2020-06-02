package data


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
 	SVELOCITY_OFFSET = uint16(0) 
    SPOSITION_OFFSET = uint16(1)
//    STURBOPUMP = 2
//	SENGINEPRE = 3
	STILTANGLE_OFFSET = uint16(2)
	STHRUST_OFFSET = uint16(3)
	SEVENT_OFFSET = uint16(4)
	SMASSPROPELLANT_OFFSET = uint16(5)
	STIME_OFFSET = uint16(6)

	INSTRUMENTS_COUNT = STIME_OFFSET + 1

	// Rocket constants
	//DRYWEIGHT = 15000
	//MAXVOL_OXYDIZER = 100000
	//MAXVOL_PROPELLANT = 200000
	
	DOWNLINK_SERVER = "localhost:2000"

	IDVELOCITY = uint32(1)
    IDPOSITION = uint32(2)
//    STURBOPUMP = 2
//	SENGINEPRE = 3
	IDTILTANGLE = uint32(3)
	IDTHRUST = uint32(4)
	IDEVENT = uint32(5)
	IDMASSPROPELLANT = uint32(6)
	IDTIME = uint32(7)
)   



// each dp is 16 bytes long
type SENSvelocity struct {
	Id    	 		uint32 //uint16
	Velocity 		float32
	Acceleration 	float32
	Stage			uint32
//	Reserved 		[4]byte //[6]byte
}

type SENSevent struct {
	Id    		uint32
	EventId 	uint32
	Time 		float32
	EventMap 	uint32
}

type SENStime struct {
	Id    		uint32
	Time 		float32
	reserved 	[8]byte
}

type SENStimestamp struct {
	Id    		uint64
	TT 			float64
}

type SENSposition struct {
	Id    		uint32 //uint16
	Range 		float32
    Altitude 	float32
	//Inclinaison float32
	Stage		uint32
	//reserved	[2]byte
}

type SENStiltAngle struct {
	Id   uint32 //uint16
	Alpha, Beta, Gamma float32
/*	Angle  float32
	RateOfChange float32
	Reserved [4]byte //[6]byte */
}

type SENSthrust struct {
	Id   	uint32 //uint16
	Thrust  float32
	Stage   int32 //int8
	Reserved [4]byte //[9]byte
}

type SENSturboPump struct {
	Id   uint32 //uint16
	Rpm  int32
	Reserved [8]byte //[10]byte
}

type SENSenginePressure struct {
	Id   	 uint32 //uint16
	Pressure float32
	Reserved [8]byte //[10]byte
}

type SENSpropellantMass struct {
	Id   	 uint32 //uint16
	Mass   	 float32
	//Mejected float32
	Stage	 uint32
	Mflow	 float32
	//reserved [2]byte
}

type DataPoint struct {
	Id    	 		uint32 //uint16
	Reserved 		[12]byte //[14]byte
}

type Pempty []byte

type CONFinfo struct {
//		DB 				  *sql.DB
//		DB_dsn     		  string `yaml:"db_dsn"`	
		REDIS_dsn  		  string `yaml:"redis_dsn"`	
		RedisHost  string `yaml:"cache-redis-host"`
		RedisPort  string `yaml:"cache-redis-port"`
		RedisPass  string `yaml:"cache-redis-pass"`
		RedisValid bool
}

var CInfo CONFinfo 