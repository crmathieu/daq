package main
import (
	"math"
)

func mod(a, b float64) float64 {
	for a >= b {
		a = a - b
	}
	return a
	/*	if a < b {
		return a
	}
	return mod(a - b, b)*/
}

func g(h float64) float64 {
	return G * Me/(h * h)
}

/*
	These next two functions for atmospheric density and pressure were obtained from a load of data points
	I found on some NASA website, to which I fitted functions using some other awesome plot-fitting website.
*/

//inline double rho(double h)
func rho(h float64) float64 {
	return 1.21147 * math.Exp(h * -1.12727e-4)
}

func P(h float64) float64 {
	return -517.18 * math.Log(0.012833 * math.Log(6.0985e28 * h + 2.0981e28))
}

/*
	Interpolating Isp at given altitude using sea-level/vacuum values, and the current atmospheric pressure.
*/


func (v *VEHICLE) Isp(h float64) float64 {
	Engine := EnginesMap[v.Stages[BOOSTER].EngineID]
	if h < 80000 {
		return Engine.Isp_sl + (1.0 / P(0)) * (P(0) - P(h * 1e-3)) * (Engine.Isp_vac - Engine.Isp_sl)
	}
	return Engine.Isp_vac
}

func (v *VEHICLE) GetThrust(H float64, stage int32) float64 {
	if stage == BOOSTER || !v.SysGuidance._MECO1 {
		return v.Isp(H - Re) * 236 * g0	// 236 kg/s = M1D rate of fuel consumption
//		return v.Isp(H - Re) * EnginesMap[v.Stages[BOOSTER].EngineID].Flow_rate * g0
	} 
	return EnginesMap[v.Stages[STAGE2].EngineID].Th_vac    //M1Dv.Th_vac
}

func (r *VEHICLE) flip(i int32) {
	r.Stages[i].gamma = r.Stages[i].alpha + math.Pi		// retrograde
}


