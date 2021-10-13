package main

import (
	"fmt"
	"os"
)

//inline void output_telemetry(char* event, FILE *f2, int i)	// i = stage
func (v *VEHICLE) output_telemetry(event string, f2 *os.File, i int32) { // i = stage

	temp_s0 := v.Stages[i].px - (vE * v.Stages[i].Clock)
	temp_S := v.Stages[i].DTF - Re
	temp_V := v.Stages[i].rvelocity
	dist := "m"
	vel := "m/s"

	if temp_s0 > 1e3 || temp_S > 1e3 {
		temp_s0 = temp_s0 * 1e-3
		temp_S = temp_S * 1e-3
		dist = "km"
	}

	if f2 != nil {
		fmt.Fprintf(f2, "%g\t%f\t%f\t%f\t%f\t%f\t%f\t%s\n",
			v.Stages[i].Clock,
			v.Stages[i].px*1e-3,
			(v.Stages[i].py-Re)*1e-3,
			(v.Stages[i].DTF-Re)*1e-3,
			temp_V,
			v.Stages[i].Acc/g0,
			v.Stages[i].mass,
			event)
	}

	if temp_V > 1e3 {
		temp_V = temp_V * 1e-3
		vel = "km/s"
	}

	fmt.Printf("T%+07.2f\t%16.16s\t%.2f%s x %.2f%s @ %.2f%s\n", v.Stages[i].Clock, event, temp_s0, dist, temp_S, dist, temp_V, vel)
	//	fmt.Printf("T%+07.2f\t%16.16s\t%.2f%s x %.2f%s @ %.2f%s\n", t, event, temp_s0, dist, temp_S, dist, temp_V, vel);
}

//inline void output_file(int i, FILE *f)
func (v *VEHICLE) output_file(i int32, f *os.File) {
	fmt.Fprintf(f, "%g\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n",
		v.Stages[i].Clock,
		v.Stages[i].px*1e-3,
		(v.Stages[i].py-Re)*1e-3,
		(v.Stages[i].DTF-Re)*1e-3,
		v.Stages[i].rvelocity,
		v.Stages[i].Acc/g0,
		v.Stages[i].mass,
		v.Stages[i].ThrottleRate)
}
