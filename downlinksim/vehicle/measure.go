package main
import (
	"github.com/crmathieu/daq/data"
	"unsafe"
//	"fmt"
	"net"
//	"context"
	"time"
)

// instruments read
/*func (v *VEHICLE) readEnginePressure() interface{}  {
	//return (*(*data.SENSenginePressure)(v.Stage[v.CurrentStage].Sensors.EnginePressureSensor[data.SENGINEPRE].Data))
	//return v.Stage[v.CurrentStage].Instruments[data.SENGINEPRE].(data.SENSenginePressure)
	return (*(*data.SENSenginePressure)(v.Stage[v.CurrentStage].Instruments[data.SENGINEPRE]))
}*/


func (v *VEHICLE) readTiltAngle() [data.DATAPOINT_SIZE]byte {	//interface{}  {
	p := (*data.SENStiltAngle)(v.Stage[v.CurrentStage].Instruments[data.STILTANGLE])
	p.Angle = float32(v.Gamma * rad)
	p.RateOfChange = float32(v.gamma_dot)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readThrust() [data.DATAPOINT_SIZE]byte {	//interface{}  {
	p := (*data.SENSthrust)(v.Stage[v.CurrentStage].Instruments[data.STHRUST])
	p.Thrust = float32(v.Stage[v.CurrentStage].Thrust)
	p.Stage = v.CurrentStage + 1
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readVelocity() [data.DATAPOINT_SIZE]byte {	//interface{}  {
	p := (*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY])
	p.Velocity = float32(v.Velocity)
	p.Acceleration = float32(1.0)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readPosition() [data.DATAPOINT_SIZE]byte {	//interface{}  {
	p := (*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION])
	p.Range = float32(v.Range)
	p.Altitude = float32(v.Altitude)
	p.Inclinaison = float32(0.0)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readPropellantMass() [data.DATAPOINT_SIZE]byte {	//interface{}  {
	p := (*data.SENSpropellantMass)(v.Stage[v.CurrentStage].Instruments[data.SMASSPROPELLANT])
	p.Mflow = v.Stage[v.CurrentStage].M_dot * THROTTLE_VALUE
	p.Mass = v.Stage[v.CurrentStage].PropellantMass + v.Stage[v.CurrentStage].DryMass - p.Mflow * (v.Clock - v.ClockAtMeco)
	p.Mejected = p.Mflow * v.Clock
//	p.Stage = v.CurrentStage + 1
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

// ReadInstruments ------------------------------------------------------------
// this is called as a goroutine to perform updates
// returns:
//	- the number of datapoints measured
//	- the size of the data 
//	- the index of the next instrument to measure in the list of instruments
// ----------------------------------------------------------------------------
func (v *VEHICLE)ReadInstruments(pOut []byte, capacity, index int) (byte, int, int) { 
	var cur =  0
	var ndp = byte(0)
	var k = index
	max := len(v.Stage[v.CurrentStage].Instruments)
	for {
		payload 	:= v.Stage[v.CurrentStage].Handlers[k].ReadSensor() //SensorsMap[sensorIndexList[k]].ReadSensor()
		n 			:= copy((pOut)[cur:], payload[:])
		cur 		= cur + n 
		ndp++
		k = (k + 1) % max
		if capacity - data.DATAPOINT_SIZE < cur {
			break
		}
	}
	return ndp, cur, k  // returns #of datapoints, current offset in payload, current index in sensors
}

// StreamData -----------------------------------------------------------------
// loops in taking instruments measurement and stream readings
// ----------------------------------------------------------------------------
func (v *VEHICLE) StreamData(c net.Conn) {

	var packet = [data.PACKET_LENGTH]byte {data.PACKET_START1, data.PACKET_START2}
	defer c.Close()

	var ndp = byte(0)
	var sensorInd, size = 0, 0
	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
			case <-ticker.C: 
				 ndp, size, sensorInd = v.ReadInstruments(packet[data.PACKET_PAYLOAD_OFFSET:], data.PACKET_PAYLOAD_LENGTH, sensorInd)
				 setPacket(&packet, ndp, size)
				 writePacket(c, (*[data.PACKET_LENGTH]byte)(unsafe.Pointer(&packet)))
		}
	}	
}

// writePacket ----------------------------------------------------------------
// writes a packet to an established connection
// ----------------------------------------------------------------------------
func writePacket(c net.Conn, pk *[data.PACKET_LENGTH]byte) (int, error) {
	//fmt.Println("------>", (*pk))
    return c.Write((*pk)[:data.PACKET_LENGTH])   
}

// setPacket ------------------------------------------------------------------
// sets a packet header and body by inserting:
// 		- the number of datapoints in the packet
//		- the timestamp
//		- the CRC calculated on the payload
// ----------------------------------------------------------------------------
func setPacket(pk *[data.PACKET_LENGTH]byte, numberDP byte, sizePayload int) { 
	// set number of datapoints in this packet
	*(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET])) = numberDP

	// set timestamp using PACKET_TT_OFFSET
    *(*int64)(unsafe.Pointer(&(*pk)[data.PACKET_TT_OFFSET])) = time.Now().UnixNano()

	// set CRC 
    *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) = data.CRC32(0, (*pk)[data.PACKET_PAYLOAD_OFFSET:], sizePayload)
}