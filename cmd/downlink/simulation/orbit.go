package main
import (
	"fmt"
	"math"
	"os"
	"time"
	"github.com/crmathieu/daq/packages/data"
)

var events *[]Event

//func main() {
func (v *VEHICLE) Simlaunch() {
	if events = v.InitGuidance(); events == nil {
		fmt.Println("Could not find flight profile data")
		os.Exit(-1)
	}

	// calculate speed boost on x axis based on earth rotation speed, latitude and Azimuth of trajectory
	vE = profile.EarthRotation * math.Cos(profile.LaunchLatitude) * math.Sin(profile.LaunchAzimuth)

	oldDRadius := float64(0) 
	newDRadius := float64(0) 
	apogee  := float64(0) 
	perigee  := float64(0)
	orbit := false

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

	//ticker := time.NewTicker(10 * time.Millisecond)
	//curstage := BOOSTER	
	for !orbit {

		select {
			//case <-ticker.C:
			default:
				// check for sequence in flight profile
				if v.CheckGuidanceEvents(events) != "MEI" {

					/*	SECO1			*/
	//				if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.SysGuidance._SECO1 {
					if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.hadEvent(data.E_SECO_1) { //}._SECO1 {
						// if no more fuel left, we trigger SECO. No more thrust
						//output_telemetry("SECO-1", f1, 1);
						fmt.Printf("\t\t\t\t\tSECO @ ---> %g degrees\n", (-3 * M_PI/2 + v.Stages[STAGE2].alpha - v.Stages[STAGE2].beta) * 180 / M_PI)
						v.SysGuidance._SECO1 = v.MSECO(STAGE2, data.E_SECO_1)
						v.Stages[STAGE2].apogee = v.Stages[STAGE2].DTF
						v.Stages[STAGE2].perigee = v.Stages[STAGE2].DTF
						//v.Stages[STAGE2].dt = 0.1
					}
			//		fmt.Println("SECO1 = ", _SECO1,", MECO1 = ",_MECO1)
					
					// 	Advance first stage	
	//				if !v.hadEvent(data.E_MECO_1) { //v.SysGuidance._MECO1 {
					if !v.hadEvent(data.E_STAGESEP) {
						v.timeStep(BOOSTER);
						//output_file(0, f);
					} else {
						// we had MECO
						if !orbit {
							v.timeStep(STAGE2)
							//output_file(1, f2)

							oldDRadius = newDRadius

							newDRadius = v.Stages[STAGE2].cx
				//			fmt.Println(newDRadius)
				//			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
							if oldDRadius < 0 && newDRadius > 0 {
								fmt.Println("Orbit!!!!")
								orbit = true
							}

							if v.hadEvent(data.E_SECO_1) {
								if v.Stages[STAGE2].DTF > apogee {
									v.Stages[STAGE2].apogee = v.Stages[STAGE2].DTF
								} 
								if v.Stages[STAGE2].DTF < perigee {
									v.Stages[STAGE2].perigee = v.Stages[STAGE2].DTF
								} 
							}
						}
					}
				}

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
	}
	fmt.Println("BOOSTER-clk:",	v.Stages[BOOSTER].Clock,", STAGE2 clk:", v.Stages[STAGE2].Clock)
	fmt.Printf("\nT%+07.2f\t%16.16s\t%.2f%s x %.2f%s\n", 
		//t,
		v.Stages[STAGE2].Clock, 
		"Orbit", 
		(v.Stages[STAGE2].perigee - Re) * 1e-3, 
		"km", 
		(v.Stages[STAGE2].apogee - Re) * 1e-3, 
		"km")
}

//func main2() {
func (v *VEHICLE) launch() {

	if events = v.InitGuidance(); events == nil {
		fmt.Println("Could not find flight profile data")
		os.Exit(-1)
	}

	// calculate speed boost on x axis based on earth rotation speed, latitude and Azimuth of trajectory
	vE = profile.EarthRotation * math.Cos(profile.LaunchLatitude) * math.Sin(profile.LaunchAzimuth)

	oldDRadius := float64(0) 
	newDRadius := float64(0) 
	apogee  := float64(0) 
	perigee  := float64(0)
	orbit := false

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

	ticker := time.NewTicker(10 * time.Millisecond)
	//curstage := BOOSTER
	for !orbit {

		select {
			case <-ticker.C:
	
			//default:
				for i := 0; i < len(*events); i++ {
//					if (*events)[i].Stage == BOOSTER && (math.Abs(v.Stages[BOOSTER].Clock - (*events)[i].T) < v.Stages[BOOSTER].dt/2) && !v.SysGuidance._stagesep {	
					if !v.hadEvent(data.E_STAGESEP) && (*events)[i].Stage == BOOSTER && (math.Abs(v.Stages[BOOSTER].Clock - (*events)[i].T) < v.Stages[BOOSTER].dt/2)  { //v.SysGuidance._stagesep {	
						// If an event in profile.txt occurs at this time, execute the event
						fmt.Println("booster event (delta = ",math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T),"), inc=",v.Stages[BOOSTER].dt/2)
						//v.execute((*events)[i].Id, nil) //f1)
						v.execute((*events)[i]) //f1)
					}
					if (*events)[i].Stage == STAGE2 && (math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T) < v.Stages[STAGE2].dt/2) { // stage 2 events
						fmt.Println("second stage event (delta = ",math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T),"), inc=",v.Stages[BOOSTER].dt/2)
//						v.execute((*events)[i].Id, nil) //f1);
						v.execute((*events)[i]) //f1);
					}
				}
				
//				v.CheckGuidanceEvents(events)

				/*	SECO1			*/
//				if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.SysGuidance._SECO1 {
				if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.hadEvent(data.E_SECO_1) { //}._SECO1 {
					//output_telemetry("SECO-1", f1, 1);
					fmt.Printf("\t\t\t\t\tSECO @ ---> %g degrees\n", (-3 * M_PI/2 + v.Stages[STAGE2].alpha - v.Stages[STAGE2].beta) * 180 / M_PI)
					v.SysGuidance._SECO1 = v.MSECO(1, data.E_SECO_1)
					apogee = v.Stages[STAGE2].DTF
					perigee = v.Stages[STAGE2].DTF
					//v.Stages[STAGE2].dt = 0.1
				}
		//		fmt.Println("SECO1 = ", _SECO1,", MECO1 = ",_MECO1)
				
				// 	Advance first stage	
//				if !v.hadEvent(data.E_MECO_1) { //v.SysGuidance._MECO1 {
				if !v.hadEvent(data.E_STAGESEP) {
					v.timeStep(BOOSTER);
					//output_file(0, f);
				} else {
					// we had MECO
					if !orbit {
						v.timeStep(STAGE2)
						//output_file(1, f2)

						oldDRadius = newDRadius

						newDRadius = v.Stages[STAGE2].cx
			//			fmt.Println(newDRadius)
			//			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
						if oldDRadius < 0 && newDRadius > 0 {
							fmt.Println("Orbit!!!!")
							orbit = true
						}

						if v.hadEvent(data.E_SECO_1) {
							if v.Stages[STAGE2].DTF > apogee {
								apogee = v.Stages[STAGE2].DTF
							} 
							if v.Stages[STAGE2].DTF < perigee {
								perigee = v.Stages[STAGE2].DTF
							} 
						}
					}
				}

				// 	Advance second stage	
//				if v.SysGuidance._MECO1 && !orbit {
/*				if v.hadEvent(data.E_MECO_1) && !orbit {
					//curstage = STAGE2
					v.timeStep(STAGE2)
					//output_file(1, f2)

					oldDRadius = newDRadius

					newDRadius = v.Stages[STAGE2].cx
		//			fmt.Println(newDRadius)
		//			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
					if oldDRadius < 0 && newDRadius > 0 {
						fmt.Println("Orbit!!!!")
						orbit = true
					}

//					if v.SysGuidance._SECO1 {
					if v.hadEvent(data.E_SECO_1) {
						if v.Stages[STAGE2].DTF > apogee {
							apogee = v.Stages[STAGE2].DTF
						} 
						if v.Stages[STAGE2].DTF < perigee {
							perigee = v.Stages[STAGE2].DTF
						} 
					}
				}
*/
//				if !v.SysGuidance._stagesep {
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
	}
	fmt.Println("BOOSTER-clk:",	v.Stages[BOOSTER].Clock,", STAGE2 clk:", v.Stages[STAGE2].Clock)
	fmt.Printf("\nT%+07.2f\t%16.16s\t%.2f%s x %.2f%s\n", 
		//t,
		v.Stages[STAGE2].Clock, 
		"Orbit", 
		(perigee - Re) * 1e-3, 
		"km", 
		(apogee - Re) * 1e-3, 
		"km")
	


/*	for !orbit {

		select {
			case <-ticker.C:
			//default:
				// check for sequence in flight profile
				if v.CheckGuidanceEvents(events) != "MEI" {

					//	SECO1			
	//				if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.SysGuidance._SECO1 {
					if (v.NoFuel(STAGE2) || v.Stages[STAGE2].VAbsolute > math.Sqrt(G * Me / v.Stages[STAGE2].DTF)) && !v.hadEvent(data.E_SECO_1) { //}._SECO1 {
						// if no more fuel left, we trigger SECO. No more thrust
						//output_telemetry("SECO-1", f1, 1);
						fmt.Printf("\t\t\t\t\tSECO @ ---> %g degrees\n", (-3 * M_PI/2 + v.Stages[STAGE2].alpha - v.Stages[STAGE2].beta) * 180 / M_PI)
						v.SysGuidance._SECO1 = v.MSECO(STAGE2, data.E_SECO_1)
						v.Stages[STAGE2].apogee = v.Stages[STAGE2].DTF
						v.Stages[STAGE2].perigee = v.Stages[STAGE2].DTF
						//v.Stages[STAGE2].dt = 0.1
					}
			//		fmt.Println("SECO1 = ", _SECO1,", MECO1 = ",_MECO1)
					
					// 	Advance first stage	
	//				if !v.hadEvent(data.E_MECO_1) { //v.SysGuidance._MECO1 {
					if !v.hadEvent(data.E_STAGESEP) {
						v.timeStep(BOOSTER);
						//output_file(0, f);
					} else {
						// we had MECO
						if !orbit {
							v.timeStep(STAGE2)
							//output_file(1, f2)

							oldDRadius = newDRadius

							newDRadius = v.Stages[STAGE2].cx
				//			fmt.Println(newDRadius)
				//			fmt.Println("OLDradius = ", oldDRadius,", newDradius = ",newDRadius)
							if oldDRadius < 0 && newDRadius > 0 {
								fmt.Println("Orbit!!!!")
								orbit = true
							}

							if v.hadEvent(data.E_SECO_1) {
								if v.Stages[STAGE2].DTF > apogee {
									v.Stages[STAGE2].apogee = v.Stages[STAGE2].DTF
								} 
								if v.Stages[STAGE2].DTF < perigee {
									v.Stages[STAGE2].perigee = v.Stages[STAGE2].DTF
								} 
							}
						}
					}
				}

				if !v.hadEvent(data.E_STAGESEP) { //}._stagesep {
					v.Stages[BOOSTER].Clock = v.Stages[BOOSTER].Clock + v.Stages[BOOSTER].dt
				}
				v.Stages[STAGE2].Clock = v.Stages[STAGE2].Clock + v.Stages[STAGE2].dt
				//t = t + dt
				//fmt.Println(t)

		}
	}
	fmt.Println("BOOSTER-clk:",	v.Stages[BOOSTER].Clock,", STAGE2 clk:", v.Stages[STAGE2].Clock)
	fmt.Printf("\nT%+07.2f\t%16.16s\t%.2f%s x %.2f%s\n", 
		//t,
		v.Stages[STAGE2].Clock, 
		"Orbit", 
		(v.Stages[STAGE2].perigee - Re) * 1e-3, 
		"km", 
		(v.Stages[STAGE2].apogee - Re) * 1e-3, 
		"km")
*/
}

func (v *VEHICLE) CheckGuidanceEvents(events *[]Event) string {

	for i := 0; i < len(*events); i++ {
		if !v.hadEvent(data.E_STAGESEP) && (*events)[i].Stage == BOOSTER && (math.Abs(v.Stages[BOOSTER].Clock - (*events)[i].T) < v.Stages[BOOSTER].dt/2)  { //v.SysGuidance._stagesep {	
			// If an event in profile.txt occurs at this time, execute the event
			fmt.Println("booster event (delta = ",math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T),"), inc=",v.Stages[BOOSTER].dt/2)
			//v.execute((*events)[i].Id, nil) //f1)
			v.execute((*events)[i]) //f1)
			return (*events)[i].Id
		}
		if (*events)[i].Stage == STAGE2 && (math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T) < v.Stages[STAGE2].dt/2) { // stage 2 events
			fmt.Println("second stage event (delta = ",math.Abs(v.Stages[STAGE2].Clock - (*events)[i].T),"), inc=",v.Stages[BOOSTER].dt/2)
			v.execute((*events)[i]) //f1);
			return (*events)[i].Id
		}
	}
	return ""
}