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
 	SVELOCITY = uint16(0) 
    SPOSITION = uint16(1)
//    STURBOPUMP = 2
//	SENGINEPRE = 3
	STILTANGLE = uint16(2)
	STHRUST = uint16(3)
	SMASSPROPELLANT = uint16(4)

	INSTRUMENTS_COUNT = SMASSPROPELLANT + 1

	// Rocket constants
	//DRYWEIGHT = 15000
	//MAXVOL_OXYDIZER = 100000
	//MAXVOL_PROPELLANT = 200000
	
	DOWNLINK_SERVER = "localhost:2000"
)   



// each dp is 16 bytes long
type SENSvelocity struct {
	Id    	 		uint16
	Velocity 		float32
	Acceleration 	float32
	Reserved 		[6]byte
}

type SENSposition struct {
	Id    		uint16
	Range 		float32
    Altitude 	float32
	Inclinaison float32
	reserved	[2]byte
}

type SENStiltAngle struct {
	Id   uint16
	Angle  float32
	RateOfChange float32
	Reserved [10]byte
}

type SENSthrust struct {
	Id   	uint16
	Thrust  float32
	Stage   int8
	Reserved [9]byte
}

type SENSturboPump struct {
	Id   uint16
	Rpm  int32
	Reserved [10]byte
}

type SENSenginePressure struct {
	Id   	 uint16
	Pressure float32
	Reserved [10]byte
}

type SENSpropellantMass struct {
	Id   	 uint16
	Mass   	 float32
	Mejected float32
	Mflow	 float32
	reserved [2]byte
}

type SENSgeneric struct {
	Id    	 		uint16
	Reserved 		[14]byte
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