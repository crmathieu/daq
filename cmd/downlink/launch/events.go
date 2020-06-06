
package main
import (
//	"math"
//	"os"
	"fmt"
	"github.com/crmathieu/daq/packages/data"
)

func (v *VEHICLE) hasEvent(event uint32) bool {
	if v.EventsMap & event != 0 {
		return true
	}
	return false
}

//inline void pitchStart(int *x)
//func pitchStart(x *int32) {
//	gamma[0] = (M_PI / 2) - 0.025;
//	*x = 1;
//}
func (r *VEHICLE) pitchStart() bool {
	r.Stages[BOOSTER].gamma = (M_PI / 2) - 0.025
	r.EventsMap = r.EventsMap | data.E_STARTPITCH
	r.LastEvent = data.E_STARTPITCH
	return true
}

func (r *VEHICLE) Throttle(event *Event) bool {
	r.Stages[event.Stage].ThrottleRate = float64(event.Rate)/100
	r.EventsMap = r.EventsMap | data.EventMapping[event.Id]
	r.LastEvent = data.EventMapping[event.Id]
	return true
}

func (r *VEHICLE) Ignition(stage int32, event uint32, num_engs int32) bool {
	r.Stages[stage].ThrottleRate = 1.0
	r.Stages[stage].RunningEngines = num_engs
	r.EventsMap = r.EventsMap | event
	r.LastEvent = event
	return true
}

func (r *VEHICLE) MSECO(stage int32, event uint32) bool {
	r.Stages[stage].ThrottleRate = 0
	r.Stages[stage].RunningEngines = 0
	r.EventsMap = r.EventsMap | event
	r.LastEvent = event
	return true
}

func (r *VEHICLE) sync_stages() {
	r.Stages[STAGE2].cx = r.Stages[BOOSTER].cx
	r.Stages[STAGE2].cy = r.Stages[BOOSTER].cy

	r.Stages[STAGE2].vAx = r.Stages[BOOSTER].vAx
	r.Stages[STAGE2].vAy = r.Stages[BOOSTER].vAy

	r.Stages[STAGE2].DTF = r.Stages[BOOSTER].DTF
	r.Stages[STAGE2].VAbsolute = r.Stages[BOOSTER].VAbsolute

	r.Stages[STAGE2].alpha = r.Stages[BOOSTER].alpha
	r.Stages[STAGE2].beta = r.Stages[BOOSTER].beta
	r.Stages[STAGE2].gamma = r.Stages[BOOSTER].gamma

//	r.Stages[STAGE2].Clock = r.Stages[BOOSTER].Clock

}

func (r *VEHICLE) stage_sep() {
	for i := 0; i < 2; i++ {
		r.Stages[i].Mass = r.Stages[i].Mr + r.Stages[i].Mf + r.Stages[i].Mp;
	}
	r.SysGuidance._stagesep = true
	r.EventsMap = r.EventsMap | data.E_STAGESEP
	r.LastEvent = data.E_STAGESEP
	go r.boosterGuidance()
}

//func (r *VEHICLE) execute(event string, f *os.File) {
func (r *VEHICLE) execute(event Event) {
	//fmt.Println(event)
	switch event.Id {
	case "MEI":
		r.SysGuidance._MEI1 = r.Ignition(BOOSTER, data.EventMapping[event.Id], 9)
		//output_telemetry(event, nil, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"--> Ignition Booster .....")
		break

	case "THROTTLE_DWN":
		r.Throttle(&event)
		r.EventsMap = r.EventsMap &^ (data.E_THROTTLE_U)
		fmt.Println("\t",r.Stages[event.Stage].Clock,"--> Throttling down at ",event.Rate)
		break

	case "THROTTLE_UP":
		r.Throttle(&event)
		r.EventsMap = r.EventsMap &^ (data.E_THROTTLE_D)
		fmt.Println("\t",r.Stages[event.Stage].Clock,"--> Throttling up at ",event.Rate)
		break

	case  "LIFT_OFF":
		r.SysGuidance._release = r.liftOff()
		//output_telemetry(event, nil, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"--> Lift off .....")
		break

	case "PITCH":
		r.SysGuidance._pitch = r.pitchStart()
		//output_telemetry(event, nil, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"--> Pitching .....")
		break

	case "MECO":
		r.SysGuidance._MECO1 = r.MSECO(BOOSTER, data.E_MECO_1)
		//output_telemetry(event, f, 1)
		r.sync_stages()
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Main engine cut off .....")
		break

	case "STAGE_SEP":
		r.stage_sep()
		//output_telemetry(event, nil, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Stage separation .....")
		break

	case "SEI":
		r.SysGuidance._SEI1 = r.Ignition(STAGE2, data.E_SEI_1, 1)
		//output_telemetry(event, nil, 1)
		fmt.Println("\t",r.Stages[STAGE2].Clock,"Second Stage Ignition .....")
		break

	case "ENTRY_BURN_ON": //"MEI-2":
		r.SysGuidance._MEI2 = r.Ignition(BOOSTER, data.E_EBURNI, 3) //data.E_MEI_2, 3)
		//output_telemetry(event, f, 0)
		r.SysGuidance._EBURN = true // to be removed $$$
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Entry burn Starts .....")
		break

	case "ENTRY_BURN_OFF": //"MECO-2":
		r.SysGuidance._MECO2 = r.MSECO(BOOSTER, data.E_EBURNO) // data.E_MECO_2)
		//output_telemetry(event, f, 0)
		r.SysGuidance._EBURN = false // to be removed $$$ 
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Entry burn Stopped .....")
		break

	case "BOOSTBACK_BURN_ON": //boost back burn starts
		r.SysGuidance._MEI2 = r.Ignition(BOOSTER, data.E_BBURNI, 3) // data.E_MEI_2, 3)
		//output_telemetry(event, f, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"BoostBack burn Starts .....")
		r.SysGuidance._BBURN = true
		break

	case "BOOSTBACK_BURN_OFF": // boost back burn stops
		r.SysGuidance._MECO2 = r.MSECO(BOOSTER, data.E_BBURNO) //data.E_MECO_2)
		//output_telemetry(event, f, 0)
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Boostback burn Stopped .....")
		r.SysGuidance._BBURN = false
		break

	case "SECO":
		r.SysGuidance._SECO1 = r.MSECO(STAGE2, data.E_SECO_1)
		//output_telemetry(event, f, 1);
		fmt.Println("\t",r.Stages[STAGE2].Clock,"Second stage engine cut off .....")
		break

	case "LANDING_BURN_ON": //"MEI-3":
		r.SysGuidance._MEI3 = r.Ignition(BOOSTER, data.E_LBURNI, 1) //data.E_MEI_3, 1)
		//output_telemetry(event, f, 0)
		r.SysGuidance._LBURN = true // to be removed $$$
		fmt.Println("\t",r.Stages[BOOSTER].Clock,"Landing burn Started .....")
		break

	default:
		fmt.Println(event,": unknown launch phase")
	}
}

