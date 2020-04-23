package main
import (
	"github.com/crmathieu/daq/data"
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

var VelocitySensor 			= data.Pvelocity{Id:data.PVELOCITY, Velx:0.0, Vely:0.0, Velz:0.0,}
var CoordinatesSensor 		= data.Pcoordinates{Id:data.PCOORDINATES, Coorx:0.0, Coory:0.0, Coorz:0.0,}
var TurboPumpRPMSensor 		= data.PturboPumpRPM{Id:data.PTURBOPUMP, Rpm:0,}
var EnginePressureSensor 	= data.PenginePressure{Id:data.PENGINEPRE, Pressure:0.0,}

/*
// instruments read
func readVelocity() interface{} {	
	return (*(*data.Pvelocity)(SensorsMap[data.PVELOCITY].Data))
}
func readCoordinates() interface{} {
	return (*(*data.Pcoordinates)(SensorsMap[data.PCOORDINATES].Data))
}
func readTurboPumpRPM() interface{}  {
	return (*(*data.PturboPumpRPM)(SensorsMap[data.PTURBOPUMP].Data))
}
func readEnginePressure() interface{}  {
	return (*(*data.PenginePressure)(SensorsMap[data.PENGINEPRE].Data))
}
*/