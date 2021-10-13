package main

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/crmathieu/daq/packages/data"
)

var events *[]Pevent

var oldDRadius = float64(0)
var newDRadius = float64(0)
var apogee = float64(0)
var perigee = float64(0)
var orbit = false
var crashed = false
var TargetOrbitalVelocity = 0.0

// launch ---------------------------------------------------------------------
// launch to orbit in realtime or calculated modes
// ----------------------------------------------------------------------------
func (v *VEHICLE) launch(realTime bool) {
	// calculate speed boost on x axis based on earth rotation speed, latitude and Azimuth of trajectory
	// profile.LaunchAzimuth = math.Asin(math.Cos(profile.OrbitInclination) / math.Cos(math.Cos(profile.LaunchLatitude)))
	// vE = profile.EarthRotation * math.Cos(profile.LaunchLatitude) * math.Sin(profile.LaunchAzimuth)
	vE = profile.EarthRotation * math.Sin(deg2rad(profile.LaunchAzimuth))

	println("Earth ROT=", vE, "m/s")

	fmt.Println(vE)

	TargetOrbitalVelocity = math.Sqrt(G * Me / (profile.OrbitInsertion + Re))

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

	var ticker *time.Ticker
	//var elasped = 0.0
	if realTime {
		ticker = time.NewTicker(10 * time.Millisecond)
		for !orbit && !crashed {
			select {
			case <-ticker.C:
				/*elasped = elasped + 10e-3
				if elasped > 45.0 {
					//fmt.Println("> 15")
					v.Stages[BOOSTER].RunningEngines = 0
				}*/
				v.addStep()
			}
		}
	} else {
		for !orbit && !crashed {
			v.addStep()
		}
	}
	v.showOrbitDetails()
	os.Exit(0)
}

// showOrbitDetails -----------------------------------------------------------
// Displays final results about orbit
// ----------------------------------------------------------------------------
func (v *VEHICLE) showOrbitDetails() {
	fmt.Println("BOOSTER-clk:", v.Stages[BOOSTER].Clock, ", STAGE2 clk:", v.Stages[STAGE2].Clock)
	fmt.Printf("\nT%+07.2f\t%16.16s\t%.2f%s x %.2f%s\n",
		//t,
		v.Stages[STAGE2].Clock,
		"Orbit",
		(perigee-Re)*1e-3,
		"km",
		(apogee-Re)*1e-3,
		"km")
	fmt.Println("MaxQ =", mQ.MaxQ, "at time", mQ.Time, ", altitude:", mQ.Alt, ",\nrange:", mQ.Range, "at speed = ", mQ.Velocity, "m/s", "\nangle=", mQ.Angle, ", density:", mQ.RhoMQ)

	// eccentricity is given with A = a(1+e) and P = a(1-e), hence A/P = (1+e)/(1-e) and
	// finally e = (A-P)/(A+P)
	fmt.Printf("Orbit eccentricity is %v\n", float64(apogee-perigee)/float64(apogee+perigee))

}

//var TargetOrbitalVelocity = math.Sqrt(G * Me / (profile.OrbitInsertion + Re))

// addStep --------------------------------------------------------------------
// Takes care of moving the rocket one step
// ----------------------------------------------------------------------------
func (v *VEHICLE) addStep() {
	// take care of flight profile events
	v.CheckGuidanceEvents(events)

	// check 2nd stage status
	if v.SysGuidance._SEI1 && !v.hadEvent(data.E_SECO_1) {
		// engine still supposed to be running
		//		if v.NoFuel(STAGE2) || (v.Stages[STAGE2].VAbsolute >= math.Sqrt(G*Me/(profile.OrbitInsertion+Re))) { //||
		//fmt.Println("VEL=", v.Stages[STAGE2].AVel*3.6, ", TGT=", TargetOrbitalVelocity*3.6)
		if v.NoFuel(STAGE2) || (v.Stages[STAGE2].avelocity >= TargetOrbitalVelocity) { //||
			fmt.Print("\n!!! We have premature SECO due to ")
			if v.NoFuel(STAGE2) {
				fmt.Println("No more fuel...")
			} else {
				fmt.Println("Target orbit velocity reached...")
			}
			fmt.Println(
				"Time ..................", v.Stages[STAGE2].Clock, "s",
				"\nRemaining fuel ........", v.Stages[STAGE2].Mf, "kg",
				"\nAltitude ..............", (v.Stages[STAGE2].altitude)*1e-3, "km",
				"\nVelocity ..............", v.Stages[STAGE2].avelocity*3.6, "km/h",
				"\nTarget Velocity .......", TargetOrbitalVelocity*3.6, "km/h",
				"\nFlight Path ...........", rad2deg(v.Stages[STAGE2].gamma), "deg",
				"\nAngular range .........", rad2deg(v.Stages[STAGE2].beta), "deg")

			if v.Stages[STAGE2].altitude > profile.OrbitInsertion {
				v.Stages[STAGE2].gamma = 0
			}
			// no fuel left or our speed exceeds the target orbital speed
			if v.NoFuel(STAGE2) {
				fmt.Println("\n\n---------> 2nd STAGE EMPTY!!!! @", v.Stages[STAGE2].Clock, "\n") //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			}
			/*fmt.Printf("\n************\nSECO @ ---> %g seconds\n", v.Stages[STAGE2].Clock)                 //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Remaining fuel ...... ", v.Stages[STAGE2].Mf, "kg")                               //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Velocity ............ ", v.Stages[STAGE2].RVel*3.6, "km/h")                       //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Altitude ............ ", (v.Stages[STAGE2].altitude)*1e-3, "km")                  //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			fmt.Println("Injection Angle ..... ", rad2deg(v.Stages[STAGE2].gamma), "degres\n************") //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
			*/
			v.SysGuidance._SECO1 = v.MSECO(STAGE2, data.E_SECO_1)
			apogee = v.Stages[STAGE2].DTF
			perigee = v.Stages[STAGE2].DTF
			//fmt.Println("AP=", apogee, "PER=", perigee)
			//v.Stages[STAGE2].dt = 0.1

		} else {
			//fmt.Println(v.Stages[STAGE2].VAbsolute * 3.6)
		}
	}

	// check for stage separation
	if !v.hadEvent(data.E_STAGESEP) {
		// 	Advance first stage
		v.timeStep(BOOSTER)
		//output_file(0, f);
	} else {
		// Advance 2nd stage
		if !orbit {
			v.timeStep(STAGE2)
			//output_file(1, f2)
			/*
				oldDRadius = newDRadius

				newDRadius = v.Stages[STAGE2].altitude + Re //cx
				//			fmt.Println(newDRadius)
				//			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
				if oldDRadius < 0 && newDRadius > 0 {
					fmt.Println("Orbit!!!!")
					orbit = true
				}
			*/
			// update apagee and perigee based on current distance to focus (center)
			if v.hadEvent(data.E_SECO_1) {
				if v.Stages[STAGE2].altitude <= 0 {
					fmt.Println("Crashed", v.Stages[STAGE2].drange, "km down range")
					crashed = true
					return
				} else {
					//fmt.Println("Altitude", v.Stages[STAGE2].altitude, "km, Force", v.Stages[STAGE2].Force)

				}

				if v.Stages[STAGE2].DTF > apogee {
					apogee = v.Stages[STAGE2].DTF
				}
				if v.Stages[STAGE2].DTF < perigee {
					perigee = v.Stages[STAGE2].DTF
				}
			}
		}
	}

	// update clock
	if !v.hadEvent(data.E_STAGESEP) { //}._stagesep {
		v.Stages[BOOSTER].Clock = v.Stages[BOOSTER].Clock + v.Stages[BOOSTER].dt
	}
	v.Stages[STAGE2].Clock = v.Stages[STAGE2].Clock + v.Stages[STAGE2].dt
	//t = t + dt
	//fmt.Println(t)
	/*				fmt.Printf("Time: %03.2f --> Altitude: %04.3f ---> Range: %04.3f --> VA: %04.0f --> VR: %04.0f\n",
						v.Stages[curstage].Clock,
					(v.Stages[curstage].cy - Re)*1e-3,
					(v.Stages[curstage].cx)*1e-3,
						v.Stages[curstage].VAbsolute,
						v.Stages[curstage].VRelative)*/
}

// CheckGuidanceEvents --------------------------------------------------------
// Detects flight profile events
// ----------------------------------------------------------------------------
func (v *VEHICLE) CheckGuidanceEvents(events *[]Pevent) string {

	for i := 0; i < len(*events); i++ {
		if !v.hadEvent(data.E_STAGESEP) && (*events)[i].Stage == BOOSTER && (math.Abs(v.Stages[BOOSTER].Clock-(*events)[i].T) < v.Stages[BOOSTER].dt/2) { //v.SysGuidance._stagesep {
			// takes care of pre stage-separation BOOSTER events
			//fmt.Println("booster event (delta = ", math.Abs(v.Stages[STAGE2].Clock-(*events)[i].T), "), inc=", v.Stages[BOOSTER].dt/2)
			//v.execute((*events)[i].Id, nil) //f1)
			v.execute((*events)[i]) //f1)
			return (*events)[i].Id
		}

		if (*events)[i].Stage == STAGE2 && (math.Abs(v.Stages[STAGE2].Clock-(*events)[i].T) < v.Stages[STAGE2].dt/2) { // stage 2 events
			// takes care of STAGE2 events
			//fmt.Println("second stage event (delta = ", math.Abs(v.Stages[STAGE2].Clock-(*events)[i].T), "), inc=", v.Stages[BOOSTER].dt/2)
			v.execute((*events)[i]) //f1);
			return (*events)[i].Id
		}
	}
	return ""
}
