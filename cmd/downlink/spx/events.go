
package main
import (
//	"math"
	"os"
	"fmt"
)

func ignition(stage int32, num_engs int32) bool {
	F9.Stages[stage].ThrottleRate = 1.0
	F9.Stages[stage].RunningEngines = num_engs
	//t = 0.0 //F9.Clock = 0.0
	F9.Stages[STAGE2].Clock = F9.Stages[BOOSTER].Clock
	return true
}

func MSECO(stage int32) bool {
	F9.Stages[stage].ThrottleRate = 0
	F9.Stages[stage].RunningEngines = 0
	return true
}

func sync_stages() {
	F9.Stages[STAGE2].cx = F9.Stages[BOOSTER].cx
	F9.Stages[STAGE2].cy = F9.Stages[BOOSTER].cy

	F9.Stages[STAGE2].vAx = F9.Stages[BOOSTER].vAx
	F9.Stages[STAGE2].vAy = F9.Stages[BOOSTER].vAy

	F9.Stages[STAGE2].PolarDistance = F9.Stages[BOOSTER].PolarDistance
	F9.Stages[STAGE2].VAbsolute = F9.Stages[BOOSTER].VAbsolute

	F9.Stages[STAGE2].alpha = F9.Stages[BOOSTER].alpha
	F9.Stages[STAGE2].beta = F9.Stages[BOOSTER].beta
	F9.Stages[STAGE2].gam = F9.Stages[BOOSTER].gam
}

func stage_sep() {
	for i := 0; i < 2; i++ {
		F9.Stages[i].Mass = F9.Stages[i].Mr + F9.Stages[i].Mf + F9.Stages[i].Mp;
	}
}

func execute(event string, f *os.File) {
	//fmt.Println(event)
	switch event {
	case "MEI-1":
		_MEI1 = ignition(0, 9)
		//output_telemetry(event, nil, 0)
		break
	case  "Liftoff":
		_release = first_step()
		//output_telemetry(event, nil, 0)
		break
	case "Pitch_Kick":
		_pitch = pitch_kick()
		//output_telemetry(event, nil, 0)
		break
	case "MECO-1":
		_MECO1 = MSECO(0)
		sync_stages()
		//output_telemetry(event, f, 1)
		break
	case "Stage_Sep":
		stage_sep()
		//output_telemetry(event, nil, 0)
		break
	case "SEI-1":
		_SEI1 = ignition(1, 1)
		//output_telemetry(event, nil, 1)
		break
	case "MEI-2":
		_MEI2 = ignition(0, 3)
		_BBURN = true
		//output_telemetry(event, f, 0)
		break
	case "MECO-2":
		_MECO2 = MSECO(0)
		_BBURN = false
		//output_telemetry(event, f, 0)
		break
	case "SECO-1":
		_SECO1 = MSECO(1)
		//output_telemetry(event, f, 1);
		break
	case "MEI-3":
		_MEI3 = ignition(0, 1)
		_LBURN = true
		//output_telemetry(event, f, 0)
		break
	default:
		fmt.Println(event,": unknown launch phase")
	}
}

