package main

import (
    "log"
    "fmt"
    "net"
    "time"
    "unsafe"
    "github.com/crmathieu/daq/data"
)

var zvelocity = float32(0.2)
var altitude = float32(1.1)

func takeMeasurements() (interface{}, uintptr) {
    var payload = unsafe.Pointer(&[48]byte{})
    var offset unsafe.Offset = 0
    for mtype := range PayloadMultiplexor {
        switch mtype {
        case data.PVELOCITY:    //payload = data.Pvelocity{Velx:0, Vely:0, Velz: zvelocity};
                                p1 := PayloadMultiplexor[mtype].Data.(data.Pvelocity)
                                //payload = PayloadMultiplexor[mtype].Data.(data.Pvelocity)
                                //payload.Velz = zvelocity
                                p1.Velz = zvelocity
                                n := copy(data_send[4:], pavio_arr[:])
                                zvelocity++;
                                return payload, PayloadMultiplexor[mtype].Length

        case data.PCOORDINATES: //payload = data.Pcoordinates{Coorx:0, Coory:0, Coorz: altitude}
                                payload := PayloadMultiplexor[mtype].Data.(data.Pcoordinates)
                                payload.Coorz = altitude
                                altitude++;
                                return payload, PayloadMultiplexor[mtype].Length

        case data.PTURBOPUMP:   //payload = data.PturboPumpRPM{Rpm:15000}
                                payload := PayloadMultiplexor[mtype].Data.(data.PturboPumpRPM)
                                payload.Rpm = 15000
                                return payload, PayloadMultiplexor[mtype].Length

        case data.PENGINEPRE:   //payload = data.PenginePressure{Pressure:180}
                                payload := PayloadMultiplexor[mtype].Data.(data.PenginePressure)
                                payload.Pressure = 180
                                return payload, PayloadMultiplexor[mtype].Length
        }
        return nil, 0
    }
}

//var nullPacket = []byte{data.PACKET_START,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1}
var nullPacket = []byte{data.PACKET_START,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1}

// streamData -----------------------------------------------------------------
// stream downlink data
// ----------------------------------------------------------------------------
func streamData(c net.Conn) {
    start := time.Now()
    totalBytes := 0

    var payload interface{} 

    var length uintptr
    for i := 0; i < 1; i++ {
        for k := range PayloadMultiplexor {
            payload, length = takeMeasurements(k)
            //length = dataPointsLength[k]
            fmt.Printf("%#v\n", payload)
            pUART := setPacket(&Packet, unsafe.Pointer(&payload), k, length)
            n, err := writePacket(c, pUART)

            // Was there an error in writing?
            if err != nil {
                log.Printf("Write error: %s", err)
                continue
            }
            //log.Println(n)
            totalBytes += n
            break; /////////////////////////
        }
    }
    log.Printf("%d bytes written in %s", totalBytes, time.Now().Sub(start))
    //time.Sleep(10 * time.Second)
    c.Close()
}

func writePacket(c net.Conn, pk *[]byte) (int, error) {
    return c.Write(*pk)   
}

/*
func produceItem(c net.Conn, packet item byte) {
    n, err := c.Write(tbuf)
    totalBytes += n

    inputb[ihead] = item
	ihead = (ihead + 1) % BUFFER_SIZE
	items.Signal()
}
*/

var cnt = uint32(1)

func SendDataPoint() {}

func setPacket(pkheader *[]byte, payload unsafe.Pointer, payloadType byte, length uintptr) *[]byte {
    // reset packet START mark
    //(*pk)[0] = data.PACKET_START
    *pkheader = ([]byte)(nullPacket)
    //*(*[32]byte)(unsafe.Pointer(pk)) = nullPacket
    fmt.Printf("PK=%#v\n", *pkheader)
    *(*byte)(unsafe.Pointer(&(*pkheader)[data.PACKET_NDP_OFFSET])) = payloadType

    // cast payload as *[]byte
    pl := (*[]byte)(payload)

    // insert CRC32 calculated on payload content
    *(*uint32)(unsafe.Pointer(&(*pkheader)[data.PACKET_CRC_OFFSET])) = data.CRC32(0, pl, length) //cnt //crc

    // insert payload in packet
//    *(*[]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = *pl
 //   *(*[12]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = *pl
    *pkheader = append(*pkheader, (*pl)[:length]...)
    return pkheader
}

func setPacketSAVE(pk *[]byte, payload unsafe.Pointer, payloadType byte, length uintptr) *[]byte {
    // reset packet START mark
    //(*pk)[0] = data.PACKET_START
    *pk = ([]byte)(nullPacket)
    //*(*[32]byte)(unsafe.Pointer(pk)) = nullPacket
    fmt.Printf("PK=%#v\n", *pk)
    *(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET])) = payloadType

    // cast payload as *[]byte
    pl := (*[]byte)(payload)

    // insert CRC32 calculated on payload content
    *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) = data.CRC32(0, pl, length) //cnt //crc

    // insert payload in packet
//    *(*[]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = *pl
 //   *(*[12]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = *pl
    *pk = append(*pk, (*pl)[:12]...)
    return pk
}

// Packet - the variable holding packet content
var Packet []byte

// PayloadMultiplexor - holds the definition of each datapoint
/*var PayloadMultiplexor = map[uint8]interface{}{
    data.PVELOCITY: data.Pvelocity{},
    data.PCOORDINATES: data.Pcoordinates{},
    data.PTURBOPUMP: data.PturboPumpRPM{},
    data.PENGINEPRE: data.PenginePressure{},
}*/

var PayloadMultiplexor = map[uint8]*DP{
    data.PVELOCITY: &DP{data.Pvelocity{}, unsafe.Sizeof(data.Pvelocity{}),},
//    data.PCOORDINATES: DP{data.Pcoordinates{}, unsafe.Sizeof(data.Pcoordinates{}),},
//    data.PTURBOPUMP: DP{data.PturboPumpRPM{}, unsafe.Sizeof(data.PturboPumpRPM{}),},
//    data.PENGINEPRE: DP{data.PenginePressure{}, unsafe.Sizeof(data.PenginePressure{}),},
}

type DP struct {
    Data    interface{}
    Length  uintptr
}

// dataPointsLength - holds length
var dataPointsLength map[uint8]uintptr

func main() {

    Packet = make([]byte, data.PACKET_PAYLOAD_OFFSET) //PACKET_LENGTH)
    conn, err := net.Dial("tcp", ":2000")
    if err != nil {
        log.Fatal(err)
    }
//    var z = float32(0.0)
    fmt.Printf("---> %d\n", unsafe.Sizeof(data.Pvelocity{}))
    dataPointsLength = make(map[uint8]uintptr)
    for k, v := range PayloadMultiplexor {
        dataPointsLength[k] = unsafe.Sizeof(v)
    }

    log.Println("Sending to localhost:2000")
    streamData(conn)
}
