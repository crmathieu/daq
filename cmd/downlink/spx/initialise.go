package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"math"
)

type Profile struct {
	Description 	string	`yaml:"description"`
	EarthRotation 	float64 `yaml:"earthrotation"`
	Events 			[]Event	`yaml:"events"`
}

type Event struct {
	Id 		string	`yaml:"id"`
	T		float64	`yaml:"time"`
	Stage	int8	`yaml:"stage"`
} 

var profile Profile
var M_PI = math.Pi

const (

	G 	= 6.67384e-11	// Gravitational Const
	Me 	= 5.97219e24	// Mass of Earth
	Re 	= 6378137		// Radius of Earth
	g0 	= 9.7976		// Gravity acceleration on surface
)

type Rocket struct {
	//Clock	float64
	Stages []RocketStage
}

type RocketStage struct {
	Clock	float64		// stage reference clock
	dt 		float64		// time increment

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
	gam				float64		// gamma = angle of thrust

	// velocity
	VRelative 		float64		// Relative Velocity
	VAbsolute		float64		// absolute velocity

	// Mass
	Mass 			float64	

	// polar acceleration
	Acc 			float64	

	// polar distance from earth center
	PolarDistance 	float64
		
	// cartesian coordinates
	cx, cy			float64

	// cartesian Absolute and relative velocity (rel. to Earth)
	vAx, vAy 		float64 
	vRx, vRy 		float64

	// cartesian acceleration
	ax, ay 			float64
}

type Engine struct {
	Isp_sl 	float64 	// Sea level ISP
	Isp_vac float64		// Vacuum Isp
	Th_sl	float64 	// Sea level Thrust
	Th_vac 	float64		// Vacuum Thrust
}

var aerodynPressure float64 				// aero pressure
var drag float64				// drag
var dm float64 				// rate of fuel consumption
//var t = float64(-10.0)		// time (initialized at -10sec)
//var dt = float64(0.001)		// time step

var vE = float64(0.0)		// inclination (rads) = 28.49*M_PI/180;
// vE = atoi(optarg)==0 ? 0 : 407.6614278; break; // Earth velocity at Cape Canaveral

var F9 = Rocket {
	//Clock: -10.0,
	Stages: []RocketStage{
		// booster
		{	Clock:-10.0,			// clock is set at -10sec before launch
			dt: 0.001,				// time increment
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
			gam: M_PI/2,
		}, 
		// stage2
		{	Clock: 0.0,
			dt: 0.001,
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
			gam: 0.0,
		},},
}

var EnginesMap = map[string]Engine{
	"M1D": Engine{
		Isp_sl:282, 
		Isp_vac:311, 
		Th_sl: 650000, 
		Th_vac: 720000,
	},
	"M1Dv": Engine{
		Isp_sl:0, 
		Isp_vac:345, 
		Th_sl:0, 
		Th_vac:801000,
	},
} 

const BOOSTER = 0
const STAGE2 = 1
const STAGE3 = 2

var _release, _pitch, _MEI1, _MEI2, _MEI3, _SEI1, _LBURN, _BBURN, _MECO1, _SECO1, _MECO2, _MECO3 bool = false, false, false, false, false, false, false,false, false, false, false, false

// Init -----------------------------------------------------------------------
// reads the flight profile to build the event table
// ----------------------------------------------------------------------------
func Init() *[]Event {

//	s[0][1] = Re;
	F9.Stages[BOOSTER].cy = Re
	F9.Stages[BOOSTER].PolarDistance = Re;

	F9.Stages[BOOSTER].Mass = F9.Stages[BOOSTER].Mr + F9.Stages[STAGE2].Mr + F9.Stages[BOOSTER].Mf + F9.Stages[STAGE2].Mf + F9.Stages[STAGE2].Mp;
	F9.Stages[STAGE2].Mass = F9.Stages[BOOSTER].Mr + F9.Stages[STAGE2].Mr + F9.Stages[BOOSTER].Mf + F9.Stages[STAGE2].Mf + F9.Stages[STAGE2].Mp;

	file := "profile.yml"
	filepath := "./profiles/" + file
	if _, err := os.Stat(filepath); err == nil {
		// file exists
		data, err := ioutil.ReadFile(filepath)
		fmt.Printf("\n---\n%s\n---\n", data)
		if err == nil {
			err = yaml.Unmarshal(data, &profile)
			if err == nil {
				fmt.Printf("Profile read for: %v\n%v\n", profile.Description, profile.Events)
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


