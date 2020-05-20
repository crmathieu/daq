package main
import (
	"math"
)

//inline void pitch_kick(int *x)
//func pitch_kick(x *int32) {
//	gam[0] = (M_PI / 2) - 0.025;
//	*x = 1;
//}
func pitch_kick() bool {
	F9.Stages[BOOSTER].gam = (M_PI / 2) - 0.025
	return true
}

//inline void angles(int i)
func angles(i int32) {
/*
	beta is the angle through which gravity pulls the vehicle. alpha is the angle of attack relative to earth.
	The 'if' statements are just for trigonometry. If you draw a picture you'll see why they're necessary.
*/

	F9.Stages[i].beta = math.Acos(F9.Stages[i].cx / F9.Stages[i].PolarDistance)
	if F9.Stages[i].cy < 0 {
		F9.Stages[i].beta = (2 * M_PI) - F9.Stages[i].beta
	}

	F9.Stages[i].alpha = math.Acos(F9.Stages[i].vRx / F9.Stages[i].VRelative)
	if F9.Stages[i].vRy < 0 {
		F9.Stages[i].alpha = (2 * M_PI) - F9.Stages[i].alpha
	}
}

func grav_turn(i int32) {
	F9.Stages[i].gam = F9.Stages[i].alpha
	if _MECO1 && math.Cos(F9.Stages[i].beta - F9.Stages[i].alpha) < 0 && !_SECO1 {
		F9.Stages[i].gam = math.Asin(-F9.Stages[i].Mass * g(F9.Stages[i].PolarDistance) * math.Sin(F9.Stages[i].beta + M_PI) / F9.Stages[i].Thrust)
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
	if F9.Stages[STAGE2].VRelative > 2200 {
		F9.Stages[STAGE2].ThrottleRate = 0.7				// Throttle down to 70%
	}
	if F9.Stages[STAGE2].VRelative > 2900 {
		F9.Stages[STAGE2].gam = F9.Stages[STAGE2].beta - M_PI/2			// go horizontal rel. to earth
	}
	if F9.Stages[STAGE2].VRelative > 3400 {
		F9.Stages[STAGE2].gam = F9.Stages[STAGE2].beta - M_PI/2 - 0.1		// Start pointing down....
	}
	if F9.Stages[STAGE2].VRelative > 3900 {
		F9.Stages[STAGE2].gam = F9.Stages[STAGE2].beta - M_PI/2 - 0.2
	}
	if F9.Stages[STAGE2].VRelative > 4400 {
		F9.Stages[STAGE2].gam = F9.Stages[STAGE2].beta - M_PI/2 - 0.3
	}
	if F9.Stages[STAGE2].VRelative > 6200 {
		F9.Stages[STAGE2].gam = F9.Stages[STAGE2].beta - M_PI/2 - 0.4
	}
}

//inline void first_step()
func first_step() bool {
	Engine := EnginesMap[F9.Stages[BOOSTER].EngineID]

	F9.Stages[BOOSTER].ForceY = (float64(F9.Stages[BOOSTER].RunningEngines) * Engine.Th_sl) - F9.Stages[BOOSTER].Mass * g0
	F9.Stages[BOOSTER].ay = F9.Stages[BOOSTER].ForceY / F9.Stages[BOOSTER].Mass

	F9.Stages[BOOSTER].vAx = vE				// Absolute velocity in x-direction = velocity of earth at surface
	F9.Stages[BOOSTER].vAy = F9.Stages[BOOSTER].vAy + (F9.Stages[BOOSTER].ay * F9.Stages[BOOSTER].dt/2)

	F9.Stages[BOOSTER].PolarDistance = math.Sqrt(F9.Stages[BOOSTER].cx * F9.Stages[BOOSTER].cx + F9.Stages[BOOSTER].cy * F9.Stages[BOOSTER].cy)
	F9.Stages[BOOSTER].Acc = math.Sqrt(F9.Stages[BOOSTER].ax * F9.Stages[BOOSTER].ax + F9.Stages[BOOSTER].ay * F9.Stages[BOOSTER].ay)
	return true
}

/*
	Leapfrog integrator for moving ma boi F9. Conservative which is nice. Euler can go poo a pineapple.
*/

//void leapfrog_step(int i) // i = stage
func leapfrog_step(i int32) { // i = stage


	if _MEI1 {
		dm = float64(F9.Stages[i].RunningEngines) * F9.Stages[i].ThrottleRate * 236 * F9.Stages[i].dt
		F9.Stages[i].Mf = F9.Stages[i].Mf - dm;
		F9.Stages[i].Mass = F9.Stages[i].Mass - dm;
	}

	if _release {
		aerodynPressure = 0.5 * rho(F9.Stages[i].PolarDistance - Re) * F9.Stages[i].VRelative * F9.Stages[i].VRelative * 1e-3				// Aerodynamic stress
		drag = (0.5) * F9.Stages[i].Cd * F9.Stages[i].CSArea * rho(F9.Stages[i].PolarDistance - Re) * F9.Stages[i].VRelative * F9.Stages[i].VRelative		// Drag
		F9.Stages[i].Thrust = float64(F9.Stages[i].RunningEngines) * F9.Stages[i].ThrottleRate * GetThrust(F9.Stages[i].PolarDistance, i)	// Thrust

		/* x-direction	*/
		F9.Stages[i].ForceX = F9.Stages[i].Thrust * math.Cos(F9.Stages[i].gam) + drag * math.Cos(F9.Stages[i].alpha + M_PI) + F9.Stages[i].Mass * g(F9.Stages[i].PolarDistance) * math.Cos(F9.Stages[i].beta + M_PI)
		F9.Stages[i].cx = F9.Stages[i].cx + F9.Stages[i].vAx * F9.Stages[i].dt
		F9.Stages[i].ax = F9.Stages[i].ForceX / F9.Stages[i].Mass

		F9.Stages[i].vAx = F9.Stages[i].vAx + F9.Stages[i].ax * F9.Stages[i].dt
		F9.Stages[i].vRx = F9.Stages[i].vAx - vE * math.Sin(F9.Stages[i].beta)

		/* y-direction	*/
		F9.Stages[i].ForceY = F9.Stages[i].Thrust * math.Sin(F9.Stages[i].gam) + drag * math.Sin(F9.Stages[i].alpha + M_PI) + F9.Stages[i].Mass * g(F9.Stages[i].PolarDistance) * math.Sin(F9.Stages[i].beta + M_PI)
		F9.Stages[i].cy = F9.Stages[i].cy + F9.Stages[i].vAy * F9.Stages[i].dt
		F9.Stages[i].ay = F9.Stages[i].ForceY / F9.Stages[i].Mass

		F9.Stages[i].vAy = F9.Stages[i].vAy + F9.Stages[i].ay * F9.Stages[i].dt
		F9.Stages[i].vRy = F9.Stages[i].vAy - vE * math.Cos(M_PI + F9.Stages[i].beta)

		F9.Stages[i].PolarDistance = math.Sqrt(F9.Stages[i].cx * F9.Stages[i].cx + F9.Stages[i].cy * F9.Stages[i].cy)
		F9.Stages[i].VAbsolute = math.Sqrt(F9.Stages[i].vAx * F9.Stages[i].vAx + F9.Stages[i].vAy * F9.Stages[i].vAy)
		F9.Stages[i].VRelative = math.Sqrt(F9.Stages[i].vRx * F9.Stages[i].vRx + F9.Stages[i].vRy * F9.Stages[i].vRy)
		F9.Stages[i].Acc = math.Sqrt(F9.Stages[i].ax * F9.Stages[i].ax + F9.Stages[i].ay * F9.Stages[i].ay)
	}

	if _release {
		angles(i)
	}	
	if F9.Stages[i].Clock > 55 {
		if i == STAGE2 {
			grav_turn(i);
		} 	
		if i == BOOSTER && !_MECO1 {
			grav_turn(i)
		}	
	}
	if _BBURN  || _LBURN {
		flip(0)
	}
	if _LBURN && mod(F9.Stages[i].Clock, 5) < F9.Stages[i].dt {
		update_landing_throttle()
	}
}


