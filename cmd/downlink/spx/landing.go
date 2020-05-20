
package main
import (
	"math"
)

/*
	This function tests a throttle value to find out, at a certain throttle, what
	height will the vehicle reach zero vertical velocity.
*/

func throttle_test(hy, ux, uy, mf, throttle float64) float64 {

	var VEL, ft, fd, mass float64
	var fx, fy, ax, ay float64

	VEL = math.Sqrt(float64(ux*ux + uy*uy))

	ft = throttle * GetThrust(hy, 0);

	//do
	for ;uy < 0; {
		mf = mf - throttle * 236 * F9.Stages[BOOSTER].dt
		mass = mf + F9.Stages[BOOSTER].Mr;						

		fd = (0.5) * F9.Stages[BOOSTER].Cd * F9.Stages[BOOSTER].CSArea * rho(hy - Re) * VEL * VEL

		flip(0);

		fx = ft * math.Cos(F9.Stages[BOOSTER].gam) + fd * math.Cos(F9.Stages[BOOSTER].alpha + M_PI) + mass * g(hy) * math.Cos(F9.Stages[BOOSTER].beta + M_PI)
		ax = fx / mass
		ux = ux + ax * F9.Stages[BOOSTER].dt
	
		fy = ft * math.Sin(F9.Stages[BOOSTER].gam) + fd * math.Sin(F9.Stages[BOOSTER].alpha + M_PI) + mass * g(hy) * math.Sin(F9.Stages[BOOSTER].beta + M_PI)
		hy = hy + uy * F9.Stages[BOOSTER].dt
		ay = fy / mass
		uy = uy + ay * F9.Stages[BOOSTER].dt

		VEL = math.Sqrt((ux * ux) + (uy * uy))
	}
	return hy - Re
}

/*
	This guy calls the above function to get as close to a hoverslam as possible. Need to keep calculating this
	as we fall cause it changes with rounding errors etc. I re-calculate it every 5 seconds, gets me a pretty good result
*/

//inline double get_landing_throttle(double H, double ux, double uy, double mf)
func get_landing_throttle(H, ux, uy, mf float64) float64 {

	var a = float64(0.7) 
	var b = float64(1.0)
	var end_H float64

	if throttle_test(H, ux, uy, mf, b) > 0 {			// Will a full-power burn keep you alive?
		if throttle_test(H, ux, uy, mf, a) < 0 {			// Yes. Will a minimum-power burn kill you?
			for {
				end_H = throttle_test(H, ux, uy, mf, (a+b)/2.0)
				if math.Abs(end_H) < 0.1 {
					return (a + b) / 2.0				// Yes. Burn at this throttle from now to do hoverslam.
				}
				if end_H < 0 {
					a = (a + b) / 2.0
				} 
				if end_H > 0 {
					b = (a + b) / 2.0
				}
			}
		} else {
			return 0.0;						// No. Don't start burn yet. 
		}
	} else {
		return 1.0;						// No. Too late. Crash unavoidable. Should have started earlier
	}

}

//inline void update_landing_throttle()
func update_landing_throttle() {
	F9.Stages[BOOSTER].ThrottleRate = get_landing_throttle(F9.Stages[BOOSTER].PolarDistance, F9.Stages[BOOSTER].vRx, F9.Stages[BOOSTER].vRy, F9.Stages[BOOSTER].Mf)
}

