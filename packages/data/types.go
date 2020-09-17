package data

const (
	// instruments offset in sensors map
	SVELOCITY_OFFSET       = uint16(0)
	SPOSITION_OFFSET       = uint16(1)
	SANGLES_OFFSET         = uint16(2)
	STHRUST_OFFSET         = uint16(3)
	SEVENT_OFFSET          = uint16(4)
	SMASSPROPELLANT_OFFSET = uint16(5)
	STIME_OFFSET           = uint16(6)

	INSTRUMENTS_COUNT = STIME_OFFSET + 1

	// Rocket constants

	// groundstation default server/port/webport
	DOWNLINK_SERVER  = "localhost"
	DOWNLINK_PORT    = "2000"
	DOWNLINK_WEBPORT = "1969"

	// datapoints ID
	IDVELOCITY       = uint32(1)
	IDPOSITION       = uint32(2)
	IDANGLES      	 = uint32(3)
	IDTHRUST         = uint32(4)
	IDEVENT          = uint32(5)
	IDMASSPROPELLANT = uint32(6)
	IDTIME           = uint32(7)
)

// Datapoint definitions (each dp is 16 bytes long)
type SENSvelocity struct {
	Id           uint32
	Velocity     float32
	Acceleration float32
//	Stage        uint32
	reserved     [4]byte
}

type SENSevent struct {
	Id       uint32
	EventId  uint32
	Time     float32
	EventMap uint32
}

type SENStime struct {
	Id       uint32
	Time     float32
	reserved [8]byte
}

type SENStimestamp struct {
	Id uint64
	TT float64
}

type SENSposition struct {
	Id       uint32
	Range    float32
	Altitude float32
//	Stage    uint32
	reserved     [4]byte
}

type SENSangles struct {
	Id      uint32
	Alpha	float32 
	Beta	float32 
	Gamma 	float32
}

type SENSthrust struct {
	Id       uint32
	Thrust   float32
//	Stage    int32
	Reserved [8]byte
}

type SENSturboPump struct {
	Id       uint32
	Rpm      int32
	Reserved [8]byte
}

type SENSenginePressure struct {
	Id       uint32
	Pressure float32
	Reserved [8]byte
}

type SENSpropellantMass struct {
	Id    uint32
	Mass  float32
//	Stage uint32
	Mflow float32
	reserved     [4]byte
}

type DataPoint struct {
	Id       uint32
	Reserved [12]byte
}

// redis extension (if needed for various reason)
type CONFinfo struct {
	REDIS_dsn  string `yaml:"redis_dsn"`
	RedisHost  string `yaml:"cache-redis-host"`
	RedisPort  string `yaml:"cache-redis-port"`
	RedisPass  string `yaml:"cache-redis-pass"`
	RedisValid bool
}

var CInfo CONFinfo
