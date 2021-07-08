package main

import (
	"unsafe"

	"github.com/crmathieu/daq/packages/data"
	//	"fmt"
	"net"
	//	"context"
	"time"
)

func (v *VEHICLE) readAngles(stage int32) [data.DATAPOINT_SIZE]byte { //interface{}  {
	//p := (*data.SENStiltAngle)(v.Stage[v.CurrentStage].Instruments[data.STILTANGLE_OFFSET])
	p := (*data.SENSangles)(v.Instruments[data.SANGLES_OFFSET])
	p.Id = Mux(stage, p.Id)
	p.Gamma = float32(v.Stages[stage].gamma)
	p.Beta = float32(v.Stages[stage].beta)
	//	p.Alpha = float32(v.Stages[stage].alpha)
	p.Zeta = float32(v.Stages[stage].zeta)
	//p.RateOfChange = float32(v.gamma_dot)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
	//	return [data.DATAPOINT_SIZE]byte{}
}

func (v *VEHICLE) readThrust(stage int32) [data.DATAPOINT_SIZE]byte { //interface{}  {
	p := (*data.SENSthrust)(v.Instruments[data.STHRUST_OFFSET])
	p.Id = Mux(stage, p.Id)
	p.Thrust = float32(v.Stages[stage].Thrust)
	//p.Stage = stage
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readVelocity(stage int32) [data.DATAPOINT_SIZE]byte { //interface{}  {
	p := (*data.SENSvelocity)(v.Instruments[data.SVELOCITY_OFFSET])
	p.Id = Mux(stage, p.Id)
	p.Velocity = float32((v.Stages[stage].VRelative) * 1e-3)
	p.Acceleration = float32(v.Stages[stage].Acc / g0)
	//p.Stage = uint32(stage)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readPosition(stage int32) [data.DATAPOINT_SIZE]byte { //interface{}  {
	p := (*data.SENSposition)(v.Instruments[data.SPOSITION_OFFSET])
	p.Id = Mux(stage, p.Id)
	//	p.Range = float32((v.Stages[stage].cx)*1e-3) //float32(v.Range)
	p.Range = float32((v.Stages[stage].drange) * 1e-3) //float32(v.Range)
	//	p.Altitude = float32((v.Stages[stage].DTF - Re) * 1e-3)

	//	p.Altitude = float32((v.Stages[stage].cy - Re)*1e-3)
	p.Altitude = float32((v.Stages[stage].altitude) * 1e-3)

	//p.Inclinaison = float32(0.0)
	//p.Stage = uint32(stage)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readPropellantMass(stage int32) [data.DATAPOINT_SIZE]byte { //interface{}  {
	p := (*data.SENSpropellantMass)(v.Instruments[data.SMASSPROPELLANT_OFFSET])
	//	p.Mflow = v.Stages[stage].M_dot * THROTTLE_VALUE
	//	p.Mass = v.Stages[stage].PropellantMass + v.Stage[v.CurrentStage].DryMass - p.Mflow * (v.Clock - v.ClockAtMeco)
	//	p.Mejected = p.Mflow * v.Clock
	//	p.Stage = uint32(stage) //v.CurrentStage + 1
	p.Id = Mux(stage, p.Id)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readEvent(stage int32) [data.DATAPOINT_SIZE]byte {
	p := (*data.SENSevent)(v.Instruments[data.SEVENT_OFFSET])
	p.Id = Mux(stage, p.Id)
	p.EventId = v.LastEvent
	p.Time = float32(v.Stages[stage].Clock)
	p.EventMap = v.EventsMap
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func (v *VEHICLE) readTime(stage int32) [data.DATAPOINT_SIZE]byte {
	p := (*data.SENStime)(v.Instruments[data.STIME_OFFSET])
	p.Id = Mux(stage, p.Id)
	p.Time = float32(v.Stages[stage].Clock)
	return *(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(p))
}

func Mux(stage int32, id uint32) uint32 {
	return (uint32(stage) << 16) | (id & 0xffff)
}

// ReadInstruments ------------------------------------------------------------
// GOROUTINE - performs instruments reading
// returns:
//	- the number of datapoints measured
//	- the size of the data
//	- the index of the next instrument to measure in the list of instruments
//  - the stage of the next instrument to measure in the list of instruments
// ----------------------------------------------------------------------------
func (v *VEHICLE) ReadInstruments(pOut []byte, capacity, index int, stage int) (byte, int, int, int) {
	var cur = 0
	var ndp = byte(0)
	var k = index
	max := len(v.Instruments)
	for {
		if v.Handlers[k].ReadSensor != nil {
			payload := v.Handlers[k].ReadSensor(int32(stage))
			n := copy((pOut)[cur:], payload[:])
			cur = cur + n
			ndp++
		}
		if k+1 >= max {
			stage = (stage + 1) % 2
		}
		k = (k + 1) % max
		if capacity-data.DATAPOINT_SIZE < cur {
			return ndp, cur, k, stage // returns #of datapoints, current offset in payload, current index in sensors
		}
	}
}

// StreamData -----------------------------------------------------------------
// Instruments reading and streaming infinite loop. Paused every 10 millisec
// ----------------------------------------------------------------------------
func (v *VEHICLE) StreamData(c net.Conn) {

	var packet = [data.PACKET_LENGTH]byte{data.PACKET_START1, data.PACKET_START2}
	defer c.Close()

	var ndp = byte(0)
	var sensorInd, size, stage = 0, 0, 0
	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			ndp, size, sensorInd, stage = v.ReadInstruments(packet[data.PACKET_PAYLOAD_OFFSET:], data.PACKET_PAYLOAD_LENGTH, sensorInd, stage)
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
