package main
import(
	"math"
	"fmt"
	//"os"
	"time"
	"github.com/crmathieu/daq/packages/data"
)

func (v *VEHICLE) boosterGuidance() {
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
	ticker := time.NewTicker(10 * time.Millisecond)
	//curstage := BOOSTER	
	for !touchdown {
		select {
			case <-ticker.C:
			//default:
			/*	Execute events		*/
			//for(i=0;i<N;i++) {
			for i := 0;i < len(*event); i++ {	
				if (*event)[i].Stage == BOOSTER && math.Abs(v.Stages[BOOSTER].Clock - (*event)[i].T) < v.Stages[BOOSTER].dt/2	{ 
					// If a booster event in profile.txt occurs at this time, execute the event
					fmt.Println("returning booster event")
					v.execute((*event)[i].Id, nil) //f1)
				}
				/*if (*event)[i].Stage == 1 && math.Abs(v.Stages[BOOSTER].Clock - (*event)[i].T) < dt/2	{
					// Stage2 events
					execute((*event)[i].Id, nil) //f1)
				}*/
			}

			/*	End Landing Burn	*/
			if (v.Stages[BOOSTER].Mf < 5 || (v.SysGuidance._LBURN && v.Stages[BOOSTER].DTF - Re < 0.01)) && !v.SysGuidance._MECO3 {
				//output_telemetry("MECO-3", nil, 0) //f1, 0)
				v.SysGuidance._MECO3 = v.MSECO(0, data.E_LBURNCO) //data.E_MECO_3) //, &_MECO3);
				v.SysGuidance._LBURN = false
				fmt.Println("\t",v.Stages[0].Clock,"--> Landing burn stopped!!! distance to ground:", v.Stages[BOOSTER].DTF - Re, ", remaoning fuel:",v.Stages[BOOSTER].Mf)

			} else {
				// touchdown
				if v.SysGuidance._release && v.Stages[BOOSTER].DTF < Re && !touchdown {			// If Alt = 0.0m
					//output_telemetry("Touchdown", nil, 0)
					touchdown = true
					//v.Stages[BOOSTER].dt = 0.1
				} /*else {
					// SECO1
					if (v.Stages[STAGE2].Mf < 5 || (v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me/v.Stages[STAGE2].DTF))) && !v.SysGuidance._SECO1 {
						//output_telemetry("SECO-1", nil, 1) //f1, 1)
						v.SysGuidance._SECO1 = v.MSECO(1) //, &_SECO1);
					}
				}*/
			}

			//	Advance First stage	
			if !touchdown {
				v.timeStep(0)
				//output_file(0, nil) //f)
			}

			//	Advance Second stage
			/*if v.SysGuidance._MECO1 {
				v.timeStep(1)
				//output_file(1, nil) //f2)
			}*/

			//t = t + dt
			v.Stages[BOOSTER].Clock = v.Stages[BOOSTER].Clock + v.Stages[BOOSTER].dt
		}
	}
	fmt.Println("TOUCHDOWN!!!!!!!!!!!!!")
	//os.Exit(1)
}
