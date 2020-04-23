package main
import(
	"github.com/crmathieu/daq/data"
	//"time"
	"fmt"
//	"unsafe"
)	

var zvelocity = float32(0.2)
var altitude = float32(1.1)
var tbrpm = int32(15000)
var engpressure = float32(150.56)

// instruments updates
func setVelocity() {
	zvelocity++
	(*(*data.Pvelocity)(SensorsMap[data.PVELOCITY].Data)).Velz = zvelocity
}
func setCoordinates() {
	altitude++
	(*(*data.Pcoordinates)(SensorsMap[data.PCOORDINATES].Data)).Coorz = altitude
}
func setTurboPumpRPM() {
	tbrpm++
	(*(*data.PturboPumpRPM)(SensorsMap[data.PTURBOPUMP].Data)).Rpm = tbrpm
}
func setEnginePressure() {
	engpressure++
	(*(*data.PenginePressure)(SensorsMap[data.PENGINEPRE].Data)).Pressure = engpressure
}

// UpdateInstruments ----------------------------------------------------------
// calls the SensorUpdate function defined in the datapoint structure in an
// infinite loop, only interrupted by a sleep function
// ----------------------------------------------------------------------------
func UpdateInstruments() {
	for {
		//time.Sleep(0 * time.Millisecond)
		for k, v := range SensorsMap {
			v.SensorUpdate()
			fmt.Println(k, "--->", v.ReadSensor())
		}
	}
}