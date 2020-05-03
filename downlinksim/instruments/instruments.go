package instruments
import (
	//"github.com/crmathieu/daq/data"
	"unsafe"
	//"fmt"
	//"net"
	//"context"
	//"time"
)

type DataPoint struct {
    Data    unsafe.Pointer //*interface{}
	Length  uintptr
	SensorUpdate func()
	ReadSensor func() interface{}
}

var SensorsMap map[uint8]DataPoint 
var sensorIndexList []uint8 

/*
var VelocitySensor 			= data.INSTVelocitySensor{Id:data.PVELOCITY, Velx:0.0, Vely:0.0, Velz:0.0,}
var CoordinatesSensor 		= data.INSTRangeSensor{Id:data.PRANGE, Coorx:0.0, Coory:0.0, Coorz:0.0,}
var TurboPumpRPMSensor 		= data.INSTturboPumpSensor{Id:data.PTURBOPUMP, Rpm:0,}
var EnginePressureSensor 	= data.INSTenginePressureSensor{Id:data.PENGINEPRE, Pressure:0.0,}
var VolumeOxydizer          = data.PvolumeOxydizer{Id:data.PVOXIDIZER, Volume: data.MAXVOL_OXYDIZER,}
var VolumePropellant      	= data.PpropellantWeight{Id:data.PVPROPELLANT, Volume: data.MAXVOL_PROPELLANT,}
*/
/*
// instruments read
func readVelocity() interface{} {	
	return (*(*data.IVelocitySensor)(SensorsMap[data.PVELOCITY].Data))
}
func readCoordinates() interface{} {
	return (*(*data.IRangeSensor)(SensorsMap[data.PRANGE].Data))
}
func readTurboPumpRPM() interface{}  {
	return (*(*data.IturboPumpSensor)(SensorsMap[data.PTURBOPUMP].Data))
}
func readEnginePressure() interface{}  {
	return (*(*data.INSTenginePressureSensor)(SensorsMap[data.PENGINEPRE].Data))
}
*/