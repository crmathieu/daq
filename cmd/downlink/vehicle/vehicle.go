package main
import (
	"github.com/crmathieu/daq/packages/data"
	"unsafe"
)

type SensorHandlers struct {
	ReadSensor 	func (int32) [data.DATAPOINT_SIZE]byte
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
	TotalStages			int32 //int8
	CurrentStage		int32 //int8
	AltitudeTarget, TicksPerSegment, TotalBurnTime, VerticalTicks float64
	EarlyTiltAngle, LateTiltAngle float64
	TargetAltitude, OrbitalVelocity float64

	Stage			[]rocketStage
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

	v.Stage[0].Instruments[data.SVELOCITY_OFFSET]		= (unsafe.Pointer)(&data.SENSvelocity{		Id:data.IDVELOCITY, 	Velocity:0.0, Acceleration:0.0,})
	v.Stage[0].Instruments[data.SPOSITION_OFFSET]		= (unsafe.Pointer)(&data.SENSposition{		Id:data.IDPOSITION, 	Range:0.0,    Inclinaison:0.0, Altitude:0.0,})
//	v.Stage[0].Instruments[data.STURBOPUMP_OFFSET]		= (unsafe.Pointer)(&data.SENSturboPump{		Id:data.IDTURBOPUMP, Rpm:0,})
//	v.Stage[0].Instruments[data.SENGINEPRE_OFFSET]		= (unsafe.Pointer)(&data.SENSenginePressure{Id:data.IDENGINEPRE, Pressure:0.0,})
	v.Stage[0].Instruments[data.STILTANGLE_OFFSET]		= (unsafe.Pointer)(&data.SENStiltAngle{		Id:data.IDTILTANGLE, Angle:0,})
	v.Stage[0].Instruments[data.STHRUST_OFFSET]			= (unsafe.Pointer)(&data.SENSthrust{		Id:data.IDTHRUST, Thrust:0,})
	v.Stage[0].Instruments[data.SMASSPROPELLANT_OFFSET]	= (unsafe.Pointer)(&data.SENSpropellantMass{Id:data.IDMASSPROPELLANT, Mflow: 0.0, Mass: F9_S1_PropellantMass,})

	v.Stage[0].Handlers[data.SVELOCITY_OFFSET]			= SensorHandlers{ReadSensor: v.readVelocity,}
	v.Stage[0].Handlers[data.SPOSITION_OFFSET]			= SensorHandlers{ReadSensor: v.readPosition, }
//	v.Stage[0].Handlers[data.STURBOPUMP_OFFSET]			= SensorHandlers{ReadSensor: v.readTurboPumpRPM,}
//	v.Stage[0].Handlers[data.SENGINEPRE_OFFSET]			= SensorHandlers{ReadSensor: v.readEnginePressure,}
	v.Stage[0].Handlers[data.STILTANGLE_OFFSET]			= SensorHandlers{ReadSensor: v.readTiltAngle,}
	v.Stage[0].Handlers[data.STHRUST_OFFSET]			= SensorHandlers{ReadSensor: v.readThrust,}
	v.Stage[0].Handlers[data.SMASSPROPELLANT_OFFSET]	= SensorHandlers{ReadSensor: v.readPropellantMass, }
	
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
	
	v.Stage[1].Instruments[data.SVELOCITY_OFFSET]		= (unsafe.Pointer)(&data.SENSvelocity{		Id:data.IDVELOCITY, 	Velocity:0.0, Acceleration:0.0,})
	v.Stage[1].Instruments[data.SPOSITION_OFFSET]		= (unsafe.Pointer)(&data.SENSposition{		Id:data.IDPOSITION, 	Range:0.0, 	  Inclinaison:0.0, Altitude:0.0,})
//	v.Stage[1].Instruments[data.STURBOPUMP_OFFSET]		= (unsafe.Pointer)(&data.SENSturboPump{		Id:data.IDTURBOPUMP, Rpm:0,})
//	v.Stage[1].Instruments[data.SENGINEPRE_OFFSET]		= (unsafe.Pointer)(&data.SENSenginePressure{Id:data.IDENGINEPRE, Pressure:0.0,})
	v.Stage[1].Instruments[data.STILTANGLE_OFFSET]		= (unsafe.Pointer)(&data.SENStiltAngle{		Id:data.IDTILTANGLE, Angle:0,})
	v.Stage[1].Instruments[data.STHRUST_OFFSET]			= (unsafe.Pointer)(&data.SENSthrust{		Id:data.IDTHRUST, Thrust:0,})
	v.Stage[1].Instruments[data.SMASSPROPELLANT_OFFSET]	= (unsafe.Pointer)(&data.SENSpropellantMass{Id:data.IDMASSPROPELLANT, Mass: F9_S2_PropellantMass,})

	v.Stage[1].Handlers[data.SVELOCITY_OFFSET]			= SensorHandlers{ReadSensor: v.readVelocity, }
	v.Stage[1].Handlers[data.SPOSITION_OFFSET]			= SensorHandlers{ReadSensor: v.readPosition, }
//	v.Stage[1].Handlers[data.STURBOPUMP_OFFSET]			= SensorHandlers{ReadSensor: v.readTurboPumpRPM,}
//	v.Stage[1].Handlers[data.SENGINEPRE_OFFSET]			= SensorHandlers{ReadSensor: v.readEnginePressure,}
	v.Stage[1].Handlers[data.STILTANGLE_OFFSET]			= SensorHandlers{ReadSensor: v.readTiltAngle, }
	v.Stage[1].Handlers[data.STHRUST_OFFSET]			= SensorHandlers{ReadSensor: v.readThrust,}
	v.Stage[1].Handlers[data.SMASSPROPELLANT_OFFSET]	= SensorHandlers{ReadSensor: v.readPropellantMass,}

	v.setFrontalArea()
	return v
}

func (v *VEHICLE) Meco() {
	if v.CurrentStage < v.TotalStages - 1 {
		v.Stage[v.CurrentStage+1].Instruments[data.SVELOCITY_OFFSET] 	= v.Stage[v.CurrentStage].Instruments[data.SVELOCITY_OFFSET]
		v.Stage[v.CurrentStage+1].Instruments[data.SPOSITION_OFFSET] 	= v.Stage[v.CurrentStage].Instruments[data.SPOSITION_OFFSET]
		v.CurrentStage++
	}
}
