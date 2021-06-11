package main

import (
	"math"

	"github.com/crmathieu/daq/packages/data"
)

// there is only one angle to consider: gamma. This is the angle between the velocity vector and
// the local horizon. Initially, gamma = PI/2 and progressively should converge to 0

// to calculate the various variable values each time increment we will:
// 1) calculate the amount of gas ejected, and the new mass of the vehicle
// 2) calculate the thrust T
// 3) calculate the drag D
// 4) calculate the velocity V
// 5) calculate the downrange X
// 6) calculate the altitude H
// 7) calculate the variation of flightpath gamma

// the idea to reach a desired orbit is to try different initial values for the pitch
// kick-off time and the burnout time until we find the values that converge to that orbit.

// We use a referential with its origin at the center of the earth. Initially THE COORDINATES
// are x = 0, y = Re

func (r *VEHICLE) liftOff() bool {

	r.Stages[BOOSTER].vAx = vE // Absolute velocity in x-direction = velocity of earth at surface
	r.Stages[BOOSTER].vAy = 0  //r.Stages[BOOSTER].vAy + (r.Stages[BOOSTER].ay * r.Stages[BOOSTER].dt/2)

	r.Stages[BOOSTER].DTF = r.Stages[BOOSTER].cy //math.Sqrt(r.Stages[BOOSTER].cx * r.Stages[BOOSTER].cx + r.Stages[BOOSTER].cy * r.Stages[BOOSTER].cy)
	r.Stages[BOOSTER].Acc = 0                    //math.Sqrt(r.Stages[BOOSTER].ax * r.Stages[BOOSTER].ax + r.Stages[BOOSTER].ay * r.Stages[BOOSTER].ay)

	r.EventsMap = r.EventsMap | data.E_LIFTOFF
	r.LastEvent = data.E_LIFTOFF

	return true
}

//inline void liftOff()
func (r *VEHICLE) liftOff2() bool {
	Engine := EnginesMap[r.Stages[BOOSTER].EngineID]

	r.Stages[BOOSTER].ForceY = (float64(r.Stages[BOOSTER].RunningEngines) * Engine.Th_sl) - r.Stages[BOOSTER].Mass*g0
	r.Stages[BOOSTER].ay = r.Stages[BOOSTER].ForceY / r.Stages[BOOSTER].Mass

	r.Stages[BOOSTER].vAx = vE // Absolute velocity in x-direction = velocity of earth at surface
	r.Stages[BOOSTER].vAy = r.Stages[BOOSTER].vAy + (r.Stages[BOOSTER].ay * r.Stages[BOOSTER].dt / 2)

	r.Stages[BOOSTER].DTF = math.Sqrt(r.Stages[BOOSTER].cx*r.Stages[BOOSTER].cx + r.Stages[BOOSTER].cy*r.Stages[BOOSTER].cy)
	r.Stages[BOOSTER].Acc = math.Sqrt(r.Stages[BOOSTER].ax*r.Stages[BOOSTER].ax + r.Stages[BOOSTER].ay*r.Stages[BOOSTER].ay)

	r.EventsMap = r.EventsMap | data.E_LIFTOFF
	r.LastEvent = data.E_LIFTOFF

	return true
}

/*
	Leapfrog integrator for moving rocket
*/

// The launch azimuth is the angle between north direction and the projection of the
// initial orbital plane onto the launch location. ... The orbital inclination is the
// angle between the orbital plane and the celestial body's reference plane. If the
// body spins then this is usually the equatorial plane
func (r *VEHICLE) timeStep(i int32) { // i = stage

	dgamma := float64(0.0)
	if r.SysGuidance._MEI1 {
		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * EnginesMap[r.Stages[i].EngineID].Flow_rate * r.Stages[i].dt
		r.Stages[i].Mf = r.Stages[i].Mf - dm
		r.Stages[i].Mass = r.Stages[i].Mass - dm
	}

	if r.SysGuidance._release {
		aerodynPressure = 0.5 * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative * 1e-3                       // Aerodynamic stress
		drag = (0.5) * r.Stages[i].Cd * r.Stages[i].CSArea * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative // Drag
		r.Stages[i].Thrust = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * r.GetThrust(r.Stages[i].DTF, i)        // Thrust

		// calculate force and velocity vectors norm
		gh := g(r.Stages[i].DTF)
		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*gh*math.Sin(r.Stages[i].gamma)

		// update downrange first dx/dt = (Re/(Re+h))*v*cos(gamma)
		// since we are using altitude in the range calculation formula, calculate range first so that both
		// altitude and range are calculated simultaneously
		// first update range as a fonction of old relative speed, old altitude and "dt" time increment
		r.Stages[i].drange = (Re/(Re+r.Stages[i].altitude))*r.Stages[i].RVel*math.Cos(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].drange

		// second update altitude: dh/dt = v.sin(gamma) as a function of old relative speed and time increment "dt"
		r.Stages[i].altitude = r.Stages[i].RVel*math.Sin(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].altitude

		// calculate dgamma/dt = -(g - v*v/(Re+h)) * cos(gamma) * 1/v
		//dgamma = -(gh - (r.Stages[i].Vel * r.Stages[i].Vel)/(Re + r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) * (1/r.Stages[i].Vel)
		// dgamma = gamma /()
		//		dgamma = (profile.OrbitPerigee - r.Stages[i].altitude)*1e-3 / (profile.BurnoutTime - r.Stages[i].Clock)

		//r.TotalTimeIncrement = qTotalTimeIncrement + dt

		dgamma = r.gravTurnTangentSteering(i)
		//if dgamma != 0 {
		//	println(dgamma)
		//}

		//		dgamma = r.gravTurn(i)

		// calculate cartesian coordinates
		r.Stages[i].alpha = r.Stages[i].drange / Re // polar angle (in rd) based on downrange value
		r.Stages[i].cx = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].alpha)
		r.Stages[i].cy = (Re + r.Stages[i].altitude) * math.Cos(r.Stages[i].alpha)

		if false {
			// x-direction
			//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
			//		r.Stages[i].ForceX = r.Stages[i].Force * math.Cos(r.Stages[i].gamma)
			//		r.Stages[i].cx = r.Stages[i].cx + r.Stages[i].vAx * r.Stages[i].dt
			r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass

			r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
			r.Stages[i].vAx = r.Stages[i].vRx + vE //* math.Sin(r.Stages[i].beta)

			// y-direction
			//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
			//		r.Stages[i].cy = r.Stages[i].cy + r.Stages[i].vAy * r.Stages[i].dt
			r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

			r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
			r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

			r.Stages[i].DTF = math.Sqrt(r.Stages[i].cx*r.Stages[i].cx + r.Stages[i].cy*r.Stages[i].cy)
			r.Stages[i].VAbsolute = math.Sqrt(r.Stages[i].vAx*r.Stages[i].vAx + r.Stages[i].vAy*r.Stages[i].vAy)
			r.Stages[i].VRelative = math.Sqrt(r.Stages[i].vRx*r.Stages[i].vRx + r.Stages[i].vRy*r.Stages[i].vRy)
			r.Stages[i].Acc = r.Stages[i].Force / r.Stages[i].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)

		} else {
			// NEW NEW
			r.Stages[i].VAbsolute = r.Stages[i].AVel
			r.Stages[i].VRelative = r.Stages[i].RVel
			r.Stages[i].DTF = Re + r.Stages[i].altitude
			r.Stages[i].Acc = r.Stages[i].Force / r.Stages[i].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)
		}

		// update velocities
		r.Stages[i].RVel = r.Stages[i].RVel + (r.Stages[i].Force/r.Stages[i].Mass)*r.Stages[i].dt
		r.Stages[i].AVel = r.Stages[i].RVel + vE

	}
	/*
		if r.SysGuidance._BBURN  || r.SysGuidance._LBURN {
			r.flip(0)
		}
		if r.SysGuidance._LBURN && mod(r.Stages[i].Clock, 5) < r.Stages[i].dt {
			r.update_landing_throttle()
		}
	*/

	/*
		r.Stages[i].cx = r.Stages[i].cx + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].cy = r.Stages[i].cy + r.Stages[i].vAy * r.Stages[i].dt
	*/

	if r.Stages[i].gamma >= math.Abs(dgamma) {
		r.Stages[i].gamma = r.Stages[i].gamma + dgamma //- math.Abs(dgamma)
	} else {
		r.Stages[i].gamma = 0
	}

	//	fmt.Println("gamma = ", rad2deg(r.Stages[i].gamma))
}

func (r *VEHICLE) timeStepSAVE(i int32) { // i = stage

	dgamma := float64(0.0)
	if r.SysGuidance._MEI1 {
		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * EnginesMap[r.Stages[i].EngineID].Flow_rate * r.Stages[i].dt
		r.Stages[i].Mf = r.Stages[i].Mf - dm
		r.Stages[i].Mass = r.Stages[i].Mass - dm
	}

	if r.SysGuidance._release {
		aerodynPressure = 0.5 * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative * 1e-3                       // Aerodynamic stress
		drag = (0.5) * r.Stages[i].Cd * r.Stages[i].CSArea * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative // Drag
		r.Stages[i].Thrust = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * r.GetThrust(r.Stages[i].DTF, i)        // Thrust

		// calculate force and velocity vectors norm
		gh := g(r.Stages[i].DTF)
		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*gh*math.Sin(r.Stages[i].gamma)
		r.Stages[i].RVel = r.Stages[i].RVel + (r.Stages[i].Force/r.Stages[i].Mass)*r.Stages[i].dt
		r.Stages[i].AVel = r.Stages[i].RVel + vE

		// calculate altitude: dh/dt = v.sin(gamma)
		r.Stages[i].altitude = r.Stages[i].RVel*math.Sin(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].altitude

		// calculate downrange dx/dt = (Re/(Re+h))*v*cos(gamma)
		r.Stages[i].drange = (Re/(Re+r.Stages[i].altitude))*r.Stages[i].RVel*math.Cos(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].drange

		// calculate dgamma/dt = -(g - v*v/(Re+h)) * cos(gamma) * 1/v
		//dgamma = -(gh - (r.Stages[i].Vel * r.Stages[i].Vel)/(Re + r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) * (1/r.Stages[i].Vel)
		// dgamma = gamma /()
		//		dgamma = (profile.OrbitPerigee - r.Stages[i].altitude)*1e-3 / (profile.BurnoutTime - r.Stages[i].Clock)

		//dgamma = r.gravTurnClassic(i)

		//dgamma = r.gravTurn(i)

		// calculate cartesian coordinates
		r.Stages[i].alpha = r.Stages[i].drange / Re // polar angle (in rd) based on downrange value
		r.Stages[i].cx = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].alpha)
		r.Stages[i].cy = (Re + r.Stages[i].altitude) * math.Cos(r.Stages[i].alpha)

		// x-direction
		//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
		//		r.Stages[i].ForceX = r.Stages[i].Force * math.Cos(r.Stages[i].gamma)
		//		r.Stages[i].cx = r.Stages[i].cx + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass

		r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
		r.Stages[i].vAx = r.Stages[i].vRx + vE //* math.Sin(r.Stages[i].beta)

		// y-direction
		//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
		//		r.Stages[i].cy = r.Stages[i].cy + r.Stages[i].vAy * r.Stages[i].dt
		r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

		r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
		r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

		r.Stages[i].DTF = math.Sqrt(r.Stages[i].cx*r.Stages[i].cx + r.Stages[i].cy*r.Stages[i].cy)
		r.Stages[i].VAbsolute = math.Sqrt(r.Stages[i].vAx*r.Stages[i].vAx + r.Stages[i].vAy*r.Stages[i].vAy)
		r.Stages[i].VRelative = math.Sqrt(r.Stages[i].vRx*r.Stages[i].vRx + r.Stages[i].vRy*r.Stages[i].vRy)
		r.Stages[i].Acc = r.Stages[i].Force / r.Stages[i].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)
	}
	/*
		if r.SysGuidance._BBURN  || r.SysGuidance._LBURN {
			r.flip(0)
		}
		if r.SysGuidance._LBURN && mod(r.Stages[i].Clock, 5) < r.Stages[i].dt {
			r.update_landing_throttle()
		}
	*/

	/*
		r.Stages[i].cx = r.Stages[i].cx + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].cy = r.Stages[i].cy + r.Stages[i].vAy * r.Stages[i].dt
	*/
	if r.Stages[i].gamma >= math.Abs(dgamma) {
		r.Stages[i].gamma = r.Stages[i].gamma - math.Abs(dgamma)
	} else {
		r.Stages[i].gamma = 0
	}

	//	fmt.Println("gamma = ", rad2deg(r.Stages[i].gamma))
}

var nearOrbit = false
var gammabits float64

// called with every tick
func (r *VEHICLE) gravTurnClassic(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep && i == BOOSTER {
			//println("NEAH!")
			return 0
		}
		r.dtIncrements = r.dtIncrements + r.Stages[i].dt
		// after 1 sec using the same rate, we calculate a new rate/s
		if r.dtIncrements > 1.0 {
			println(r.dtIncrements)
			r.dtIncrements = 0
			gh := g(r.Stages[i].DTF)
			// calculate the rate dgamma/dt = -(g - v*v/(Re+h)) * cos(gamma) * 1/v
			r.dgammaPerSec = -(gh - (r.Stages[i].AVel*r.Stages[i].AVel)/(Re+r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) * (1 / r.Stages[i].AVel) //* r.Stages[i].dt
			r.dgamma = r.dgammaPerSec * r.Stages[i].dt
			return r.dgammaPerSec
			println("dgammaPerSec = ", rad2deg(r.dgammaPerSec), "deg/s - r.dgamma = ", r.dgamma)
		}
		return 0 //r.dgamma
	}
	return 0
}

var GAMMA0 = M_PI/2 - 0.05

func (r *VEHICLE) gravTurnTangentSteering(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER || !r.SysGuidance._SEI1 {
				return 0
			}
		}
		//r.dtIncrements = r.dtIncrements + r.Stages[i].dt
		// implements: tanθ(t)=tan(θ0 *(1 - altitude/orbitInsertion))
		// gamma := math.Atan(math.Tan(GAMMA0 * (1 - altitude/profile.OrbitInsertion)))
		if r.Stages[i].altitude <= profile.OrbitInsertion {
			gamma := math.Atan(math.Tan(GAMMA0 * (1 - (r.Stages[i].altitude)/profile.OrbitInsertion)))
			return gamma - r.Stages[i].gamma
		} else {
			r.Stages[i].gamma = 0
			return 0
		}
	}
	return 0
}
