package main

import (
	"fmt"
	"math"
	//"os"
	"time"

	"github.com/crmathieu/daq/packages/data"
)

var touchdown = false

// boosterGuidance ------------------------------------------------------------
// Takes care of landing the booster after stage separation
// ----------------------------------------------------------------------------
func (v *VEHICLE) boosterGuidance(realTime bool) {

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

	var ticker *time.Ticker

	if realTime {
		ticker = time.NewTicker(10 * time.Millisecond)
		for !touchdown {
			select {
			case <-ticker.C:
				v.addLandingStep()
			}
		}
	} else {
		for !touchdown {
			v.addLandingStep()
		}
	}
	//v.showOrbitDetails()

	//ticker := time.NewTicker(10 * time.Millisecond)
	//curstage := BOOSTER
	/*
		for !touchdown {
			select {
			case <-ticker.C:

				for i := 0; i < len(*events); i++ {
					// execute BOOSTER events
					if (*events)[i].Stage == BOOSTER && math.Abs(v.Stages[BOOSTER].Clock-(*events)[i].T) < v.Stages[BOOSTER].dt/2 {
						fmt.Println("returning booster event")
						//					v.executes((*event)[i].Id, nil) //f1)
						v.execute((*events)[i]) //f1)
					}
				}

				//	End Landing Burn
				if (v.NoFuel(BOOSTER) || (v.SysGuidance._LBURN && v.Stages[BOOSTER].DTF-Re < 0.01)) && !v.SysGuidance._MECO3 {
					//output_telemetry("MECO-3", nil, 0) //f1, 0)
					v.SysGuidance._MECO3 = v.MSECO(0, data.E_LBURNO) //data.E_MECO_3) //, &_MECO3);
					v.SysGuidance._LBURN = false
					fmt.Println("\t", v.Stages[0].Clock, "--> Landing burn stopped!!! distance to ground:", v.Stages[BOOSTER].DTF-Re, ", remaoning fuel:", v.Stages[BOOSTER].Mf)

				} else {
					// touchdown
					if v.SysGuidance._release && v.Stages[BOOSTER].DTF < Re && !touchdown { // If Alt = 0.0m
						//output_telemetry("Touchdown", nil, 0)
						touchdown = true
					}
				}

				//	Advance First stage
				if !touchdown {
					v.timeStep(BOOSTER)
					//output_file(0, nil) //f)
				}

				//t = t + dt
				v.Stages[BOOSTER].Clock = v.Stages[BOOSTER].Clock + v.Stages[BOOSTER].dt
			}
		}*/
	fmt.Println("TOUCHDOWN!!!!!!!!!!!!!", v.Stages[BOOSTER].Clock)
	//os.Exit(1)
}

func (v *VEHICLE) addLandingStep() {
	for i := 0; i < len(*events); i++ {
		// execute BOOSTER events
		if (*events)[i].Stage == BOOSTER && math.Abs(v.Stages[BOOSTER].Clock-(*events)[i].T) < v.Stages[BOOSTER].dt/2 {
			fmt.Println("returning booster event")
			//					v.executes((*event)[i].Id, nil) //f1)
			v.execute((*events)[i]) //f1)
		}
	}

	/*	End Landing Burn	*/
	if (v.NoFuel(BOOSTER) || (v.SysGuidance._LBURN && v.Stages[BOOSTER].DTF-Re < 0.01)) && !v.SysGuidance._MECO3 {
		//output_telemetry("MECO-3", nil, 0) //f1, 0)
		v.SysGuidance._MECO3 = v.MSECO(0, data.E_LBURNO) //data.E_MECO_3) //, &_MECO3);
		v.SysGuidance._LBURN = false
		fmt.Println("\t", v.Stages[BOOSTER].Clock, "--> Landing burn stopped!!! distance to ground:", v.Stages[BOOSTER].DTF-Re, ", remaoning fuel:", v.Stages[BOOSTER].Mf)

	} else {
		// touchdown
		if v.SysGuidance._release && v.Stages[BOOSTER].DTF < Re && !touchdown { // If Alt = 0.0m
			//output_telemetry("Touchdown", nil, 0)
			touchdown = true
		}
	}

	//	Advance First stage
	if !touchdown {
		v.landingTimeStep()
		//output_file(0, nil) //f)
	}

	//t = t + dt
	v.Stages[BOOSTER].Clock = v.Stages[BOOSTER].Clock + v.Stages[BOOSTER].dt
}
