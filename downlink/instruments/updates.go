package instruments
import(
	"github.com/crmathieu/daq/data"
	"time"
	//"fmt"
//	"unsafe"
)	

var zvelocity = float32(0.2)
var altitude = float32(1.1)
var tbrpm = int32(15000)
var engpressure = float32(150.56)

// instruments updates
func setVelocity() {
	zvelocity++
	(*(*data.INSTVelocitySensor)(SensorsMap[data.PVELOCITY].Data)).Velocity = zvelocity
}
func setCoordinates() {
	altitude++
	(*(*data.INSTRangeSensor)(SensorsMap[data.PRANGE].Data)).Altitude = altitude
}
func setTurboPumpRPM() {
	tbrpm++
	(*(*data.INSTturboPumpSensor)(SensorsMap[data.PTURBOPUMP].Data)).Rpm = tbrpm
}
func setEnginePressure() {
	engpressure++
	(*(*data.INSTenginePressureSensor)(SensorsMap[data.PENGINEPRE].Data)).Pressure = engpressure
}

var UPDATE_TICK = 10 * time.Millisecond

// RunInstrumentsUpdate -------------------------------------------------------
// calls the SensorUpdate function defined in the datapoint structure in an
// infinite loop, only interrupted by a sleep function
// ----------------------------------------------------------------------------
func RunInstrumentsUpdate() {
	for {
		time.Sleep(UPDATE_TICK)
		for _, v := range SensorsMap {
			v.SensorUpdate()
			//fmt.Println(k, "--->", v.ReadSensor())
		}
	}
}