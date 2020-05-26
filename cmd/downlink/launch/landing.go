
package main
import (
	"math"
)

/*
	This function tests a throttle value to find out, at a certain throttle, what
	height will the vehicle reach zero vertical velocity.
*/

//func (r *VEHICLE) throttle_test(hy, ux, uy, mf, throttle float64) float64 {
func (r *VEHICLE) throttle_test(throttle float64) float64 {
// r.Stages[BOOSTER].DTF, r.Stages[BOOSTER].vRx, r.Stages[BOOSTER].vRy, r.Stages[BOOSTER].Mf
//	hy                                ux                    uy                       mf
	var hy, ux, uy, mf float64
	var VEL, ft, fd, mass float64
	var fx, fy, ax, ay float64

	hy = r.Stages[BOOSTER].DTF
	ux = r.Stages[BOOSTER].vRx
	uy = r.Stages[BOOSTER].vRy
	mf = r.Stages[BOOSTER].Mf
	VEL = math.Sqrt(float64(ux*ux + uy*uy))

	ft = throttle * r.GetThrust(hy, 0);

	//do
	for ;uy < 0; {
		mf = mf - throttle * 236 * r.Stages[BOOSTER].dt
		mass = mf + r.Stages[BOOSTER].Mr;						

		fd = (0.5) * r.Stages[BOOSTER].Cd * r.Stages[BOOSTER].CSArea * rho(hy - Re) * VEL * VEL

		r.flip(0);

		fx = ft * math.Cos(r.Stages[BOOSTER].gamma) + fd * math.Cos(r.Stages[BOOSTER].alpha + M_PI) + mass * g(hy) * math.Cos(r.Stages[BOOSTER].beta + M_PI)
		ax = fx / mass
		ux = ux + ax * r.Stages[BOOSTER].dt
	
		fy = ft * math.Sin(r.Stages[BOOSTER].gamma) + fd * math.Sin(r.Stages[BOOSTER].alpha + M_PI) + mass * g(hy) * math.Sin(r.Stages[BOOSTER].beta + M_PI)
		hy = hy + uy * r.Stages[BOOSTER].dt
		ay = fy / mass
		uy = uy + ay * r.Stages[BOOSTER].dt

		VEL = math.Sqrt((ux * ux) + (uy * uy))
	}
	return hy - Re
}

/*
	This guy calls the above function to get as close to a hoverslam as possible. Need to keep calculating this
	as we fall cause it changes with rounding errors etc. I re-calculate it every 5 seconds, gets me a pretty good result
*/

//inline double get_landing_throttle(double H, double ux, double uy, double mf)
//func (r *VEHICLE) get_landing_throttle(H, ux, uy, mf float64) float64 {
func (r *VEHICLE) get_landing_throttle() float64 {

	var a = float64(0.7) 
	var b = float64(1.0)
	var end_H float64

	if r.throttle_test(b) > 0 {			// Will a full-power burn keep you alive?
		if r.throttle_test(a) < 0 {			// Yes. Will a minimum-power burn kill you?
			for {
				end_H = r.throttle_test((a+b)/2.0)
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
func (r *VEHICLE) update_landing_throttle() {
	r.Stages[BOOSTER].ThrottleRate = r.get_landing_throttle()
}

