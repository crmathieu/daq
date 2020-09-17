package main
import (
	"math"
	"github.com/crmathieu/daq/packages/data"
)


//inline void angles(int i)
func (r *VEHICLE) angles(i int32) {
/*
	beta is the angle through which gravity pulls the vehicle. alpha is the angle of attack relative to earth.
	The 'if' statements are just for trigonometry. If you draw a picture you'll see why they're necessary.
*/

	r.Stages[i].beta = math.Acos(r.Stages[i].cx / r.Stages[i].DTF)
	if r.Stages[i].cy < 0 {
		r.Stages[i].beta = (2 * M_PI) - r.Stages[i].beta
	}

	r.Stages[i].alpha = math.Acos(r.Stages[i].vRx / r.Stages[i].VRelative)
	if r.Stages[i].vRy < 0 {
		r.Stages[i].alpha = (2 * M_PI) - r.Stages[i].alpha
	}
}

func (r *VEHICLE) grav_turn(i int32) {
	r.Stages[i].gamma = r.Stages[i].alpha
	if r.SysGuidance._MECO1 && math.Cos(r.Stages[i].beta - r.Stages[i].alpha) < 0 && !r.SysGuidance._SECO1 {
		r.Stages[i].gamma = math.Asin(-r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta + M_PI) / r.Stages[i].Thrust)
	}
/*
	The line above ensures that if the gravity turn makes the 2nd stage point down towards earth,
	(which should happen when the velocity is also tangential to earth)
	it will correct the angle and point such that it's upward thrust cancels out
	the downward force of gravity - theoretically keeping it at constant altitude
	until it reaches orbital velocity. Should never be used if gravity turn executed
	properly though, as it's super inefficient.
*/
/*
	The next few lines are OG2 course corrections
	Trial and error. Lots of corrections for ultra steep trajectory.
*/

/*	if r.Stages[STAGE2].VRelative > 2200 {
		r.Stages[STAGE2].ThrottleRate = 0.7				// Throttle down to 70%
	}
*/	if r.Stages[STAGE2].VRelative > 2900 {
		r.Stages[STAGE2].gamma = r.Stages[STAGE2].beta - M_PI/2			// go horizontal rel. to earth
	}
	if r.Stages[STAGE2].VRelative > 3400 {
		r.Stages[STAGE2].gamma = r.Stages[STAGE2].beta - M_PI/2 - 0.1		// Start pointing down....
	}
	if r.Stages[STAGE2].VRelative > 3900 {
		r.Stages[STAGE2].gamma = r.Stages[STAGE2].beta - M_PI/2 - 0.2
	}
	if r.Stages[STAGE2].VRelative > 4400 {
		r.Stages[STAGE2].gamma = r.Stages[STAGE2].beta - M_PI/2 - 0.3
	}
	if r.Stages[STAGE2].VRelative > 6200 {
		r.Stages[STAGE2].gamma = r.Stages[STAGE2].beta - M_PI/2 - 0.4
	}
}

//inline void liftOff()
func (r *VEHICLE) liftOff() bool {
	Engine := EnginesMap[r.Stages[BOOSTER].EngineID]

	r.Stages[BOOSTER].ForceY = (float64(r.Stages[BOOSTER].RunningEngines) * Engine.Th_sl) - r.Stages[BOOSTER].Mass * g0
	r.Stages[BOOSTER].ay = r.Stages[BOOSTER].ForceY / r.Stages[BOOSTER].Mass

	r.Stages[BOOSTER].vAx = vE				// Absolute velocity in x-direction = velocity of earth at surface
	r.Stages[BOOSTER].vAy = r.Stages[BOOSTER].vAy + (r.Stages[BOOSTER].ay * r.Stages[BOOSTER].dt/2)

	r.Stages[BOOSTER].DTF = math.Sqrt(r.Stages[BOOSTER].cx * r.Stages[BOOSTER].cx + r.Stages[BOOSTER].cy * r.Stages[BOOSTER].cy)
	r.Stages[BOOSTER].Acc = math.Sqrt(r.Stages[BOOSTER].ax * r.Stages[BOOSTER].ax + r.Stages[BOOSTER].ay * r.Stages[BOOSTER].ay)
	
	r.EventsMap = r.EventsMap | data.E_LIFTOFF
	r.LastEvent = data.E_LIFTOFF

	return true
}

/*
	Leapfrog integrator for moving ma boi r. Conservative which is nice. Euler can go poo a pineapple.
*/

//void timeStep(int i) // i = stage
func (r *VEHICLE) timeStep(i int32) { // i = stage

	if r.SysGuidance._MEI1 {
		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * 236 * r.Stages[i].dt
//		dm = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * EnginesMap[r.Stages[i].EngineID].Flow_rate * r.Stages[i].dt
		r.Stages[i].Mf = r.Stages[i].Mf - dm;
		r.Stages[i].Mass = r.Stages[i].Mass - dm;
	}

	if r.SysGuidance._release {
		aerodynPressure = 0.5 * rho(r.Stages[i].DTF - Re) * r.Stages[i].VRelative * r.Stages[i].VRelative * 1e-3				// Aerodynamic stress
		drag = (0.5) * r.Stages[i].Cd * r.Stages[i].CSArea * rho(r.Stages[i].DTF - Re) * r.Stages[i].VRelative * r.Stages[i].VRelative		// Drag
		r.Stages[i].Thrust = float64(r.Stages[i].RunningEngines) * r.Stages[i].ThrottleRate * r.GetThrust(r.Stages[i].DTF, i)	// Thrust

		/* x-direction	*/
		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) + drag * math.Cos(r.Stages[i].alpha + M_PI) + r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta + M_PI)
//		r.Stages[i].ForceX = r.Stages[i].Thrust * math.Cos(r.Stages[i].gamma) - drag * math.Cos(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Cos(r.Stages[i].beta)
		r.Stages[i].cx = r.Stages[i].cx + r.Stages[i].vAx * r.Stages[i].dt
		r.Stages[i].ax = r.Stages[i].ForceX / r.Stages[i].Mass

		r.Stages[i].vAx = r.Stages[i].vAx + r.Stages[i].ax * r.Stages[i].dt
		r.Stages[i].vRx = r.Stages[i].vAx - vE * math.Sin(r.Stages[i].beta)

		/* y-direction	*/
		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) + drag * math.Sin(r.Stages[i].alpha + M_PI) + r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta + M_PI)
//		r.Stages[i].ForceY = r.Stages[i].Thrust * math.Sin(r.Stages[i].gamma) - drag * math.Sin(r.Stages[i].alpha) - r.Stages[i].Mass * g(r.Stages[i].DTF) * math.Sin(r.Stages[i].beta)
		r.Stages[i].cy = r.Stages[i].cy + r.Stages[i].vAy * r.Stages[i].dt
		r.Stages[i].ay = r.Stages[i].ForceY / r.Stages[i].Mass

		r.Stages[i].vAy = r.Stages[i].vAy + r.Stages[i].ay * r.Stages[i].dt
		r.Stages[i].vRy = r.Stages[i].vAy - vE * math.Cos(M_PI + r.Stages[i].beta)

		r.Stages[i].DTF = math.Sqrt(r.Stages[i].cx * r.Stages[i].cx + r.Stages[i].cy * r.Stages[i].cy)
		r.Stages[i].VAbsolute = math.Sqrt(r.Stages[i].vAx * r.Stages[i].vAx + r.Stages[i].vAy * r.Stages[i].vAy)
		r.Stages[i].VRelative = math.Sqrt(r.Stages[i].vRx * r.Stages[i].vRx + r.Stages[i].vRy * r.Stages[i].vRy)
		r.Stages[i].Acc = math.Sqrt(r.Stages[i].ax * r.Stages[i].ax + r.Stages[i].ay * r.Stages[i].ay)
	}

	if r.SysGuidance._release {
		r.angles(i)
	}	
	if r.Stages[i].Clock > 55 {
		if i == STAGE2 {
			r.grav_turn(i);
		} 	
		if i == BOOSTER && !r.SysGuidance._MECO1 {
			r.grav_turn(i)
		}	
	}
	if r.SysGuidance._BBURN  || r.SysGuidance._LBURN {
		r.flip(0)
	}
	if r.SysGuidance._LBURN && mod(r.Stages[i].Clock, 5) < r.Stages[i].dt {
		r.update_landing_throttle()
	}
}


