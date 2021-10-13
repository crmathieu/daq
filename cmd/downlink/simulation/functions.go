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
	return G * Me / (h * h)
}

/*
	These next two functions for atmospheric density and pressure were obtained from a load of data points
	I found on some NASA website, to which I fitted functions using some other awesome plot-fitting website.
*/

//inline double rho(double h)
func (r *VEHICLE) UpdateDensity(i int32) { //}, h float64) {
	if r.Stages[i].altitude <= 11000 {
		//Temp(l+1) = 15.04-0.00649*y(l);
		r.Stages[i].t_temp = 15.04 - 0.00649*r.Stages[i].altitude
		//p(l+1) = 101.29*((Temp(l)+273.1)/288.08)^5.258;
		r.Stages[i].t_atmoPressure = 101.29 * math.Pow(((r.Stages[i].temp+273.1)/288.08), 5.258)
	} else if 11000 < r.Stages[i].altitude && r.Stages[i].altitude <= 25000 {
		//Temp(l+1) = -56.46;
		r.Stages[i].t_temp = -56.46
		//p(l+1) = 22.65*exp(1.73-0.000157*y(l));
		r.Stages[i].t_atmoPressure = 22.65 * math.Exp(1.73-0.000157*r.Stages[i].altitude)
	} else {
		//Temp(l+1) = -131.21 + 0.00299*y(l);
		r.Stages[i].t_temp = -131.21 + 0.00299*r.Stages[i].altitude
		//p(l+1) = 2.488*((Temp(l)+273.1)/216.6)^-11.388;
		r.Stages[i].t_atmoPressure = 2.488 * math.Pow(((r.Stages[i].temp+273.1)/216.6), -11.388)
	}
	r.Stages[i].t_density = r.Stages[i].atmoPressure / (0.2869 * (r.Stages[i].temp + 273.1))
}

func (r *VEHICLE) UpdateDrag(i int32) {

	//    D(l+1) = 0.5*rho(l)*v(l)^2*Aw*cd;
	r.Stages[i].t_drag = 0.5 * r.Stages[i].density * math.Pow(r.Stages[i].rvelocity, 2) * r.Stages[i].Cd * r.Stages[i].CSArea
}

func rho(h float64) float64 {
	/*	var hscale = 7.5e3
		return 1.225 * math.Exp(-h/hscale)
		//return 1.21147 * math.Exp(h*-1.12727e-4)
	*/
	Temp := 0.0
	if h < 11000 {
		Temp = (15 - 0.0065*h) + 273
	} else {
		Temp = -56.5 + 273
	}
	if h < 100000 {
		return 1.2260 * math.Exp(-1.e-5*1.2260*9.81*h*(14.4+273)/Temp)
	}
	return 0
}

func P(h float64) float64 {
	return -517.18 * math.Log(0.012833*math.Log(6.0985e28*h+2.0981e28))
}

/*
	Interpolating Isp at given altitude using sea-level/vacuum values, and the current atmospheric pressure.
*/
func (v *VEHICLE) Isp(i int32) float64 {
	Engine := EnginesMap[v.Stages[i].EngineID]
	if v.Stages[i].altitude < 80000 {
		return Engine.Isp_sl + (1.0/P(0))*(P(0)-P(v.Stages[i].altitude*1e-3))*(Engine.Isp_vac-Engine.Isp_sl)
	}
	return Engine.Isp_vac
}

func (v *VEHICLE) Isp2(h float64) float64 {
	Engine := EnginesMap[v.Stages[BOOSTER].EngineID]
	if h < 80000 {
		return Engine.Isp_sl + (1.0/P(0))*(P(0)-P(h*1e-3))*(Engine.Isp_vac-Engine.Isp_sl)
	}
	return Engine.Isp_vac
}

func (v *VEHICLE) GetThrust(i int32) float64 {
	// H is distance to focus (center)
	if i == BOOSTER || !v.SysGuidance._MECO1 {
		//return v.Isp(H-Re) * 236 * g0 // 236 kg/s = M1D rate of fuel consumption
		return v.Isp(i) * EnginesMap[v.Stages[i].EngineID].Flow_rate * g0
	}
	return EnginesMap[v.Stages[STAGE2].EngineID].Th_vac //M1Dv.Th_vac
}

func (v *VEHICLE) GetThrust2(H float64, stage int32) float64 {
	// H is distance to focus (center)
	if stage == BOOSTER || !v.SysGuidance._MECO1 {
		//return v.Isp(H-Re) * 236 * g0 // 236 kg/s = M1D rate of fuel consumption
		return v.Isp2(H-Re) * EnginesMap[v.Stages[BOOSTER].EngineID].Flow_rate * g0
	}
	return EnginesMap[v.Stages[STAGE2].EngineID].Th_vac //M1Dv.Th_vac
}

func (r *VEHICLE) flip(i int32) {
	//	r.Stages[i].gamma = r.Stages[i].alpha + math.Pi		// retrograde
	r.Stages[i].gamma = math.Pi - r.Stages[i].gamma // retrograde
}
