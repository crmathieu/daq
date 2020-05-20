
package main
import (
	"fmt"
	"os"
)

//inline void output_telemetry(char* event, FILE *f2, int i)	// i = stage
func output_telemetry(event string, f2 *os.File, i int32) {	// i = stage


	temp_s0 := F9.Stages[i].cx - (vE * F9.Stages[i].Clock)
	temp_S := F9.Stages[i].PolarDistance - Re 
	temp_V := F9.Stages[i].VRelative
	dist := "m"
	vel := "m/s"
	
	if(temp_s0 > 1e3 || temp_S > 1e3) {
		temp_s0 = temp_s0 * 1e-3
		temp_S = temp_S * 1e-3
		dist = "km";
	}

	if f2 != nil {
		fmt.Fprintf(f2, "%g\t%f\t%f\t%f\t%f\t%f\t%f\t%s\n", 
			F9.Stages[i].Clock, 
			F9.Stages[i].cx * 1e-3,
			(F9.Stages[i].cy - Re) * 1e-3,
			(F9.Stages[i].PolarDistance - Re) * 1e-3, 
			temp_V, 
			F9.Stages[i].Acc / g0, 
			F9.Stages[i].Mass, 
			event)
	}

	if(temp_V > 1e3) {
		temp_V = temp_V * 1e-3;
		vel = "km/s";
	}

	fmt.Printf("T%+07.2f\t%16.16s\t%.2f%s x %.2f%s @ %.2f%s\n", F9.Stages[i].Clock, event, temp_s0, dist, temp_S, dist, temp_V, vel);
//	fmt.Printf("T%+07.2f\t%16.16s\t%.2f%s x %.2f%s @ %.2f%s\n", t, event, temp_s0, dist, temp_S, dist, temp_V, vel);
}

//inline void output_file(int i, FILE *f)
func output_file(i int32, f *os.File) {
	fmt.Fprintf(f, "%g\t%f\t%f\t%f\t%f\t%f\t%f\t%f\n", 
		F9.Stages[i].Clock, 
		F9.Stages[i].cx * 1e-3,
		(F9.Stages[i].cy - Re) * 1e-3,
		(F9.Stages[i].PolarDistance - Re) * 1e-3, 
		F9.Stages[i].VRelative, 
		F9.Stages[i].Acc / g0, 
		F9.Stages[i].Mass, 
		F9.Stages[i].ThrottleRate)
}

