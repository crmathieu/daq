package main
import (
	"github.com/crmathieu/daq/data"
	"unsafe"
	//"fmt"
	"net"
//	"context"
	"time"
)

// instruments read
func (v *VEHICLE) readEnginePressure() interface{}  {
	//return (*(*data.SENSenginePressure)(v.Stage[v.CurrentStage].Sensors.EnginePressureSensor[data.SENGINEPRE].Data))
	//return v.Stage[v.CurrentStage].Instruments[data.SENGINEPRE].(data.SENSenginePressure)
	return (*(*data.SENSenginePressure)(v.Stage[v.CurrentStage].Instruments[data.SENGINEPRE]))
}
func (v *VEHICLE) readVelocity() interface{} {	
//	return (*(*data.IVelocitySensor)(v.Stage[v.CurrentStage].Sensors.VelocitySensor[data.SVELOCITY].Data))
//	return v.Stage[v.CurrentStage].Instruments[data.SVELOCITY].(data.IVelocitySensor)
	return (*(*data.SENSvelocity)(v.Stage[v.CurrentStage].Instruments[data.SVELOCITY]))
}
func (v *VEHICLE) readPosition() interface{} {
//	return (*(*data.IRangeSensor)(v.Stage[v.CurrentStage].Sensors.CoordinatesSensor[data.SPOSITION].Data))
//	return v.Stage[v.CurrentStage].Instruments[data.SPOSITION].(data.IRangeSensor)
	return (*(*data.SENSposition)(v.Stage[v.CurrentStage].Instruments[data.SPOSITION]))
}
func (v *VEHICLE) readTurboPumpRPM() interface{}  {
//	return (*(*data.IturboPumpSensor)(v.Stage[v.CurrentStage].Sensors.TurboPumpRPMSensor[data.STURBOPUMP].Data))
//	return v.Stage[v.CurrentStage].Instruments[data.STURBOPUMP].(data.IturboPumpSensor)
	return (*(*data.SENSturboPump)(v.Stage[v.CurrentStage].Instruments[data.STURBOPUMP]))
}
func (v *VEHICLE) readPropellantMass() interface{}  {
//	return (*(*data.IturboPumpSensor)(v.Stage[v.CurrentStage].Sensors.TurboPumpRPMSensor[data.STURBOPUMP].Data))
//	return v.Stage[v.CurrentStage].Instruments[data.STURBOPUMP].(data.IturboPumpSensor)
	return (*(*data.SENSpropellantMass)(v.Stage[v.CurrentStage].Instruments[data.SMASSPROPELLANT]))
}

/*
func init() {
	SensorsMap = map[uint8]DataPoint{
		data.SVELOCITY: 	DataPoint{	Data: (unsafe.Pointer)(&VelocitySensor), 		
								Length: unsafe.Sizeof(data.SENSvelocity{}), 		
								SetSensor: setVelocity,
								ReadSensor: readVelocity},
		data.SPOSITION: 	DataPoint{	Data: (unsafe.Pointer)(&CoordinatesSensor), 		
								Length: unsafe.Sizeof(data.SENSposition{}), 
								SetSensor: setPosition,
								ReadSensor: readPosition},
		data.STURBOPUMP: 	DataPoint{	Data: (unsafe.Pointer)(&TurboPumpRPMSensor), 	
								Length: unsafe.Sizeof(data.SENSturboPump{}), 
								SetSensor: setTurboPumpRPM,								
								ReadSensor: readTurboPumpRPM},
		data.SENGINEPRE: 	DataPoint{	Data: (unsafe.Pointer)(&EnginePressureSensor), 	
								Length: unsafe.Sizeof(data.SENSenginePressure{}), 
								SetSensor: setEnginePressure,			
								ReadSensor: readEnginePressure},
	}
	sensorIndexList = []uint8{data.SVELOCITY, data.SPOSITION, data.STURBOPUMP, data.SENGINEPRE}
}
*/
/*
func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

    conn, err := net.Dial("tcp", data.DOWNLINK_PORT)
    if err != nil {
		fmt.Println(err)
		return
	}
	go UpdateInstruments()

	StreamData(conn)
}
*/

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
		if capacity > cur {
			payload 	:= v.Stage[v.CurrentStage].Handlers[k].ReadSensor() //SensorsMap[sensorIndexList[k]].ReadSensor()
			unspl 		:= unsafe.Pointer(&payload)
			unsplArr 	:= *((*[data.DATAPOINT_SIZE]byte)(unspl))
			n 			:= copy(pOut[cur:], unsplArr[:])
			cur 		= cur + n 
			ndp++
		} else {
			break
		}
		k = (k + 1) % max
	}
	return ndp, cur, k  // returns #of datapoints, current offset in payload, current index in sensors
}

// StreamData -----------------------------------------------------------------
// loops in taking instruments measurement and stream readings
// ----------------------------------------------------------------------------
func (v *VEHICLE) StreamData(c net.Conn) {

	var packet = [data.PACKET_LENGTH]byte {data.PACKET_START}
	defer c.Close()

	//start := time.Now()
    //totalBytes := 0

	var ndp = byte(0)
	var sensorInd, size = 0, 0
	for {
		select {
			//case <-ctx.Done(): c.Close(); fmt.Println("Ctrl-C entered..."); return
			default: ndp, size, sensorInd = v.ReadInstruments(packet[data.PACKET_PAYLOAD_OFFSET:], data.PACKET_PAYLOAD_LENGTH, sensorInd)
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

	// set number of dapapoints in this packet
	*(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET])) = numberDP

	// set timestamp using PACKET_TT_OFFSET
    *(*int64)(unsafe.Pointer(&(*pk)[data.PACKET_TT_OFFSET])) = time.Now().UnixNano()

    // cast payload as *[]byte
	pl := (*[]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET]))

    // insert CRC32 calculated on payload content
    *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) = data.CRC32(0, pl, sizePayload)
}