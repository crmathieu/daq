//package measurements
package main 
import (
	"github.com/crmathieu/daq/data"
	"unsafe"
	"fmt"
	"net"
	"context"
	"time"
)

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


func init() {
	SensorsMap = map[uint8]DataPoint{
		data.PVELOCITY: 	DataPoint{	Data: (unsafe.Pointer)(&VelocitySensor), 		
								Length: unsafe.Sizeof(data.Pvelocity{}), 		
								SensorUpdate: setVelocity,
								ReadSensor: readVelocity},
		data.PCOORDINATES: 	DataPoint{	Data: (unsafe.Pointer)(&CoordinatesSensor), 		
								Length: unsafe.Sizeof(data.Pcoordinates{}), 
								SensorUpdate: setCoordinates,
								ReadSensor: readCoordinates},
		data.PTURBOPUMP: 	DataPoint{	Data: (unsafe.Pointer)(&TurboPumpRPMSensor), 	
								Length: unsafe.Sizeof(data.PturboPumpRPM{}), 
								SensorUpdate: setTurboPumpRPM,								
								ReadSensor: readTurboPumpRPM},
		data.PENGINEPRE: 	DataPoint{	Data: (unsafe.Pointer)(&EnginePressureSensor), 	
								Length: unsafe.Sizeof(data.PenginePressure{}), 
								SensorUpdate: setEnginePressure,			
								ReadSensor: readEnginePressure},
	}
	sensorIndexList = []uint8{data.PVELOCITY, data.PCOORDINATES, data.PTURBOPUMP, data.PENGINEPRE}
}

func main() {
    conn, err := net.Dial("tcp", ":2000")
    if err != nil {
		fmt.Println(err)
		return
	}
	go UpdateInstruments()

	StreamData(conn)
}

// ReadInstruments ------------------------------------------------------------
// this is called as a goroutine to perform updates
// returns:
//	- the number of datapoints measured
//	- the size of the data 
//	- the index of the next instrument to measure in the list of instruments
// ----------------------------------------------------------------------------
func ReadInstruments(pOut []byte, capacity, index int) (byte, int, int) { 
	var cur =  0
	var ndp = byte(0)
	var k = index

	max := len(sensorIndexList)
	for {
		if capacity > cur {
			payload 	:= SensorsMap[sensorIndexList[k]].ReadSensor()
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
func StreamData(c net.Conn) {

	var packet = [data.PACKET_LENGTH]byte {data.PACKET_START}
//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//start := time.Now()
    //totalBytes := 0

	var ndp = byte(0)
	var sensorInd, size = 0, 0
	for {
		select {
			case <-ctx.Done(): c.Close(); fmt.Println("Ctrl-C entered..."); return
			default: ndp, size, sensorInd = ReadInstruments(packet[data.PACKET_PAYLOAD_OFFSET:], data.PACKET_PAYLOAD_LENGTH, sensorInd)
					 setPacket(&packet, ndp, size)
					 writePacket(c, (*[data.PACKET_LENGTH]byte)(unsafe.Pointer(&packet)))
		}
	}	
}

// writePacket ----------------------------------------------------------------
// writes a packet to an established connection
// ----------------------------------------------------------------------------
func writePacket(c net.Conn, pk *[data.PACKET_LENGTH]byte) (int, error) {
//	fmt.Println("------>", (*pk))
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