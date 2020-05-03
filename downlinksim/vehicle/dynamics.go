package main

// preliminary physics formulas:
// in circular motion 
// ------------------
// speed		v = r.w  w = d(alpha)/dt where alpha is angle created by vehicle between t -> t+dt
// 
// acceleration	
//				a = dv/dt = d(r.w)/dt = r.dw/dt + w.dr/dt
//					dr/dt = v => a = r.dw/dt + w.v = r.dw/dt + w.(r.w) = r.dw/dt + r.w*w
// 				if uniform motion: dw/dt = 0 hence
//					a = r.w*w is the acceleration due to change in the direction for a uniform motion
//					a = r.w*w + r.dw/dt  for non-uniform motion
// Centripetal force
//				F = m.a = m.(r.w*w) = (m.r*r.w*w)/r = m.v*v/r 

/*
 * Equations to determine the vehicle position, speed and direction
 * 
 *	dv/dt = (T - D)/m - g.sin(gamma)
 *  
 *	v.d(gamma)/dt = - (g - v**2/(r0+h)).cos(gamma)
 *
 *	dx/dt = v.cos(gamma).r0/(r0+h)
 *
 *	dh/dt = v.sin(gamma)
 *
 * where:
 *	gamma: angle between local horizontal and rocket axis
 *	r0: earth radius
 *	h: altitude
 *	v: speed
 *  D: drag
 * 	T: thrust
 *
 * 		The drag D = q.A.CD  where 	q  = (rho.v**2)/2  (rho = air density at this altitude) Aerodynamic pressure
 *									A  = frontal surface of vehicle
 *									CD = coefficient linked to vehicle geometry 
 *
 *		various shape's drag coefficient CD:
 *		-----------------------------------
 *		- Sphere: 			0.47
 *		- 1/2 sphere: 		0.42
 *		- Cone:				0.50
 *		- Cube:				1.05
 *		- Angled cube:  	0.80 
 *		- Long cylinder:	0.82
 *		- Short cylinder:	1.15
 *		- Wing:				0.04
 *		- 1/2 Wing:			0.09
 *
 *		The Thrust T = dm/dt * Ve + (pe - pa) * Ae    	where  	dm/dt 		= rate of ejected mass flow
 *																Ve 			= exhaust gas ejection speed 
 *																pe			= exhaust gas pressure
 *																pa			= atmospheric pressure
 *																Ae			= area of exit
 *		The massflow mdot = - T/(Isp * GRAVITY_ACC)		where   GRAVITY_ACC = gravitational force at sea level
*/

import(
	"github.com/crmathieu/daq/data"
	"time"
	"math"
	"fmt"
//	"unsafe"
)	
const (
//	GRAVITY_ACC = 9.80665
	CD = 0.42
	//EARTHRADIUS = float64(6356766.0)
	RHO = 1.225	
	H0 = 7500
)
var vdot, mdot, m, v, F, D, beta, betadot, xrange, hrange float32


var zvelocity = float32(0.2)
var altitude = float32(1.1)
var tbrpm = int32(15000)
var engpressure = float32(150.56)
var propmass = float32(F9_S1_PropellantMass)

// placeholder instruments updates
func (v *VEHICLE) setPosition() {
	/*altitude++
	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude = altitude*/
}
func (v *VEHICLE) setTurboPumpRPM() {
	tbrpm++
	(*(*data.SENSturboPump)(v.Stage[v.CurrentStage].Instruments[data.STURBOPUMP])).Rpm = tbrpm
}
func (v *VEHICLE) setEnginePressure() {
	engpressure++
	(*(*data.SENSenginePressure)(v.Stage[v.CurrentStage].Instruments[data.SENGINEPRE])).Pressure = engpressure
}
func (v *VEHICLE) setPropellantMass() {
	/*propmass--
	(*(*data.SENSpropellantMass)(v.Stage[v.CurrentStage].Instruments[data.SMASSPROPELLANT])).Mass = propmass*/
}

var UPDATE_TICK = 1000 * time.Millisecond

// RunInstrumentsUpdate -------------------------------------------------------
// calls the SetSensor function defined in the datapoint structure in an
// infinite loop, only interrupted by a sleep function
// ----------------------------------------------------------------------------
func (v *VEHICLE) RunInstrumentsUpdate() {
	for {
		time.Sleep(UPDATE_TICK)
		for _, i := range v.Stage[v.CurrentStage].Handlers {
			i.SetSensor()
			//fmt.Println(k, "--->", v.ReadSensor())
		}
	}
}

func (v *VEHICLE) launch() {
	var meco = false
	var mecoTimer = 3 // seconds
	liftoff(v)
	//return
	for {
		time.Sleep(UPDATE_TICK)
		v.Clock++
		if meco == true {
			mecoTimer--
			fmt.Println("Meco wait", mecoTimer)
			if mecoTimer <= 0 {
				meco = false
				if v.CurrentStage + 1 < int8(len(v.Stage)) {
					v.CurrentStage++
				} else {
					fmt.Println("Shutting down!")
					break
				}
			}
		}
		meco = v.setRates(meco)
		v.updateRocketSensors()
		fmt.Println(v.Clock, "--> Vel:",v.Velocity, "m/s -- ", "Alt:",v.Altitude/1000,	"km -- Downrange:", v.Range/1000,"km", "Gamma=",v.Gamma*rad,"deg")
		//fmt.Println(v.Clock, "-->",(v.Velocity)*3.6, "k/h -- ", "Alt:",v.Altitude/1000,	"km, Downrange:", (*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range/1000,"km, Gamma=",v.Gamma)

	}	
}


func (v *VEHICLE) setFrontalArea() {
//	v.FrontalArea = 4 * math.Pi * math.Pow(float64(v.Stage[v.CurrentStage].Diameter/2), 2)
	v.FrontalArea = math.Pi * math.Pow(float64(v.Stage[v.CurrentStage].Diameter/2), 2)/4
}

// getMdot --------------------------------------------------------------------
// calculate massflow based on thrust and ISP: 
//		- T/(Isp * GRAVITY_ACC) 
// Note that the number is a negative mass being consumed
// ---------------------------------------------------------------------------- 
func (v *VEHICLE) getMdot(throttle float32) float32 {
	return throttle * v.Stage[v.CurrentStage].Thrust/(v.Stage[v.CurrentStage].SpecificImp * float32(v.G)) / 100
}


func (v *VEHICLE) setMdot(stageNumber int) {
	v.Stage[stageNumber].M_dot = v.Stage[stageNumber].Thrust/(v.Stage[v.CurrentStage].SpecificImp * float32(v.G))
}

// getGravAcceleration --------------------------------------------------------
// calculate acceleration of gravity at current altitude: 
//		g = GRAVITY_ACC * (r0/(r0+h))**2
// ----------------------------------------------------------------------------  
func (v *VEHICLE) setGravityAcceleration() {
	v.G = GRAVITY_ACC * math.Pow(EARTHRADIUS / (EARTHRADIUS + float64((*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude)), 2)
}

// setPolarDot ----------------------------------------------------------------
// calculate the rate of change between original referencial vertical and
// local position's vertical	
// 		polarDot = v * sin(gamma)/(r0+h) (gamma = angle between local vertical
//										  and rocket direction)
// ----------------------------------------------------------------------------  

// setGammaDot ----------------------------------------------------------------
// calculate the rate of change between vertical of local referencial and
// rocket direction	
// 		gammaDot = G * sin(gamma)/v - polarDot = G * sin(gamma)/v - v * sin(gamma)/(r0+h)
//
// 		v.d(gamma)/dt = - (g - v**2/(r0+h)).cos(gamma) => gammaDot = ((v/(r0+h)) - g/v) * cos(gamma)
// ----------------------------------------------------------------------------  

func (v *VEHICLE) setGammaDot() {
//	v.Velocity := float64((*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity)
//	h := float64((*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude)
	v.gamma_dot = ((math.Pow(v.Velocity, 2)/(EARTHRADIUS + v.Altitude)) - v.G) * math.Cos(v.Gamma)/v.Velocity
}

const GAMMA0 = 89
func (v *VEHICLE) updateCurveAngle(maxtime float32) {
	//tan φ = (1-t/T) * tan θ0
	v.Gamma = rad2deg(math.Atan((1 - float64(v.Clock/maxtime)) * math.Tan(deg2rad(GAMMA0))))
}

func (v *VEHICLE) setISAparams() {
	v.Altitude = float64((*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude)
	v.Velocity = float64((*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity)
	v.G = GRAVITY_ACC / math.Pow(1 + (v.Altitude/EARTHRADIUS), 2)
	v.Rho = RHO * math.Exp(-v.Altitude/H0)
}
func (v *VEHICLE) setMass(deltaT float32) {
	v.Stage[v.CurrentStage].PropellantMass = v.Stage[v.CurrentStage].PropellantMass - (v.getMdot(100) * deltaT)
}

func (v *VEHICLE) updateGravityTurn() {
	v.setISAparams()
	v.setDrag()

	// update propellant depletion
	//v.Stage[v.CurrentStage].PropellantMass = v.Stage[v.CurrentStage].PropellantMass - (v.getMdot() * DELTAt))
	v.setMass(DELTAt)
	m := float64(v.Stage[v.CurrentStage].PropellantMass + v.Stage[v.CurrentStage].DryMass) 
	
	// apply: dv/dt = (T - D)/m - g.sin(gamma)
	vdot := ((float64(v.Stage[v.CurrentStage].Thrust) - v.Drag)/ m) - v.G * math.Sin(deg2rad(v.Gamma))
	if vdot < 0 {
		vdot = 0
	}
	// apply: dx/dt = v.cos(gamma).r0/(r0+h)
 	xdot := (v.Velocity + vdot) * math.Cos(deg2rad(v.Gamma)) * (EARTHRADIUS/(EARTHRADIUS+v.Altitude)) * DELTAt
	// apply: dh/dt = v.sin(gamma)	
	hdot := (v.Velocity + vdot) * math.Sin(deg2rad(v.Gamma)) * DELTAt

	// update turn angle
	/*if v.Velocity > 0 {
		v.setGammaDot()
		v.Gamma = v.Gamma - (v.gammaDot * DELTAt)
	}*/
	v.updateCurveAngle(60)

	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude = float32(v.Altitude) + float32(hdot)
	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range = (*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range + float32(xdot)
	(*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity = float32(v.Velocity) + float32(vdot)
}

// setDrag --------------------------------------------------------------------
// calculate the drag force
// 		drag = 1/2 * rho * A * Cd * v.v  where:
//				rho = air density at a given altitude
//				A   = frontal surface exposed to drag
//				Cd	= drag coefficient, only dependent of shape (Cd = 0.42)
//				v	= velocity
// ----------------------------------------------------------------------------  
func (v *VEHICLE) setDrag() {
	v.Drag = 0.5 * v.Rho * v.FrontalArea * CD * math.Pow(v.Velocity, 2)
//	v.Drag = 0.5 * v.Rho * v.FrontalArea * CD * math.Pow(float64((*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity),2)
//	fmt.Println("--",rho,"--",v.FrontalArea,"--",CD,"--",(*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity,"--",v.Drag)
}

// steps to calculate parameters value after t+dt:
// 1) increment time by deltaT
// 2) calculate how much mass is being shot out during deltaT
// 3) calculate velocity
// 4) calculate altitude
// 4) update ISA parameters
const DELTAt = 1 //0.1   // in seconds

func (v *VEHICLE) setVelocity() {

	fmt.Println(v.Clock, "-->",(v.Velocity)*3.6, "k/h -- ", "Alt:",v.Altitude/1000,	"km, Downrange:", (*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range/1000,"km, Gamma=",v.Gamma)
	if v.Clock < 2.0 {
		v.setMass(DELTAt)
//		fmt.Println(v.Clock, "-->",0, "k/h -- ", "Alt:",v.Altitude,	"Downrange:", v.Range)
		v.Clock = v.Clock + DELTAt
		return
	}
	v.updateGravityTurn()
	/*
	v.setDrag()
	v.Stage[v.CurrentStage].PropellantMass = v.Stage[v.CurrentStage].PropellantMass - (v.getMdot() * DELTAt)
	m := float64(v.Stage[v.CurrentStage].PropellantMass + v.Stage[v.CurrentStage].DryMass) 
	V := float64((*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity)
	
	// apply: dv/dt = (T - D)/m - g.sin(gamma)
	vdot := ((float64(v.Stage[v.CurrentStage].Thrust) - v.Drag)/ m) - v.G * math.Sin(deg2rad(v.Gamma))
	if vdot < 0 {
		vdot = 0
	}
	// apply: dx/dt = v.cos(gamma).r0/(r0+h)
// 	xdot := (V+vdot) * math.Sin(v.Gamma) * DELTAt
 	xdot := (V+vdot) * math.Cos(deg2rad(v.Gamma)) * (EARTHRADIUS/(EARTHRADIUS+float64((*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude)))
	// apply: dh/dt = v.sin(gamma)	
//	hdot := (V+vdot) * math.Cos(v.Gamma) * DELTAt
	hdot := (V+vdot) * math.Sin(deg2rad(v.Gamma))

	(*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])).Velocity = float32(V + vdot)

	//v.setPolarDot()
	v.setGammaDot()
	//v.PolarAngle = v.PolarAngle + v.polarDot
	v.Gamma = v.Gamma + v.gammaDot
	//fmt.Println("gamma=",v.Gamma)

	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude = 
	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Altitude + float32(hdot)

	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range = 
	(*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])).Range + float32(xdot)
	*/
//	fmt.Println(v.Drag)
//	v.Stage[v.CurrentStage].PropellantMass = v.Stage[v.CurrentStage].PropellantMass - (v.getMdot() * DELTAt)
	v.Clock = v.Clock + DELTAt
	if v.Gamma == 90 && v.Clock > 5.0  {
		v.Gamma = v.Gamma - 5
		fmt.Println("NEW GAMMA:",v.Gamma)
	}
}

func deg2rad(degre float64) float64 {
	return math.Pi * degre / 180
}

func rad2deg(rad float64) float64 {
	return rad * (180/math.Pi)
}


const DOWNRANGE_PITCH = 130	// in meters	

var deg = math.Pi / 180
var rad = 180/ math.Pi

//var totalBurnTime = 0
const TOTAL_SEGMENTS = 20
const TARGET_ORBIT = 350000
func liftoff (v *VEHICLE) {
	v.G = GRAVITY_ACC
	v.Gamma = 89.85 * deg
	for k := range v.Stage {
		v.setMdot(k)
		v.Stage[k].BurnTime = v.Stage[k].PropellantMass / v.Stage[k].M_dot
		v.TotalBurnTime = v.TotalBurnTime + float64(v.Stage[k].BurnTime)
	}
	v.TargetAltitude = TARGET_ORBIT // in m
	v.OrbitalVelocity = math.Sqrt(UNIVERSALG * EARTHMASS / (v.TargetAltitude + EARTHRADIUS))
	
	v.TicksPerSegment = v.TotalBurnTime / TOTAL_SEGMENTS // each tick range
	v.EarlyTiltAngle = (20 / v.TicksPerSegment) * deg
	v.LateTiltAngle = (80 / ((TOTAL_SEGMENTS-1) * v.TicksPerSegment)) * deg
	fmt.Println("tickRange:", v.TicksPerSegment, "sec, earlyTA:", v.EarlyTiltAngle*rad, ", lateTA:", v.LateTiltAngle*rad, ", Vorbital:", v.OrbitalVelocity,"m/s, orbite altitude:", v.TargetAltitude/1000,"km")
}

// rocket burn duration divided into 20 segments
// segment 1 produces tilt of 20 degres
// segments 2-20 produces an 70 degre tilt
// 3 phases:
// 1 - rocket vertical (gamma = 90) until Hv is reached
// 2 - rocket tilts between 0 -> 20 degres clockwork
// 3 - rocket tilts between 21 -> 90
// in phase 3 the rate to apply can be:
// 		- constant with each segment contributing to a ratio of (80*19/20) degres (but doesn't allow for a target altitude)
//		- adapt to how much is left to reach a target altitude using this algo:
//			for each iteration:
//				calculate delta = TargetAltitude - currentAltitude
//				calculate velocity v
//				calculate acceleration
//				calculate velocity normal component : vn = v.sin(gamma)
//				calculate velocity tangent comp: vt = v.cos(gamma)
//				calculate time left: T = delta/vn
//				calculate dgamma = gamma/T

func(v *VEHICLE) setRates(mecoStatus bool) bool {
	const THROTTLE_VALUE = 0.8
	var m float32
	var meco = false
	m_dot := v.Stage[v.CurrentStage].M_dot * THROTTLE_VALUE

	if !mecoStatus && (v.Clock - v.ClockAtMeco) < v.Stage[v.CurrentStage].BurnTime {
		m = v.Stage[v.CurrentStage].DryMass + v.Stage[v.CurrentStage].PropellantMass - m_dot * (v.Clock - v.ClockAtMeco)
	} else {
		m = v.Stage[v.CurrentStage].DryMass + v.Stage[v.CurrentStage].PropellantMass - m_dot * v.Stage[v.CurrentStage].BurnTime
		v.Stage[v.CurrentStage].Thrust = 0
		v.ClockAtMeco = v.Clock
		meco = true
	}
	v.G = GRAVITY_ACC / math.Pow(1 + (v.Altitude/EARTHRADIUS), 2)
	v.Rho = RHO * math.Exp(-v.Altitude/H0)
	v.Drag = 0.5 * v.Rho * v.FrontalArea * CD * math.Pow(v.Velocity, 2)
	if v.Altitude <= DOWNRANGE_PITCH {
		v.VerticalTicks++
		v.gamma_dot = 0
		v.v_dot = ((float64(v.Stage[v.CurrentStage].Thrust) - v.Drag)/ float64(m)) - v.G
		v.x_dot = 0
		v.h_dot = v.Velocity
		v.vG_dot = - v.G
	} else {
		v.v_dot = ((float64(v.Stage[v.CurrentStage].Thrust) - v.Drag)/ float64(m)) - v.G * math.Sin(v.Gamma)
		v.x_dot = (EARTHRADIUS/(EARTHRADIUS+v.Altitude)) * v.Velocity * math.Cos(v.Gamma)
		v.vG_dot = -v.G * math.Sin(v.Gamma)
		v.h_dot = v.Velocity * math.Sin(v.Gamma)
		v.updateGammaDot()
	}
	v.vD_dot = - v.Drag / float64(m)
	return meco
}

func (v *VEHICLE) updateRocketSensors() {
	v.Velocity = v.Velocity + v.v_dot
	v.Altitude = v.Altitude + v.h_dot
	v.Range = v.Range + v.x_dot
	//v.Gamma = v.Gamma + v.gamma_dot
	v.updateGamma()
	v.vG = v.vG + v.vG_dot
	v.vD = v.vD + v.vD_dot
}

func (v * VEHICLE) updateGammaDot() {
//	v.gamma_dot = -(1/v.Velocity) * (v.G - (math.Pow(v.Velocity, 2)/(EARTHRADIUS + v.Altitude))) * math.Cos(v.Gamma)
	if v.Clock - float32(v.VerticalTicks) < float32(v.TicksPerSegment) {
		v.gamma_dot = v.EarlyTiltAngle
	} else {
		v.gamma_dot = v.LateTiltAngle
	}
}
func (v * VEHICLE) updateGamma() {
	//tan φ = (1-t/T) * tan θ0
	//v.Gamma = math.Atan((1 - float64(v.Clock)/v.TotalBurnTime) * math.Tan(GAMMA0*deg))
	v.Gamma = v.Gamma - v.gamma_dot
}