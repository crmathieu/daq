package main

import (
	//	"math"
	//	"os"
	"fmt"

	"github.com/crmathieu/daq/packages/data"
)

func GetEventData(P *Profile, Id string) (*Pevent, int) {
	for i, event := range P.Events {
		if event.Id == Id {
			return &event, i
		}
	}
	return nil, -1
}

func (v *VEHICLE) hadEvent(eventid uint32) bool {
	if v.EventsMap&eventid != 0 {
		return true
	}
	return false
}

//inline void pitchStart(int *x)
//func pitchStart(x *int32) {
//	gamma[0] = (M_PI / 2) - 0.025;
//	*x = 1;
//}
func (r *VEHICLE) pitchStart(event *Pevent) bool {
	r.Stages[BOOSTER].gamma = (M_PI / 2) - deg2rad(event.Gamma0) // 0.025
	asc.aPhases[0].startingTime = r.Stages[BOOSTER].Clock
	r.EventsMap = r.EventsMap | data.E_STARTPITCH
	r.LastEvent = data.E_STARTPITCH
	return true
}

func (r *VEHICLE) Throttle(event *Pevent) bool {
	//	r.Stages[event.Stage].ThrottleRate = float64(event.Rate) / 100
	r.Stages[event.Stage].ThrottleRate = float64(event.Rate)
	//	r.EventsMap = r.EventsMap | data.EventMapping[event.Id]
	r.EventsMap = r.EventsMap | data.EventInfoMapping[event.Id].Id
	r.LastEvent = data.EventInfoMapping[event.Id].Id
	return true
}

//func (r *VEHICLE) Ignition(stage int32, eventid uint32, throttleRate float64, num_engs int32) bool {
func (r *VEHICLE) Ignition(event *Pevent, num_engs int32) bool {

	r.Stages[event.Stage].ThrottleRate = event.Rate //1.0
	r.Stages[event.Stage].RunningEngines = num_engs
	r.EventsMap = r.EventsMap | data.EventInfoMapping[event.Id].Id //eventid
	r.LastEvent = data.EventInfoMapping[event.Id].Id               //eventid
	return true
}

func (r *VEHICLE) sync_stages() {

	r.Stages[STAGE2].altitude = r.Stages[BOOSTER].altitude
	r.Stages[STAGE2].drange = r.Stages[BOOSTER].drange

	r.Stages[STAGE2].RVel = r.Stages[BOOSTER].RVel
	r.Stages[STAGE2].AVel = r.Stages[BOOSTER].AVel
	r.Stages[STAGE2].VRelative = r.Stages[BOOSTER].VRelative
	r.Stages[STAGE2].VAbsolute = r.Stages[BOOSTER].VAbsolute

	r.Stages[STAGE2].px = r.Stages[BOOSTER].px
	r.Stages[STAGE2].py = r.Stages[BOOSTER].py
	r.Stages[STAGE2].pz = r.Stages[BOOSTER].pz

	r.Stages[STAGE2].vx = r.Stages[BOOSTER].vx
	r.Stages[STAGE2].vy = r.Stages[BOOSTER].vy
	r.Stages[STAGE2].vz = r.Stages[BOOSTER].vz

	r.Stages[STAGE2].ax = r.Stages[BOOSTER].ax
	r.Stages[STAGE2].ay = r.Stages[BOOSTER].ay
	r.Stages[STAGE2].az = r.Stages[BOOSTER].az

	r.Stages[STAGE2].vAx = r.Stages[BOOSTER].vAx
	r.Stages[STAGE2].vAy = r.Stages[BOOSTER].vAy
	r.Stages[STAGE2].vRx = r.Stages[BOOSTER].vRx
	r.Stages[STAGE2].vRy = r.Stages[BOOSTER].vRy

	r.Stages[STAGE2].DTF = r.Stages[BOOSTER].DTF

	r.Stages[STAGE2].alpha = r.Stages[BOOSTER].alpha
	r.Stages[STAGE2].beta = r.Stages[BOOSTER].beta
	r.Stages[STAGE2].gamma = r.Stages[BOOSTER].gamma

	// suddenly the booster is a lot lighter
	r.Stages[BOOSTER].Mass = r.Stages[BOOSTER].Mass - r.Stages[STAGE2].Mass
}

func (r *VEHICLE) stage_sep() {
	fmt.Println("\n************\nBooster fuel left:", r.Stages[BOOSTER].Mf, "\n************")
	r.sync_stages()
	/*	for i := 0; i < 2; i++ {
		r.Stages[i].Mass = r.Stages[i].Mr + r.Stages[i].Mf + r.Stages[i].Mp;
	}*/

	if r.NoFuel(BOOSTER) {
		fmt.Println("\n\n---------> BOOSTER EMPTY!!!! @", r.Stages[BOOSTER].Clock, "\n")
	}
	fmt.Printf("\n************\nMECO @ ---> %g seconds\n", r.Stages[BOOSTER].Clock)
	fmt.Println("Remaining fuel ...... ", r.Stages[BOOSTER].Mf, "kg")
	fmt.Println("Velocity ............ ", r.Stages[BOOSTER].RVel*3.6, "km/h")
	fmt.Println("Altitude ............ ", (r.Stages[BOOSTER].altitude)*1e-3, "km") //(-3*M_PI/2+v.Stages[STAGE2].alpha-v.Stages[STAGE2].beta)*180/M_PI)
	fmt.Println("Injection Angle ..... ", rad2deg(r.Stages[BOOSTER].gamma), "degres\n************")

	r.SysGuidance._stagesep = true
	r.EventsMap = r.EventsMap | data.E_STAGESEP
	r.LastEvent = data.E_STAGESEP
	//go r.boosterGuidance(simulation)
}

// engine cut off
func (r *VEHICLE) MSECO(stage int32, eventid uint32) bool {
	r.Stages[stage].ThrottleRate = 0
	r.Stages[stage].RunningEngines = 0
	r.EventsMap = r.EventsMap | eventid
	r.LastEvent = eventid
	return true
}

//func (r *VEHICLE) execute(event string, f *os.File) {
func (r *VEHICLE) execute(event Pevent) {
	//fmt.Println(event)
	switch event.Id {
	case "MEI":
		//		r.SysGuidance._MEI1 = r.Ignition(BOOSTER, data.EventMapping[event.Id], 1.0, 9)
		//		r.SysGuidance._MEI1 = r.Ignition(BOOSTER, data.EventInfoMapping[event.Id].Id, 1.0, 9)
		r.SysGuidance._MEI1 = r.Ignition(&event, 9)
		//output_telemetry(event, nil, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "--> Ignition Booster .....")
		break

	case "THROTTLE_DWN":
		r.Throttle(&event)
		r.EventsMap = r.EventsMap &^ (data.E_THROTTLE_U)
		fmt.Println("\t", r.Stages[event.Stage].Clock, "--> Throttling down at ", event.Rate)
		break

	case "THROTTLE_UP":
		r.Throttle(&event)
		r.EventsMap = r.EventsMap &^ (data.E_THROTTLE_D)
		fmt.Println("\t", r.Stages[event.Stage].Clock, "--> Throttling up at ", event.Rate)
		break

	case "LIFT_OFF":
		r.SysGuidance._release = r.liftOff()
		//output_telemetry(event, nil, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "--> Lift off .....")
		break

	case "PITCH":
		r.SysGuidance._pitch = r.pitchStart(&event)
		//output_telemetry(event, nil, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "--> Pitching .....", event.Gamma0)
		break

	case "MECO":
		r.SysGuidance._MECO1 = r.MSECO(BOOSTER, data.E_MECO_1)
		//output_telemetry(event, f, 1)
		//r.sync_stages()
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Main engine cut off .....")
		break

	case "STAGE_SEP":
		r.stage_sep()
		//output_telemetry(event, nil, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Stage separation .....")
		break

	case "SEI":
		//		r.SysGuidance._SEI1 = r.Ignition(STAGE2, data.E_SEI_1, 1.0, 1)
		r.SysGuidance._SEI1 = r.Ignition(&event, 1)
		//output_telemetry(event, nil, 1)
		fmt.Println("\t", r.Stages[STAGE2].Clock, "Second Stage Ignition .....")
		break

	case "ENTRY_BURN_ON": //"MEI-2":
		//		r.SysGuidance._MEI2 = r.Ignition(BOOSTER, data.E_EBURNI, 1.0, 3) //data.E_MEI_2, 3)
		r.SysGuidance._MEI2 = r.Ignition(&event, 3) //data.E_MEI_2, 3)
		//output_telemetry(event, f, 0)
		r.SysGuidance._EBURN = true // to be removed $$$
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Entry burn Starts .....")
		break

	case "ENTRY_BURN_OFF": //"MECO-2":
		r.SysGuidance._MECO2 = r.MSECO(BOOSTER, data.E_EBURNO) // data.E_MECO_2)
		//output_telemetry(event, f, 0)
		r.SysGuidance._EBURN = false // to be removed $$$
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Entry burn Stopped .....")
		break

	case "BOOSTBACK_BURN_ON": //boost back burn starts
		//		r.SysGuidance._MEI2 = r.Ignition(BOOSTER, data.E_BBURNI, 1.0, 3) // data.E_MEI_2, 3)
		r.SysGuidance._MEI2 = r.Ignition(&event, 3) // data.E_MEI_2, 3)
		//output_telemetry(event, f, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "BoostBack burn Starts .....")
		r.SysGuidance._BBURN = true
		break

	case "BOOSTBACK_BURN_OFF": // boost back burn stops
		r.SysGuidance._MECO2 = r.MSECO(BOOSTER, data.E_BBURNO) //data.E_MECO_2)
		//output_telemetry(event, f, 0)
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Boostback burn Stopped .....")
		r.SysGuidance._BBURN = false
		break

	case "SECO":
		r.SysGuidance._SECO1 = r.MSECO(STAGE2, data.E_SECO_1)
		//output_telemetry(event, f, 1);
		fmt.Println("\t", r.Stages[STAGE2].Clock, "Second stage engine cut off .....")
		break

	case "LANDING_BURN_ON": //"MEI-3":
		//		r.SysGuidance._MEI3 = r.Ignition(BOOSTER, data.E_LBURNI, 1.0, 1) //data.E_MEI_3, 1)
		r.SysGuidance._MEI3 = r.Ignition(&event, 1) //data.E_MEI_3, 1)
		//output_telemetry(event, f, 0)
		r.SysGuidance._LBURN = true // to be removed $$$
		fmt.Println("\t", r.Stages[BOOSTER].Clock, "Landing burn Started .....")
		break

	default:
		fmt.Println(event, ": unknown launch phase")
	}
}
