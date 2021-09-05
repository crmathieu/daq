package main

import (
	"fmt"
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

	fmt.Printf("\n************\nLLIFT OFF @ ---> %g seconds\n", r.Stages[BOOSTER].Clock)
	fmt.Println("BOOSTER fuel amount ...... ", r.Stages[BOOSTER].Mf, "kg")

	r.Stages[BOOSTER].Thrust = float64(r.Stages[BOOSTER].RunningEngines) * r.Stages[BOOSTER].ThrottleRate * r.GetThrust(r.Stages[BOOSTER].DTF, BOOSTER) // Thrust
	fmt.Println("BOOSTER Thrust ........... ", r.Stages[BOOSTER].Thrust, "Nm")

	r.Stages[BOOSTER].vAx = vE // Absolute velocity in x-direction = velocity of earth at surface
	r.Stages[BOOSTER].vAy = 0  //r.Stages[BOOSTER].vAy + (r.Stages[BOOSTER].ay * r.Stages[BOOSTER].dt/2)

	r.Stages[BOOSTER].DTF = Re //r.Stages[BOOSTER].py //math.Sqrt(r.Stages[BOOSTER].px * r.Stages[BOOSTER].px + r.Stages[BOOSTER].py * r.Stages[BOOSTER].py)
	r.Stages[BOOSTER].Acc = 0  //math.Sqrt(r.Stages[BOOSTER].ax * r.Stages[BOOSTER].ax + r.Stages[BOOSTER].ay * r.Stages[BOOSTER].ay)

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

	r.Stages[BOOSTER].DTF = math.Sqrt(r.Stages[BOOSTER].px*r.Stages[BOOSTER].px + r.Stages[BOOSTER].py*r.Stages[BOOSTER].py)
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
var accu = 0
var zob = 0.0

func (r *VEHICLE) timeStepWEIRD(i int32) { // i = stage

	dgamma := float64(0.0)
	if r.SysGuidance._stagesep && i == 0 {
		fmt.Println("#############################################")
	}

	if r.SysGuidance._release {
		accu += 1
		//if accu > 20000 {
		//	os.Exit(0)
		//}
		//fmt.Println("DTF=", r.Stages[i].DTF, "for stage", i)
		density := rho(r.Stages[i].DTF - Re)
		//fmt.Println("density=", density)
		aerodynPressure = 0.5 * density * r.Stages[i].VRelative * r.Stages[i].VRelative * 1e-3 // Aerodynamic stress
		fmt.Println("Current Q=", aerodynPressure, " vs MaxQ=", r.MaxQ)
		if i == BOOSTER && aerodynPressure > r.MaxQ && !r.SysGuidance._MECO1 {
			fmt.Println("########################")
			r.MaxQ = aerodynPressure
			r.AltMaxQ = r.Stages[i].altitude
			zob = r.Stages[i].VRelative
		}
		drag = (0.5) * r.Stages[i].Cd * r.Stages[i].CSArea * density * r.Stages[i].VRelative * r.Stages[i].VRelative          // Drag
		r.Stages[i].Thrust = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * r.GetThrust(r.Stages[i].DTF, i) // Thrust
		// calculate force and velocity vectors norm
		gh := g(r.Stages[i].DTF)
		//if r.Stages[i].vAy < 0 {
		//	drag *= -1
		//}
		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*gh*math.Sin(r.Stages[i].gamma)
		//fmt.Println("Thrust", r.Stages[i].Thrust, "drag=", drag, "FORCE=", r.Stages[i].Force)

		//if r.Stages[i].Force < 0 {
		//	fmt.Println(r.Stages[i].Force)
		//}
		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].AVel*math.Cos(r.Stages[i].gamma), 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)

		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].vAx, 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)
		r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) * math.Sin(deg2rad(profile.LaunchAzimuth)) / r.Stages[i].Mass
		r.Stages[i].az = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) * math.Cos(deg2rad(profile.LaunchAzimuth)) / r.Stages[i].Mass
		r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass
		//fmt.Println("AX=", r.Stages[i].ax, ", AY=", r.Stages[i].ay, ", AZ=", r.Stages[i].az)

		// update downrange first dx/dt = (Re/(Re+h))*v*cos(gamma)
		// since we are using altitude in the range calculation formula, calculate range first so that both
		// altitude and range are calculated simultaneously
		// first update range as a fonction of old relative speed, old altitude and "dt" time increment

		//r.Stages[i].drange = (Re/(Re+r.Stages[i].altitude))*r.Stages[i].RVel*math.Cos(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].drange

		// second update altitude: dh/dt = v.sin(gamma) as a function of old relative speed and time increment "dt"
		//r.Stages[i].altitude = r.Stages[i].RVel*math.Sin(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].altitude

		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vx*r.Stages[i].dt
		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vy*r.Stages[i].dt
		r.Stages[i].pz = r.Stages[i].pz + r.Stages[i].vz*r.Stages[i].dt
		//fmt.Println("CX=", r.Stages[i].px, ", CY=", r.Stages[i].py, ", CZ=", r.Stages[i].pz)
		// calculate range from X value and azimuth angle

		r.Stages[i].drange = r.Stages[i].px / math.Sin(profile.LaunchAzimuth)
		r.Stages[i].altitude = r.Stages[i].py
		//fmt.Println("DT", r.Stages[i].dt, "CX=", r.Stages[i].px, ", CY=", r.Stages[i].py, ", CZ=", r.Stages[i].pz, "VX=", r.Stages[i].vx, ", VY=", r.Stages[i].vy, ", VZ=", r.Stages[i].vz)

		r.Stages[i].vx = r.Stages[i].vx + r.Stages[i].ax*r.Stages[i].dt
		r.Stages[i].vy = r.Stages[i].vy + r.Stages[i].ay*r.Stages[i].dt
		r.Stages[i].vz = r.Stages[i].vz + r.Stages[i].az*r.Stages[i].dt
		//fmt.Println("VX=", r.Stages[i].vx, ", VY=", r.Stages[i].vy, ", VZ=", r.Stages[i].vz)

		// calculate dgamma/dt = -(g - v*v/(Re+h)) * cos(gamma) * 1/v
		//dgamma = -(gh - (r.Stages[i].Vel * r.Stages[i].Vel)/(Re + r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) * (1/r.Stages[i].Vel)
		// dgamma = gamma /()
		//		dgamma = (profile.OrbitPerigee - r.Stages[i].altitude)*1e-3 / (profile.BurnoutTime - r.Stages[i].Clock)

		//r.TotalTimeIncrement = qTotalTimeIncrement + dt

		//		dgamma = r.gravTurnTangentSteering(i)

		dgamma = r.gravTurnMultiPhaseTangentSteering(i)
		fmt.Println(i, r.Stages[i].gamma)
		//dgamma = r.gravTurnClassic(i)

		//if dgamma != 0 {
		//	println(dgamma)
		//}

		//		dgamma = r.gravTurn(i)

		// calculate cartesian coordinates
		//r.Stages[i].beta = r.Stages[i].drange / Re // polar angle (in rd) based on downrange value
		//r.Stages[i].px = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].beta)
		//r.Stages[i].py = (Re + r.Stages[i].altitude) * math.Cos(r.Stages[i].beta)

		// update velocities
		r.Stages[i].RVel = r.Stages[i].RVel + (r.Stages[i].Force/r.Stages[i].Mass)*r.Stages[i].dt
		r.Stages[i].RVel = math.Sqrt(math.Pow(r.Stages[i].vx, 2) + math.Pow(r.Stages[i].vy, 2) + math.Pow(r.Stages[i].vz, 2))

		//fmt.Println("RVEL=", r.Stages[i].RVel, "From coordinate=", math.Sqrt(math.Pow(r.Stages[i].vx, 2)+math.Pow(r.Stages[i].vy, 2)+math.Pow(r.Stages[i].vz, 2)))
		r.Stages[i].AVel = r.Stages[i].RVel + vE

		if false {
			// x-direction
			//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
			//		r.Stages[i].ForceX = r.Stages[i].Force * math.Cos(r.Stages[i].gamma)
			//		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
			r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass

			r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
			r.Stages[i].vAx = r.Stages[i].vRx + vE //* math.Sin(r.Stages[i].beta)

			// y-direction
			//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
			//		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
			r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

			r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
			r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

			r.Stages[i].DTF = math.Sqrt(r.Stages[i].px*r.Stages[i].px + r.Stages[i].py*r.Stages[i].py)
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
		//r.Stages[i].RVel = r.Stages[i].RVel + (r.Stages[i].Force/r.Stages[i].Mass)*r.Stages[i].dt
		//fmt.Println("RVEL=", r.Stages[i].RVel, "From coordinate=", math.Sqrt(math.Pow(r.Stages[i].vx, 2)+math.Pow(r.Stages[i].vy, 2)+math.Pow(r.Stages[i].vz, 2)))
		//r.Stages[i].AVel = r.Stages[i].RVel + vE
		//		r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
		//		r.Stages[i].vAx = r.Stages[i].vRx + vE

		//		r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
		//		r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

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
		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
	*/

	if r.Stages[i].gamma >= math.Abs(dgamma) {
		r.Stages[i].gamma = r.Stages[i].gamma + dgamma //- math.Abs(dgamma)
	} else {
		r.Stages[i].gamma = 0
	}
	//fmt.Println("GAMMA=", r.Stages[i].gamma)
	if r.SysGuidance._MEI1 {
		if r.Stages[i].Mf < 0.05 {
			//fmt.Printf("STAGE %v ", r.Stages[i]) //, "is EMPTY!!!!")
			r.Stages[i].RunningEngines = 0
			fmt.Printf("\n************\nPrematured MECO @ ---> %g seconds\n", r.Stages[i].Clock) //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Remaining fuel ...... ", r.Stages[i].Mf, "kg")                          //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Velocity ............ ", r.Stages[i].RVel*3.6, "km/h")                  //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			r.SysGuidance._MECO1 = r.MSECO(BOOSTER, data.E_MECO_1)
			r.SysGuidance._MEI1 = false
			return
		}
		//fmt.Println(int(r.Stages[i].Clock), "STAGE", i, "fuel left:", r.Stages[i].Mf)
		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * EnginesMap[r.Stages[i].EngineID].Flow_rate * r.Stages[i].dt
		r.Stages[i].Mf = r.Stages[i].Mf - dm
		r.Stages[i].Mass = r.Stages[i].Mass - dm
		//if r.SysGuidance._release {
		//	fmt.Println("dm=", dm, "fuel=", r.Stages[i].Mf, "Total-mass=", r.Stages[i].Mass)
		//}
	}

	//	fmt.Println("gamma = ", rad2deg(r.Stages[i].gamma))
}

//var angle = 0.0
//var ztime = 0.0

/* MAXQ parameters
{
          "time": 74.533,
          "velocity": 446.391,
          "altitude": 12.677,
          "velocity_y": 363.125,
          "velocity_x": 259.625,
          "acceleration": 19.022,
          "downrange_distance": 4.419,
          "angle": 54.426,
          "q": 29106.796978492315
        },
*/
type MAXQinfo struct {
	MaxQ, Alt, Velocity, Range, Angle, Time, RhoMQ float64
}

var mQ = MAXQinfo{}
var mecoAltitude float64

func (r *VEHICLE) OpenLoop() {

	// openLoop takes velocity and altitude and estimates how much acceleration
	// the vehicle needs to meet its orbit requirements (Vor, Aor). The function
	// interpolates the trajectory as a linear function
	//r.Stages[STAGE2].ThrottleRate = 1 - ((r.Stages[STAGE2].altitude - mecoAltitude) / profile.OrbitInsertion)
	//r.Stages[STAGE2].ThrottleRate = math.Exp(-math.Pow(r.Stages[STAGE2].altitude/(profile.OrbitInsertion), 2))

	// reverse sigmoid s = np.exp(var*x) / (1 + np.exp(var*x))
	//diff := -profile.OrbitInsertion + r.Stages[STAGE2].altitude
	//r.Stages[STAGE2].ThrottleRate = math.Exp(diff) / (1 + math.Exp(diff))

	// calculate deltaV to reach orbit:
	//deltaAltitude := Profile.OrbitInsertion - r.Stages[i].altitude
	deltaV := TargetOrbitalVelocity - r.Stages[STAGE2].AVel
	if deltaV < TargetOrbitalVelocity*0.01 { //*0.0045 { //220 {
		//fmt.Println(deltaV)
		//r.Stages[STAGE2].ThrottleRate = (profile.OrbitInsertion - r.Stages[STAGE2].altitude) / profile.OrbitInsertion

		r.Stages[STAGE2].ThrottleRate = deltaV / TargetOrbitalVelocity

		//r.Stages[STAGE2].ThrottleRate = 1 - (math.Log((r.Stages[STAGE2].altitude * r.Stages[STAGE2].AVel) / (profile.OrbitInsertion * TargetOrbitalVelocity)))
		//r.Stages[STAGE2].ThrottleRate = 1 - ((r.Stages[STAGE2].altitude - mecoAltitude) / profile.OrbitInsertion)

	}
	/*
		vy/vx -> 0
		H-h/H -> 0
	*/
}

func (r *VEHICLE) timeStep(i int32) { // i = stage

	dgamma := float64(0.0)
	if r.SysGuidance._stagesep && i == 0 {
		fmt.Println("#############################################")
	}

	if r.SysGuidance._release {
		aerodynPressure = 0.5 * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative // * 1e-3 // Aerodynamic stress

		//fmt.Println("Current Q=", aerodynPressure, " vs MaxQ=", r.MaxQ)
		if i == BOOSTER && !r.SysGuidance._MECO1 && aerodynPressure > mQ.MaxQ {
			//fmt.Println("########################")
			mQ.MaxQ = aerodynPressure
			mQ.Alt = r.Stages[i].altitude
			mQ.Velocity = r.Stages[i].VRelative
			mQ.Range = r.Stages[i].drange
			mQ.Angle = rad2deg(r.Stages[i].gamma)
			mQ.Time = r.Stages[i].Clock
			mQ.RhoMQ = rho(r.Stages[i].DTF - Re)
		}

		/*		if aerodynPressure > r.MaxQ {
				if i == STAGE2 {
					panic("something is wrong")
				}
				r.MaxQ = aerodynPressure
				r.AltMaxQ = r.Stages[BOOSTER].altitude
			}*/
		drag = (0.5) * r.Stages[i].Cd * r.Stages[i].CSArea * rho(r.Stages[i].DTF-Re) * r.Stages[i].VRelative * r.Stages[i].VRelative // Drag
		r.Stages[i].Thrust = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * r.GetThrust(r.Stages[i].DTF, i)        // Thrust
		// calculate force and velocity vectors norm
		gh := g(r.Stages[i].DTF)
		if r.Stages[i].vAy < 0 {
			drag *= -1
		}
		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*gh*math.Sin(r.Stages[i].gamma)
		//if r.Stages[i].Force < 0 {
		//	fmt.Println(r.Stages[i].Force)
		//}
		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].AVel*math.Cos(r.Stages[i].gamma), 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)

		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].vAx, 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)
		///////////////////////

		r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) * math.Sin(deg2rad(profile.LaunchAzimuth)) / r.Stages[i].Mass
		r.Stages[i].az = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) * math.Cos(deg2rad(profile.LaunchAzimuth)) / r.Stages[i].Mass
		r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

		///////////////////////
		// THE FOLLLOWING WORKS AND NEEDS TO BE UNCOMMENTED
		//r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass
		//r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass
		///////////////
		// update downrange first dx/dt = (Re/(Re+h))*v*cos(gamma)
		// since we are using altitude in the range calculation formula, calculate range first so that both
		// altitude and range are calculated simultaneously
		// first update range as a fonction of old relative speed, old altitude and "dt" time increment
		r.Stages[i].drange = (Re/(Re+r.Stages[i].altitude))*r.Stages[i].RVel*math.Cos(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].drange

		// second update altitude: dh/dt = v.sin(gamma) as a function of old relative speed and time increment "dt"
		r.Stages[i].altitude = r.Stages[i].RVel*math.Sin(r.Stages[i].gamma)*r.Stages[i].dt + r.Stages[i].altitude

		dgamma = r.TwoPhaseTangentSteering(i)

		// calculate cartesian coordinates
		r.Stages[i].beta = r.Stages[i].drange / Re // polar angle (in rd) based on downrange value
		// when beta greater than 2xPI, we have orbit
		if r.Stages[i].beta >= 2*math.Pi {
			orbit = true
		}
		r.Stages[i].px = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].beta) * math.Sin(deg2rad(profile.LaunchAzimuth))
		r.Stages[i].py = (Re + r.Stages[i].altitude) * math.Cos(r.Stages[i].beta)
		r.Stages[i].pz = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].beta) * math.Cos(deg2rad(profile.LaunchAzimuth))

		// NEW NEW
		r.Stages[i].VAbsolute = r.Stages[i].AVel
		r.Stages[i].VRelative = r.Stages[i].RVel
		r.Stages[i].DTF = Re + r.Stages[i].altitude
		r.Stages[i].Acc = r.Stages[i].Force / r.Stages[i].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)

		// update velocities
		r.Stages[i].RVel = r.Stages[i].RVel + (r.Stages[i].Force/r.Stages[i].Mass)*r.Stages[i].dt
		r.Stages[i].AVel = r.Stages[i].RVel + vE
	}

	if r.Stages[i].gamma >= math.Abs(dgamma) {
		r.Stages[i].gamma = r.Stages[i].gamma + dgamma //- math.Abs(dgamma)
	} else {
		r.Stages[i].gamma = 0
	}

	if r.SysGuidance._MEI1 {
		if r.Stages[i].Mf < 0.05 {
			//fmt.Printf("STAGE %v ", r.Stages[i]) //, "is EMPTY!!!!")
			r.Stages[i].RunningEngines = 0
			fmt.Printf("\n************\nPrematured MECO @ ---> %g seconds\n", r.Stages[i].Clock) //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Remaining fuel ...... ", r.Stages[i].Mf, "kg")                          //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Velocity ............ ", r.Stages[i].RVel*3.6, "km/h")                  //(-3*M_PI/2+v.Stages[STAGE2].beta-v.Stages[STAGE2].beta)*180/M_PI)
			r.SysGuidance._MECO1 = r.MSECO(BOOSTER, data.E_MECO_1)
			r.SysGuidance._MEI1 = false
			return
		}
		r.OpenLoop()
		//fmt.Println(int(r.Stages[i].Clock), "STAGE", i, "fuel left:", r.Stages[i].Mf)
		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * EnginesMap[r.Stages[i].EngineID].Flow_rate * r.Stages[i].dt
		r.Stages[i].Mf = r.Stages[i].Mf - dm
		r.Stages[i].Mass = r.Stages[i].Mass - dm
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

		dgamma = r.gravTurnMultiPhaseTangentSteering(i)

		//dgamma = r.gravTurn(i)

		// calculate cartesian coordinates
		r.Stages[i].beta = r.Stages[i].drange / Re // polar angle (in rd) based on downrange value
		r.Stages[i].px = (Re + r.Stages[i].altitude) * math.Sin(r.Stages[i].beta)
		r.Stages[i].py = (Re + r.Stages[i].altitude) * math.Cos(r.Stages[i].beta)

		// x-direction
		//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
		//		r.Stages[i].ForceX = r.Stages[i].Force * math.Cos(r.Stages[i].gamma)
		//		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass

		r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
		r.Stages[i].vAx = r.Stages[i].vRx + vE //* math.Sin(r.Stages[i].beta)

		// y-direction
		//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
		//		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
		r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

		r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
		r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

		r.Stages[i].DTF = math.Sqrt(r.Stages[i].px*r.Stages[i].px + r.Stages[i].py*r.Stages[i].py)
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
		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
	*/
	if r.Stages[i].gamma >= math.Abs(dgamma) {
		r.Stages[i].gamma = r.Stages[i].gamma + dgamma //- math.Abs(dgamma)
	} else {
		r.Stages[i].gamma = 0
	}

	//	fmt.Println("gamma = ", rad2deg(r.Stages[i].gamma))
}

var nearOrbit = false
var gammabits float64

// gravTurnClassic ------------------------------------------------------------
// steering program using differential equation to calculate instant rate if
// change of flightpath angle
// ----------------------------------------------------------------------------
func (r *VEHICLE) gravTurnClassic(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER || !r.SysGuidance._SEI1 {
				return 0
			}
		}
		gh := g(r.Stages[i].DTF)
		instantDeviation := -(1 / r.Stages[i].RVel) * (gh - (math.Pow(r.Stages[i].RVel, 2))/(Re+r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) //* r.Stages[i].dt
		return instantDeviation * r.Stages[i].dt
	}
	return 0
}

var GAMMA0 = M_PI / 2 //- 0.01
const GAMMA0_DEG = float64(90)

var phase1Gamma0 = GAMMA0
var PHASEROLL15 = deg2rad(15)
var PHASEROLL05 = deg2rad(5)

var phase2Gamma0 = GAMMA0 - PHASEROLL05

const PHASE1_ALTITUDE = 15000

var switch2phase2 = false

const PHASE2_ALTITUDE = 15000

type AscentPhase struct {
	startingTime     float64
	endingTime       float64
	startingAltitude float64
	endingAltitude   float64
	angleDeviation   float64
}

type AscentSet struct {
	currentPhase  int
	nphases       int
	deviationLeft float64
	aPhases       []AscentPhase
}

var asc AscentSet

// TwoPhaseTangentSteering ----------------------------------------------------
// steering program using time-to-MECO as first phase and distance-to-orbit as
// second phase
// ----------------------------------------------------------------------------
func (r *VEHICLE) TwoPhaseTangentSteering(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER || !r.SysGuidance._SEI1 {
				return 0
			}
		}

		// implements: tanθ(t)=tan(θ0) *(1 - altitude/orbitInsertion)
		if asc.currentPhase == 0 {
			// during first phase, we use time as a dividing factor (essentially time between PITCH and MECO)
			if r.Stages[i].Clock <= asc.aPhases[asc.currentPhase].endingTime {
				gamma := asc.deviationLeft - deg2rad(asc.aPhases[asc.currentPhase].angleDeviation)*((r.Stages[i].Clock-asc.aPhases[asc.currentPhase].startingTime)/(asc.aPhases[asc.currentPhase].endingTime-asc.aPhases[asc.currentPhase].startingTime))
				return gamma - r.Stages[i].gamma
			}
			// as we passed endTime, we enter the second phase where the dividing factor is altitude (essentially
			// between MECO altitude and target orbit altitude)
			asc.deviationLeft = asc.deviationLeft - deg2rad(asc.aPhases[asc.currentPhase].angleDeviation)
			asc.currentPhase += 1
			asc.aPhases[asc.currentPhase].startingAltitude = r.Stages[i].altitude
		}
		if r.Stages[i].altitude <= asc.aPhases[asc.currentPhase].endingAltitude {
			gamma := asc.deviationLeft - deg2rad(asc.aPhases[asc.currentPhase].angleDeviation)*((r.Stages[i].altitude-asc.aPhases[asc.currentPhase].startingAltitude)/(asc.aPhases[asc.currentPhase].endingAltitude-asc.aPhases[asc.currentPhase].startingAltitude))
			return gamma - r.Stages[i].gamma
		}
		return 0
	}
	return 0
}

// TwoPhaseTangentSteering ----------------------------------------------------
// steering program using multiple rate of changes based on altitude
// ----------------------------------------------------------------------------
func (r *VEHICLE) gravTurnMultiPhaseTangentSteering(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER { //|| !r.SysGuidance._SEI1 {
				return 0
			}
		}

		for true {
			if r.Stages[i].altitude <= asc.aPhases[asc.currentPhase].endingAltitude {
				gamma := asc.deviationLeft - deg2rad(asc.aPhases[asc.currentPhase].angleDeviation)*((r.Stages[i].altitude-asc.aPhases[asc.currentPhase].startingAltitude)/(asc.aPhases[asc.currentPhase].endingAltitude-asc.aPhases[asc.currentPhase].startingAltitude))
				if asc.currentPhase == len(asc.aPhases)-1 {
					//gamma = gamma + profile.InjectionAngle
				}
				return gamma - r.Stages[i].gamma
			} else {
				asc.deviationLeft = asc.deviationLeft - deg2rad(asc.aPhases[asc.currentPhase].angleDeviation)
				asc.currentPhase += 1
				if asc.currentPhase < asc.nphases {
					continue
				}
				asc.currentPhase -= 1
				break
			}
		}
	}
	return 0
}

// gravTurnSimpleTangentSteering ----------------------------------------------
// steering program using constant rate of change based on altitude
// ----------------------------------------------------------------------------
func (r *VEHICLE) gravTurnSimpleTangentSteering(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER || !r.SysGuidance._SEI1 {
				return 0
			}
		}
		// implements: tanθ(t)=tan(θ0) *(1 - altitude/orbitInsertion)
		// gamma := math.Atan(math.Tan(GAMMA0 * (1 - altitude/profile.OrbitInsertion)))
		if r.Stages[i].altitude <= profile.OrbitInsertion {
			gamma := GAMMA0 * (1 - (r.Stages[i].altitude)/(profile.OrbitInsertion))
			return gamma - r.Stages[i].gamma
		} else {
			r.Stages[i].gamma = 0
			return 0
		}
	}
	return 0
}

var switchPhase = false
var GammaPhase2 = 0.0

type AngleTime struct {
	Time  float64 `json:"time"`
	Angle float64 `json:"angle"`
}

var gammaSet = []AngleTime{{0, 90}, {1, 90}, {2, 90}, {3, 90}, {4, 90}, {5, 90}, {6, 90}, {7, 90}, {8, 90}, {9, 90}, {10, 90}, {11, 90}, {12, 90}, {13, 90}, {14, 90}, {15, 90}, {16, 90},
	{17, 90}, {18, 89.9}, {19, 89.8}, {20, 89.7}, {21, 89.6}, {22, 89.5}, {23, 89.4}, {24, 89.3}, {25, 89.2}, {26, 89.1}, {27, 89.0},
	//{28, 90}, {29, 90}, {30, 90},{31, 90}, {32, 90}, {33, 90}, {34, 90}, {35, 90}, {36, 90}, {37, 90},
	{28, 88.5}, {29, 88}, {30, 87}, {31, 86}, {32, 85}, {33, 84}, {34, 83}, {35, 82}, {36, 81}, {37, 80},

	{38, 79.092}, {39, 78.21}, {40, 77.5}, {41, 76.401}, {42, 75.777},
	{43, 75.129}, {44, 75.158}, {45, 75.028}, {46, 74.866}, {47, 74.631}, {48, 74.018}, {49, 73.323}, {50, 73.178}, {51, 72.853}, {52, 71.997},
	{53, 70.736}, {54, 69.397}, {55, 68.181}, {56, 67.221}, {57, 66.594}, {58, 66.607}, {59, 66.609}, {60, 66.395},

	{61, 65.764}, {62, 65.017}, {63, 64.052}, {64, 62.937}, {65, 61.85}, {66, 60.889}, {67, 60.007}, {68, 59.18}, {69, 58.383},
	{70, 57.679}, {71, 56.833}, {72, 56.145}, {73, 55.56}, {74, 54.997}, {75, 54.448}, {76, 53.896}, {77, 53.113}, {78, 52.304},
	{79, 51.386}, {80, 50.585}, {81, 49.614}, {82, 48.791}, {83, 48.052}, {84, 47.37}, {85, 46.741}, {86, 46.214}, {87, 45.605},
	{88, 45.037}, {89, 44.329}, {90, 43.524}, {91, 42.568}, {92, 41.52}, {93, 40.33}, {94, 39.085}, {95, 37.951}, {104, 36.893},
	{105, 36.604}, {106, 36.24}, {107, 35.908}, {108, 35.628}, {109, 35.409}, {110, 35.217}, {111, 34.992}, {112, 34.653},
	{113, 34.166}, {114, 33.533}, {115, 32.829}, {116, 32.133}, {117, 31.505}, {118, 30.949}, {119, 30.479}, {120, 30.176},
	{121, 30.011}, {122, 29.998}, {123, 30.071}, {124, 30.105}, {125, 30.211}, {126, 30.285}, {127, 30.31}, {128, 30.244},
	{129, 30.102}, {130, 29.833}, {131, 29.574}, {132, 29.355}, {133, 29.152}, {134, 28.976}, {135, 28.848}, {136, 28.799},
	{137, 28.802}, {138, 28.817}, {139, 28.84}, {140, 28.852}, {141, 28.853}, {142, 28.845}, {143, 28.785}, {144, 28.658},
	{145, 28.416}, {146, 28.139}, {147, 27.834}, {148, 27.545}, {149, 27.253}, {150, 26.974}, {151, 26.713}, {152, 26.476},
	{153, 26.24}, {154, 25.991}, {155, 25.719},
}

var curAngleIndex, nextAngleIndex = 0, 1
var angleIncrement = tinc
var tlast = 0

// gravTurnProgrammedTwoPhaseSteering -----------------------------------------
// steering program using pre-define rate of change in first phase, and rate of
// change based on altitude in second phase
// ----------------------------------------------------------------------------
func (r *VEHICLE) gravTurnProgrammedTwoPhaseSteering(i int32) float64 {
	if r.SysGuidance._pitch {
		// after stage sep, we don't care about booster gravity turn
		if r.SysGuidance._stagesep {
			// if we have a stage separation, we don't care of the booster, but if it is the
			// second stage, make sure to have engine ignition before continuing the steering program
			if i == BOOSTER || !r.SysGuidance._SEI1 {
				return 0
			}
		}
		if r.Stages[i].Clock <= 155.0 {
			// first phase
			tcur := math.Round(r.Stages[i].Clock)
			if int(tcur) < len(gammaSet)-1 {
				gamma := deg2rad(gammaSet[int(tcur)].Angle)
				gammaNext := deg2rad(gammaSet[int(tcur+1)].Angle)
				r.Stages[i].gamma = r.Stages[i].gamma - ((gamma - gammaNext) * tinc)
			}
			return 0

			//return deg2rad(90) - gamma // - r.Stages[i].gamma
		} else {

			if !switchPhase {
				switchPhase = true
				// init initial altitude and angle and with current values
				GammaPhase2 = r.Stages[i].gamma
			}
			// second phase
			if r.Stages[i].altitude <= profile.OrbitInsertion {
				gamma := GammaPhase2 * (1 - (r.Stages[i].altitude)/(profile.OrbitInsertion))
				return gamma - r.Stages[i].gamma
			} else {
				r.Stages[i].gamma = 0
				return 0
			}
		}
	}
	return 0
}

func (r *VEHICLE) landingTimeStep() { // i = stage

	dgamma := float64(0.0)

	if r.SysGuidance._release {
		aerodynPressure = 0.5 * rho(r.Stages[BOOSTER].DTF-Re) * r.Stages[BOOSTER].VRelative * r.Stages[BOOSTER].VRelative * 1e-3                                   // Aerodynamic stress
		drag = (0.5) * r.Stages[BOOSTER].Cd * r.Stages[BOOSTER].CSArea * rho(r.Stages[BOOSTER].DTF-Re) * r.Stages[BOOSTER].VRelative * r.Stages[BOOSTER].VRelative // Drag
		r.Stages[BOOSTER].Thrust = float64(r.Stages[BOOSTER].RunningEngines) * r.Stages[BOOSTER].ThrottleRate * r.GetThrust(r.Stages[BOOSTER].DTF, BOOSTER)        // Thrust

		// calculate force and velocity vectors norm
		gh := g(r.Stages[BOOSTER].DTF)
		r.Stages[BOOSTER].Force = (r.Stages[BOOSTER].Thrust + drag) - r.Stages[BOOSTER].Mass*gh*math.Sin(r.Stages[BOOSTER].gamma)
		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].AVel*math.Cos(r.Stages[i].gamma), 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)

		//		r.Stages[i].Force = (r.Stages[i].Thrust - drag) - r.Stages[i].Mass*(gh-math.Pow(r.Stages[i].vAx, 2)/(Re+r.Stages[i].altitude))*math.Sin(r.Stages[i].gamma)
		r.Stages[BOOSTER].ax = r.Stages[BOOSTER].Force * math.Cos(r.Stages[BOOSTER].gamma) / r.Stages[BOOSTER].Mass

		// update downrange first dx/dt = (Re/(Re+h))*v*cos(gamma)
		// since we are using altitude in the range calculation formula, calculate range first so that both
		// altitude and range are calculated simultaneously
		// first update range as a fonction of old relative speed, old altitude and "dt" time increment
		r.Stages[BOOSTER].drange = (Re/(Re+r.Stages[BOOSTER].altitude))*r.Stages[BOOSTER].RVel*math.Cos(r.Stages[BOOSTER].gamma)*r.Stages[BOOSTER].dt + r.Stages[BOOSTER].drange

		// second update altitude: dh/dt = v.sin(gamma) as a function of old relative speed and time increment "dt"
		r.Stages[BOOSTER].altitude = r.Stages[BOOSTER].RVel*math.Sin(r.Stages[BOOSTER].gamma)*r.Stages[BOOSTER].dt + r.Stages[BOOSTER].altitude

		// calculate dgamma/dt = -(g - v*v/(Re+h)) * cos(gamma) * 1/v
		//dgamma = -(gh - (r.Stages[i].Vel * r.Stages[i].Vel)/(Re + r.Stages[i].altitude)) * math.Cos(r.Stages[i].gamma) * (1/r.Stages[i].Vel)
		// dgamma = gamma /()
		//		dgamma = (profile.OrbitPerigee - r.Stages[i].altitude)*1e-3 / (profile.BurnoutTime - r.Stages[i].Clock)

		//r.TotalTimeIncrement = qTotalTimeIncrement + dt

		//		dgamma = r.gravTurnTangentSteering(i)

		// -----------> dgamma = r.gravTurnMultiPhaseTangentSteering(BOOSTER)

		//dgamma = r.gravTurnClassic(i)
		//if dgamma != 0 {
		//	println(dgamma)
		//}

		//		dgamma = r.gravTurn(i)

		// calculate cartesian coordinates
		r.Stages[BOOSTER].beta = r.Stages[BOOSTER].drange / Re // polar angle (in rd) based on downrange value
		r.Stages[BOOSTER].px = (Re + r.Stages[BOOSTER].altitude) * math.Sin(r.Stages[BOOSTER].beta)
		r.Stages[BOOSTER].py = (Re + r.Stages[BOOSTER].altitude) * math.Cos(r.Stages[BOOSTER].beta)

		if false {
			/*			i := 0
						// x-direction
						//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
						//		r.Stages[i].ForceX = r.Stages[i].Force * math.Cos(r.Stages[i].gamma)
						//		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
						r.Stages[i].ax = r.Stages[i].Force * math.Cos(r.Stages[i].gamma) / r.Stages[i].Mass

						r.Stages[i].vRx = r.Stages[i].vRx + r.Stages[i].ax*r.Stages[i].dt
						r.Stages[i].vAx = r.Stages[i].vRx + vE //* math.Sin(r.Stages[i].beta)

						// y-direction
						//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].beta) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
						//		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
						r.Stages[i].ay = r.Stages[i].Force * math.Sin(r.Stages[i].gamma) / r.Stages[i].Mass

						r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay*r.Stages[i].dt
						r.Stages[i].vRy = r.Stages[i].vAy //- vE * math.Cos(M_PI + r.Stages[i].beta)

						r.Stages[i].DTF = math.Sqrt(r.Stages[i].px*r.Stages[i].px + r.Stages[i].py*r.Stages[i].py)
						r.Stages[i].VAbsolute = math.Sqrt(r.Stages[i].vAx*r.Stages[i].vAx + r.Stages[i].vAy*r.Stages[i].vAy)
						r.Stages[i].VRelative = math.Sqrt(r.Stages[i].vRx*r.Stages[i].vRx + r.Stages[i].vRy*r.Stages[i].vRy)
						r.Stages[i].Acc = r.Stages[i].Force / r.Stages[i].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)
			*/
		} else {
			// NEW NEW
			r.Stages[BOOSTER].VAbsolute = r.Stages[BOOSTER].AVel
			r.Stages[BOOSTER].VRelative = r.Stages[BOOSTER].RVel
			r.Stages[BOOSTER].DTF = Re + r.Stages[BOOSTER].altitude
			r.Stages[BOOSTER].Acc = r.Stages[BOOSTER].Force / r.Stages[BOOSTER].Mass //math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)
		}

		// update velocities
		r.Stages[BOOSTER].RVel = r.Stages[BOOSTER].RVel + (r.Stages[BOOSTER].Force/r.Stages[BOOSTER].Mass)*r.Stages[BOOSTER].dt
		r.Stages[BOOSTER].AVel = r.Stages[BOOSTER].RVel + vE
		r.Stages[BOOSTER].vRx = r.Stages[BOOSTER].vRx + r.Stages[BOOSTER].ax*r.Stages[BOOSTER].dt
		r.Stages[BOOSTER].vAx = r.Stages[BOOSTER].vRx + vE

		r.Stages[BOOSTER].VAbsolute = r.Stages[BOOSTER].AVel
		r.Stages[BOOSTER].VRelative = r.Stages[BOOSTER].RVel

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
		r.Stages[i].px = r.Stages[i].px + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].py = r.Stages[i].py + r.Stages[i].vAy * r.Stages[i].dt
	*/

	if r.Stages[BOOSTER].gamma >= math.Abs(dgamma) {
		r.Stages[BOOSTER].gamma = r.Stages[BOOSTER].gamma + dgamma //- math.Abs(dgamma)
	} else {
		r.Stages[BOOSTER].gamma = 0
	}

	if r.Stages[BOOSTER].RunningEngines > 0 {
		dm = float64(r.Stages[BOOSTER].RunningEngines) * r.Stages[BOOSTER].ThrottleRate * EnginesMap[r.Stages[BOOSTER].EngineID].Flow_rate * r.Stages[BOOSTER].dt
		r.Stages[BOOSTER].Mf = r.Stages[BOOSTER].Mf - dm
		r.Stages[BOOSTER].Mass = r.Stages[BOOSTER].Mass - dm
	}

	//	fmt.Println("gamma = ", rad2deg(r.Stages[i].gamma))
}
