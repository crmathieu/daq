package main
import (
	"fmt"
	"math"
	//"os"
)

var event *[]Event

func main() {

	oldDRadius := float64(0) 
	newDRadius := float64(0) 
	apogee  := float64(0) 
	perigee  := float64(0)
	orbit := false

	event = Init()
	vE = profile.EarthRotation
	fmt.Println(vE)

/*************************************************************************************************/
/*	Launch/Pitch Kick/ Gravity Turn				*/
/*	First Stage: Takeoff					*/
/*************************************************************************************************/

/*	f, err := os.Create("./Stage1.dat")
	if err != nil {
		panic("Cannot create file")
	}
	defer f.Close()
	f1, err := os.Create("./Points.dat")
	if err != nil {
		panic("Cannot create file")
	}
	defer f1.Close()
	f2, err := os.Create("./Stage3.dat")
	if err != nil {
		panic("Cannot create file")
	}
	defer f2.Close()
*/

	for !orbit {
		/*	Execute events		*/
		for i := 0;i < len(*event); i++ {
			if (*event)[i].Stage == BOOSTER && (math.Abs(F9.Stages[BOOSTER].Clock - (*event)[i].T) < dt/2) && !_MECO1 {	
				// If an event in profile.txt occurs at this time, execute the event
				execute((*event)[i].Id, nil) //f1)
			}
			if (*event)[i].Stage == STAGE2 && (math.Abs(F9.Stages[STAGE2].Clock - (*event)[i].T) < dt/2) { // stage 2 events
				execute((*event)[i].Id, nil) //f1);
			}
		}

		/*	SECO1			*/
		if (F9.Stages[STAGE2].Mf < 5 || F9.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / F9.Stages[STAGE2].PolarDistance)) && !_SECO1 {
			//output_telemetry("SECO-1", f1, 1);
			fmt.Printf("\t\t\t\t\t@ ---> %g degrees\n", (-3 * M_PI/2 + F9.Stages[STAGE2].alpha - F9.Stages[STAGE2].beta) * 180 / M_PI)
			_SECO1 = MSECO(1)
			apogee = F9.Stages[STAGE2].PolarDistance
			perigee = F9.Stages[STAGE2].PolarDistance
			dt = 0.1
		}
//		fmt.Println("SECO1 = ", _SECO1,", MECO1 = ",_MECO1)
		/* 	Advance first stage	*/
		if !_MECO1 {
			leapfrog_step(0);
			//output_file(0, f);
		}

		/* 	Advance second stage	*/
		if _MECO1 && !orbit {
			leapfrog_step(1)
			//output_file(1, f2)

			oldDRadius = newDRadius

			newDRadius = F9.Stages[STAGE2].cx
			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
			if oldDRadius < 0 && newDRadius > 0 {
				fmt.Println("Orbit!!!!")
				orbit = true
			}

			if _SECO1 {
				if F9.Stages[STAGE2].PolarDistance > apogee {
					apogee = F9.Stages[STAGE2].PolarDistance
				} 
				if F9.Stages[STAGE2].PolarDistance < perigee {
					perigee = F9.Stages[STAGE2].PolarDistance
				} 
			}
		}

		F9.Stages[BOOSTER].Clock = F9.Stages[BOOSTER].Clock + dt
		F9.Stages[STAGE2].Clock = F9.Stages[STAGE2].Clock + dt
		//t = t + dt
		//fmt.Println(t)
	}

	fmt.Printf("\nT%+07.2f\t%16.16s\t%.2f%s x %.2f%s\n", 
		//t,
		F9.Stages[STAGE2].Clock, 
		"Orbit", 
		(perigee - Re) * 1e-3, 
		"km", 
		(apogee - Re) * 1e-3, 
		"km")
}