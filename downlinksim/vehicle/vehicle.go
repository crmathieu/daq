package main
import (
	"github.com/crmathieu/daq/data"
	"unsafe"
)
/*
type DataPoint struct {
    Data    		unsafe.Pointer //*interface{}
	Length  		uintptr
	SetSensor 	func (*VEHICLE)()
	ReadSensor 		func (*VEHICLE)() //interface{}
}
*/

type SensorHandlers struct {
	ReadSensor 	func () interface{}
	SetSensor 	func ()
}

type rocketStage struct {
	Diameter		float32
	PropellantMass	float32
	DryMass			float32
	Thrust			float32
	SpecificImp		float32
	BurnTime		float32
	M_dot			float32
	Instruments		[]unsafe.Pointer //interface{}
	Handlers		[]SensorHandlers
}

type VEHICLE struct {
	Clock, ClockAtMeco	float32
	vG, vG_dot, vD, vD_dot float64
	v_dot, x_dot, h_dot float64
	Velocity, Altitude, Range float64
	PayloadMass float32  //2482
	FrontalArea float64
	Gamma, gamma_dot float64
	Drag, G, Rho float64
	TotalStages			int8
	CurrentStage		int8
	AltitudeTarget, TicksPerSegment, TotalBurnTime, VerticalTicks float64
	EarlyTiltAngle, LateTiltAngle float64
	TargetAltitude, OrbitalVelocity float64

	Stage				[]rocketStage
}

func NewVehicle() *VEHICLE {
	var v = &VEHICLE{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 
					2482, 0.0, 
					90, 0.0, 0.0, 
					GRAVITY_ACC, RHO, 2, 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,  
					[]rocketStage{
		{}, {}},
	}

	// stage-1
	v.Stage[0].Diameter 		= F9_S1_DIAMETER
	v.Stage[0].PropellantMass 	= F9_S1_PropellantMass 
	v.Stage[0].DryMass 			= F9_S1_dryMass + F9_S2_dryMass + F9_S2_PropellantMass + v.PayloadMass
	v.Stage[0].SpecificImp		= F9_S1_ISP
	v.Stage[0].Thrust 			= F9_S1_Thrust
	v.Stage[0].BurnTime 		= F9_S1_burntime
	v.Stage[0].M_dot			= 0
	v.Stage[0].Instruments		= make([]unsafe.Pointer, data.INSTRUMENTS_COUNT)
	v.Stage[0].Handlers			= make([]SensorHandlers, data.INSTRUMENTS_COUNT)

	v.Stage[0].Instruments[data.SVELOCITY]		= (unsafe.Pointer)(&data.SENSvelocity{		Id:data.SVELOCITY, 	Velocity:0.0, Acceleration:0.0,})
	v.Stage[0].Instruments[data.SPOSITION]		= (unsafe.Pointer)(&data.SENSposition{		Id:data.SPOSITION, 	Range:0.0,    Inclinaison:0.0, Altitude:0.0,})
	v.Stage[0].Instruments[data.STURBOPUMP]		= (unsafe.Pointer)(&data.SENSturboPump{		Id:data.STURBOPUMP, Rpm:0,})
	v.Stage[0].Instruments[data.SENGINEPRE]		= (unsafe.Pointer)(&data.SENSenginePressure{Id:data.SENGINEPRE, Pressure:0.0,})
	v.Stage[0].Instruments[data.SMASSPROPELLANT]= (unsafe.Pointer)(&data.SENSpropellantMass{Id:data.SMASSPROPELLANT, Mass: F9_S1_PropellantMass,})

	v.Stage[0].Handlers[data.SVELOCITY]			= SensorHandlers{ReadSensor: v.readVelocity, 		SetSensor: v.setVelocity,}
	v.Stage[0].Handlers[data.SPOSITION]			= SensorHandlers{ReadSensor: v.readPosition, 		SetSensor: v.setPosition,}
	v.Stage[0].Handlers[data.STURBOPUMP]		= SensorHandlers{ReadSensor: v.readTurboPumpRPM,	SetSensor: v.setTurboPumpRPM,}
	v.Stage[0].Handlers[data.SENGINEPRE]		= SensorHandlers{ReadSensor: v.readEnginePressure, 	SetSensor:v.setEnginePressure,}
	v.Stage[0].Handlers[data.SMASSPROPELLANT]	= SensorHandlers{ReadSensor: v.readPropellantMass, 	SetSensor: v.setPropellantMass,}
	
	// stage-2
	v.Stage[1].Diameter 		= F9_S2_DIAMETER
	v.Stage[1].PropellantMass 	= F9_S2_PropellantMass
	v.Stage[1].DryMass 			= F9_S2_dryMass + v.PayloadMass
	v.Stage[1].SpecificImp		= F9_S2_ISP
	v.Stage[1].Thrust 			= F9_S2_Thrust
	v.Stage[1].BurnTime 		= F9_S2_burntime
	v.Stage[1].M_dot			= 0
	v.Stage[1].Instruments		= make([]unsafe.Pointer, data.INSTRUMENTS_COUNT)
	v.Stage[1].Handlers			= make([]SensorHandlers, data.INSTRUMENTS_COUNT)
	
	v.Stage[1].Instruments[data.SVELOCITY]		= (unsafe.Pointer)(&data.SENSvelocity{		Id:data.SVELOCITY, 	Velocity:0.0, Acceleration:0.0,})
	v.Stage[1].Instruments[data.SPOSITION]		= (unsafe.Pointer)(&data.SENSposition{		Id:data.SPOSITION, 	Range:0.0, 	  Inclinaison:0.0, Altitude:0.0,})
	v.Stage[1].Instruments[data.STURBOPUMP]		= (unsafe.Pointer)(&data.SENSturboPump{		Id:data.STURBOPUMP, Rpm:0,})
	v.Stage[1].Instruments[data.SENGINEPRE]		= (unsafe.Pointer)(&data.SENSenginePressure{Id:data.SENGINEPRE, Pressure:0.0,})
	v.Stage[1].Instruments[data.SMASSPROPELLANT]= (unsafe.Pointer)(&data.SENSpropellantMass{Id:data.SMASSPROPELLANT, Mass: F9_S2_PropellantMass,})

	v.Stage[1].Handlers[data.SVELOCITY]			= SensorHandlers{ReadSensor: v.readVelocity, 		SetSensor: 	v.setVelocity,}
	v.Stage[1].Handlers[data.SPOSITION]			= SensorHandlers{ReadSensor: v.readPosition, 		SetSensor: 	v.setPosition,}
	v.Stage[1].Handlers[data.STURBOPUMP]		= SensorHandlers{ReadSensor: v.readTurboPumpRPM,	SetSensor: 	v.setTurboPumpRPM,}
	v.Stage[1].Handlers[data.SENGINEPRE]		= SensorHandlers{ReadSensor: v.readEnginePressure, 	SetSensor:	v.setEnginePressure,}
	v.Stage[1].Handlers[data.SMASSPROPELLANT]	= SensorHandlers{ReadSensor: v.readPropellantMass, 	SetSensor: 	v.setPropellantMass,}

	v.setFrontalArea()
	return v
}

func (v *VEHICLE) Meco() {
	if v.CurrentStage < v.TotalStages - 1 {
		v.Stage[v.CurrentStage+1].Instruments[data.SVELOCITY] 		= v.Stage[v.CurrentStage].Instruments[data.SVELOCITY]
		v.Stage[v.CurrentStage+1].Instruments[data.SPOSITION] 	= v.Stage[v.CurrentStage].Instruments[data.SPOSITION]
		v.CurrentStage++
	}
}
/*
type SENSORUPDATE interface {
	func (v *VEHICLE) setTurboPumpRPM() 	interface{}
	func (v *VEHICLE) setEnginePressure() 	interface{}
	func (v *VEHICLE) setVelocity() 		interface{}
	func (v *VEHICLE) setPosition() 		interface{} 
}*/