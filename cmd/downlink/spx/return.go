package main
import(
	"math"
	"fmt"
)

func boosterGuidance() {
	touchdown := false

	/*************************************************************************************************/
	/*	Launch/Pitch Kick/Gravity turn				*/
	/*	First Stage: Flip/Entry burn/Landing burn		*/
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

	//do{
	for !touchdown {
		/*	Execute events		*/
		//for(i=0;i<N;i++) {
		for i := 0;i < len(*event); i++ {	
			if (*event)[i].Stage == BOOSTER && math.Abs(F9.Stages[BOOSTER].Clock - (*event)[i].T) < F9.Stages[BOOSTER].dt/2	{ 
				// If a booster event in profile.txt occurs at this time, execute the event
				execute((*event)[i].Id, nil) //f1)
			}
			/*if (*event)[i].Stage == 1 && math.Abs(F9.Stages[BOOSTER].Clock - (*event)[i].T) < dt/2	{
				// Stage2 events
				execute((*event)[i].Id, nil) //f1)
			}*/
		}

		/*	End Landing Burn	*/
		if (F9.Stages[BOOSTER].Mf < 5 || (_LBURN && F9.Stages[BOOSTER].PolarDistance - Re < 0.01)) && !_MECO3 {
			output_telemetry("MECO-3", nil, 0) //f1, 0)
			_MECO3 = MSECO(0) //, &_MECO3);
			_LBURN = false
		} else {
			// touchdown
			if _release && F9.Stages[BOOSTER].PolarDistance < Re && !touchdown {			// If Alt = 0.0m
				output_telemetry("Touchdown", nil, 0)
				touchdown = true
				F9.Stages[BOOSTER].dt = 0.1
			} else {
				// SECO1
				if (F9.Stages[STAGE2].Mf < 5 || (F9.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me/F9.Stages[STAGE2].PolarDistance))) && !_SECO1 {
					output_telemetry("SECO-1", nil, 1) //f1, 1)
					_SECO1 = MSECO(1) //, &_SECO1);
				}
			}
		}

		//	Advance First stage	
		if !touchdown {
			leapfrog_step(0)
			output_file(0, nil) //f)
		}

		//	Advance Second stage
		if _MECO1 {
			leapfrog_step(1)
			output_file(1, nil) //f2)
		}

		//t = t + dt
		F9.Stages[BOOSTER].Clock = F9.Stages[BOOSTER].Clock + F9.Stages[BOOSTER].dt

	}
	fmt.Println()
}
