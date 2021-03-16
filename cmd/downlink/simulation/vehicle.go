package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"github.com/crmathieu/daq/packages/data"
	"io/ioutil"
	"os"
	"math"
	"unsafe"
)

type SensorHandlers struct {
	ReadSensor 	func (int32) [data.DATAPOINT_SIZE]byte
}
func NullSensor (a int32) [data.DATAPOINT_SIZE]byte {return [data.DATAPOINT_SIZE]byte{}}

type Profile struct {
	Description 	string	`yaml:"description"`
	EarthRotation 	float64 `yaml:"earthrotation"`
	LaunchAzimuth	float64 `yaml:"launchazimuth"`
	LaunchLatitude	float64 `yaml:"launchlatitude"`
	BurnoutTime		float64 `yaml:"burnout"`
	OrbitPerigee    float64 `yaml:"orbitperigee"`
	Events 			[]Event	`yaml:"events"`
}

type Event struct {
	Id 		string	`yaml:"id"`
	T		float64	`yaml:"time"`
	Stage	int8	`yaml:"stage"`
	Rate    int8    `yaml:"rate,omitempty"`
} 

var profile Profile
var M_PI = math.Pi

const (

	G 	 = 6.67384e-11	// Gravitational Const
	Me 	 = 5.97219e24	// Mass of Earth
	Re 	 = 6378137		// Radius of Earth
	g0 	 = 9.807 //9.7976		// Gravity acceleration on surface
	tinc = 0.01 		// time increment
)

type VEHICLE struct {
	//Clock	float64
	SysGuidance Guidance
	Instruments		[]unsafe.Pointer 
	Handlers		[]SensorHandlers
	EventsMap		uint32
	LastEvent		uint32
	Stages []RocketStage
}

type RocketStage struct {
	Clock	float64		// stage reference clock
	dt 		float64		// time increment
	gt      float64		// gravity turn time increment

	// drag parameters
	Cd		float64 // Coeff. of drag
	CSArea  float64	// cross-sectional area
	
	// mass
	Mr	float64 // Dry mass  
	Mf  float64 // Fuel mass
	Mp 	float64	// Payload mass

	// plumbing
	EngineID		string
	RunningEngines	int32

	// guidance
	ThrottleRate	float64
	Thrust			float64

	// Force x,y on stage
	ForceX  		float64
	ForceY  		float64

	// angles: 
	// - beta is the angle through which gravity pulls the vehicle. 
	// - alpha is the angle of attack relative to earth.

	alpha 			float64 	// alpha = angle of velocity
	beta 			float64		// beta = angle with gravity 
	gamma			float64		// gamma = angle of thrust

	// velocity
	VRelative 		float64		// Relative Velocity
	VAbsolute		float64		// absolute velocity

	// Mass
	Mass 			float64	

	// polar variables
	Acc, Force, RVel, AVel float64	
	altitude,drange float64

	// Distance to focus (earth center)
	DTF 			float64
		
	// cartesian coordinates
	cx, cy			float64

	// cartesian Absolute and relative velocity (rel. to Earth)
	vAx, vAy 		float64 
	vRx, vRy 		float64

	// cartesian acceleration
	ax, ay 			float64

	apogee, perigee float64

}

type Engine struct {
	Isp_sl 				float64 	// Sea level ISP
	Isp_vac 			float64		// Vacuum Isp
	Th_sl				float64 	// Sea level Thrust
	Th_vac 				float64		// Vacuum Thrust
	Min_ThrottleRate 	float64  
	Flow_rate 			float64     // fuel flow rate (kg/s)
}

var aerodynPressure float64 				// aero pressure
var drag float64				// drag
var dm float64 				// rate of fuel consumption

//var t = float64(-10.0)		// time (initialized at -10sec)
//var dt = float64(0.001)		// time step

var vE = float64(0.0)		// inclination (rads) = 28.49*M_PI/180;
// vE = atoi(optarg)==0 ? 0 : 407.6614278; break; // Earth velocity at Cape Canaveral

/*
var Veh = VEHICLE {
	//Clock: -10.0,
	Stages: []RocketStage{
		// booster
		{	Clock:-10.0,			// clock is set at -10sec before launch
			dt: tinc, //0.001,				// time increment
			Cd:0.3, 				// drag coefficient 
			CSArea:10.52, 			// cross section area in m*m
			Mr:20000, 
			Mf:390000, 
			Mp:0,
			RunningEngines:0,
			EngineID: "M1D",
			ThrottleRate: 0.0,
			Thrust: 0.0,
			ForceX: 0.0,
			ForceY: 0.0,
			alpha: 0.0,
			beta: M_PI/2,
			gamma: M_PI/2,
		}, 
		// stage2
		{	Clock: -10.0,
			dt: tinc, //0.001,
			Cd:0.3, 		
			CSArea:10.52, 
			Mr:4900, 
			Mf:75700, 
			Mp:1200,
			RunningEngines:0,
			EngineID: "M1Dv",
			ThrottleRate: 0.0,
			Thrust: 0.0,
			ForceX: 0.0,
			ForceY: 0.0,
			alpha: 0.0,
			beta: 0.0,
			gamma: 0.0,
		},},
}
*/
var EnginesMap = map[string]Engine{
	"M1D": Engine{
		Isp_sl:282, 
		Isp_vac:311, 
		Th_sl: 650000, 
		Th_vac: 720000,
		Min_ThrottleRate: 0.40,
		Flow_rate: 235.4,
	},
	"M1Dv": Engine{
		Isp_sl:0, 
		Isp_vac:345, 
		Th_sl:0, 
		Th_vac: 801000,
		Min_ThrottleRate: 0.39,
		Flow_rate: 235.4,
	},
	"M1DB5": Engine{
		Isp_sl:282, 
		Isp_vac:311, 
		Th_sl: 854000, // from wikipedia --650000, 
		Th_vac: 981000, // from wikipedia --720000,
		Min_ThrottleRate: 0.40,
		Flow_rate: 235.4, // unsure

	},
	"M1DvB5": Engine{
		Isp_sl:0, 
		Isp_vac:348, // from wikipedia --345, 
		Th_sl:0, 
		Th_vac:948000, // from wikipedia --801000,
		Min_ThrottleRate: 0.39,
		Flow_rate: 235.4, // unsure
	},
} 

const BOOSTER = 0
const STAGE2 = 1
const STAGE3 = 2

type Guidance struct {
	_release, _pitch, _stagesep, _MEI1, _MEI2, _MEI3, _SEI1, _LBURN, _BBURN, _MECO1, _SECO1, _MECO2, _MECO3, _EBURN bool
}


//var _release, _pitch, _MEI1, _MEI2, _MEI3, _SEI1, _LBURN, _BBURN, _MECO1, _SECO1, _MECO2, _MECO3 bool = false, false, false, false, false, false, false,false, false, false, false, false

// InitGuidance ---------------------------------------------------------------
// reads the flight profile to build the event table
// ----------------------------------------------------------------------------
func (r *VEHICLE) InitGuidance() *[]Event {

	r.initInstruments()

//	s[0][1] = Re;
	r.Stages[BOOSTER].cy = Re
	r.Stages[BOOSTER].DTF = Re
	r.Stages[BOOSTER].AVel = vE
	r.Stages[BOOSTER].RVel = 0

	r.Stages[STAGE2].Mass = r.Stages[STAGE2].Mr + r.Stages[STAGE2].Mf + r.Stages[STAGE2].Mp
	r.Stages[BOOSTER].Mass = r.Stages[BOOSTER].Mr + r.Stages[BOOSTER].Mf + r.Stages[STAGE2].Mass

	file := "profile.yml"
	filepath := "./profiles/" + file
	if _, err := os.Stat(filepath); err == nil {
		// file exists
		data, err := ioutil.ReadFile(filepath)
		fmt.Printf("\n---\n%s\n---\n", data)
		if err == nil {
			err = yaml.Unmarshal(data, &profile)
			if err == nil {
				fmt.Printf("Profile read for: %v\nlatitude:%v\nazimuth:%v\nburnout time: %v\ntarget orbit: %v km\nEvents: %v\n", profile.Description, profile.LaunchLatitude, profile.LaunchAzimuth, profile.BurnoutTime, profile.OrbitPerigee*1e-3, profile.Events)
				return &profile.Events
			} 	
			fmt.Println("Error unmarshalling flight profile: ", err.Error())
		} else {
			fmt.Println("Error reading flight profile:", err.Error())
		}
	} else {
		fmt.Println("Flight Profile ''" + file + "' does not exist")
	}
	return nil
}


func (v *VEHICLE) initInstruments() {

	v.Instruments = make([]unsafe.Pointer, data.INSTRUMENTS_COUNT)
	v.Handlers	  = make([]SensorHandlers, data.INSTRUMENTS_COUNT)

	v.Instruments[data.SVELOCITY_OFFSET]		= (unsafe.Pointer)(&data.SENSvelocity{		Id:data.IDVELOCITY, 	Velocity:0.0, Acceleration:0.0, })
	v.Instruments[data.SPOSITION_OFFSET]		= (unsafe.Pointer)(&data.SENSposition{		Id:data.IDPOSITION, 	Range:0.0,    Altitude:0.0,})
	v.Instruments[data.SEVENT_OFFSET]			= (unsafe.Pointer)(&data.SENSevent{			Id:data.IDEVENT, 		EventId:0,    Time:0, })
	v.Instruments[data.STIME_OFFSET]			= (unsafe.Pointer)(&data.SENStime{			Id:data.IDTIME, 		Time:0, })
	v.Instruments[data.SANGLES_OFFSET]			= (unsafe.Pointer)(&data.SENSangles{		Id:data.IDANGLES, 		Alpha:0.0, Beta:0.0, Gamma:0.0,})
	//v.Instruments[data.STHRUST_OFFSET]			= (unsafe.Pointer)(&data.SENSthrust{		Id:data.IDTHRUST, Thrust:0,})
	//v.Instruments[data.SMASSPROPELLANT_OFFSET]	= (unsafe.Pointer)(&data.SENSpropellantMass{Id:data.IDMASSPROPELLANT, Mflow: 0.0, Mass: 0.0,})

	v.Handlers[data.SVELOCITY_OFFSET]		= SensorHandlers{ReadSensor: v.readVelocity,}
	v.Handlers[data.SPOSITION_OFFSET]		= SensorHandlers{ReadSensor: v.readPosition, }
	v.Handlers[data.SEVENT_OFFSET]			= SensorHandlers{ReadSensor: v.readEvent,}
	v.Handlers[data.STIME_OFFSET]			= SensorHandlers{ReadSensor: v.readTime,}
	v.Handlers[data.SANGLES_OFFSET]			= SensorHandlers{ReadSensor: v.readAngles,} //v.readTiltAngle,}
	v.Handlers[data.STHRUST_OFFSET]			= SensorHandlers{ReadSensor: nil,} //v.readThrust,}
	v.Handlers[data.SMASSPROPELLANT_OFFSET]	= SensorHandlers{ReadSensor: nil,} //v.readPropellantMass, }

}

func NewVehicle() *VEHICLE {
	return &VEHICLE {
	//Clock: -10.0,
	Stages: []RocketStage{
		// booster
		{	Clock:-10.0,			// clock is set at -10sec before launch
			dt: tinc, //0.001,				// time increment
			Cd:0.3, 				// drag coefficient 
			CSArea:10.52, 			// cross section area in m*m
			Mr:20000, 
			Mf:390000, 
			Mp:0,
			RunningEngines:0,
			EngineID: "M1DB5",
			ThrottleRate: 1.0,
			Thrust: 0.0,
			ForceX: 0.0,
			ForceY: 0.0,
			alpha: 0.0,
			//beta: M_PI/2,
			gamma: M_PI/2,
		}, 
		// stage2
		{	Clock: -10.0,
			dt: tinc, //0.001,
			Cd:0.3, 		
			CSArea:10.52, 
			Mr:4900, 
			Mf:75700, 
			Mp:1200,
			RunningEngines:0,
			EngineID: "M1DvB5",
			ThrottleRate: 1.0,
			Thrust: 0.0,
			ForceX: 0.0,
			ForceY: 0.0,
			alpha: 0.0,
			//beta: M_PI/2,
			gamma: M_PI/2,
		},},
	}
}

func (v *VEHICLE) NoFuel(stage int) bool {
	if v.Stages[stage].Mf < 5 {
		return true
	}
	return false
}