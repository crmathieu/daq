package main

// to install protobuf: 
// - download protobuf package
// - go install google.golang.org/protobuf/cmd/protoc-gen-go
// also: go get -u github.com/golang/protobuf/proto
import (
    "log"
    "fmt"
    "net"
    "time"
    "unsafe"
    "github.com/crmathieu/daq/data"
    "github.com/golang/protobuf/proto"
    "github.com/crmathieu/daq/daqproto"
)

//var ps = []byte{8, 16, 32}

func generatePackets(c net.Conn) {
    start := time.Now()
    //tbuf := make([]byte, 4096)
    totalBytes := 0
    var elevation float32
//    payload := make([]byte, data.PACKET_PAYLOAD_LENGTH)
//    payload := make([]byte, data.PACKET_PAYLOAD_LENGTH)



    for i := 0; i < 1000000; i++ {
        payload := &daqproto.Packet{
            DataPointID: daqproto.VELOCITY,
            DataPoint: &daqproto.Packet_Velocity {
                    Velocity: &daqproto.Packet_DPvelocity {
                        Velx: 0.0,
                        Vely: 0.0,
                        Velz: elevation,
                    },
            },
        }

        out, err := proto.Marshal(payload)
        //fmt.Printf("%x\n", payload) //*(*[]byte)(unsafe.Pointer(&payload)))
        elevation++

        //pUART := setPacket(&Packet, unsafe.Pointer(&payload))
        //n, err := writePacket(c, pUART) //&UARTregister)

        n, err := writePacket(c, &out) //&UARTregister)
        // Was there an error in writing?
        if err != nil {
            log.Printf("Write error: %s", err)
            break
        }
        //log.Println(n)
        totalBytes += n
    }
    log.Printf("%d bytes written in %s", totalBytes, time.Now().Sub(start))
    //time.Sleep(10 * time.Second)
    c.Close()
}

func generatePackets2(c net.Conn) {
    start := time.Now()
    //tbuf := make([]byte, 4096)
    totalBytes := 0

//    payload := make([]byte, data.PACKET_PAYLOAD_LENGTH)
//    payload := make([]byte, data.PACKET_PAYLOAD_LENGTH)
    payload := data.PL_DYN{}
    for i := 0; i < 2; i++ {
        payload.CoorX++
        payload.CoorY++
        payload.CoorZ++
        payload.VelX++
        payload.VelY++
        payload.VelZ++
        fmt.Printf("%x\n", payload) //*(*[]byte)(unsafe.Pointer(&payload)))

//        pUART := setPacket(&UARTregister, &payload)
        pUART := setPacket(&Packet, unsafe.Pointer(&payload))
        n, err := writePacket(c, pUART) //&UARTregister)
        // Was there an error in writing?
        if err != nil {
            log.Printf("Write error: %s", err)
            break
        }
        //log.Println(n)
        totalBytes += n
    }
    log.Printf("%d bytes written in %s", totalBytes, time.Now().Sub(start))
    //time.Sleep(10 * time.Second)
    c.Close()
}

func writePacket(c net.Conn, pk *[]byte) (int, error) {
//    return c.Write(append([]byte{PACKET_START}, UARTregister[ptype].Buffer[:ptype]...))
    return c.Write(*pk)   
}

//func writePacket2(c net.Conn, ptype byte) (int, error) {
//    return c.Write(append([]byte{PACKET_START}, UARTregister[ptype].Buffer[:ptype]...))
//}
/*
func produceItem(c net.Conn, packet item byte) {
    n, err := c.Write(tbuf)
    totalBytes += n

    inputb[ihead] = item
	ihead = (ihead + 1) % BUFFER_SIZE
	items.Signal()
}
*/

//var UARTregister map[byte]*data.GSbuf
//var UARTregister []byte
//var pl =[data.PACKET_PAYLOAD_LENGTH]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}

var cnt = uint32(1)

//type PL_DYN struct {
//    CoorX, CoorY, CoorZ float64
//    VelX, VelY, VelZ    float64
//}

//func setPacket(pk *[]byte, payload *[]byte) *[]byte {
func setPacket(pk *[]byte, payload unsafe.Pointer) *[]byte {
    //var crc = uint32(333)
    (*pk)[0] = data.PACKET_START
    pl := (*[]byte)(payload)
    //plen := len(*pl)
    *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) = 0 //data.CRC32(0, pl) //cnt //crc

//    *(*[data.PACKET_PAYLOAD_LENGTH]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = [data.PACKET_PAYLOAD_LENGTH]byte{}
//    *(*[plen]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = [plen]byte(payload)
    *(*[]byte)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) = *pl
    return pk
}

var Packet = make([]byte, data.PACKET_LENGTH)
func main() {

//    UARTregister = make([]byte, data.PACKET_LENGTH)

//    var crc = uint32(333)
//    UARTregister[0] = data.PACKET_START
//    *(*uint32)(unsafe.Pointer(&UARTregister[data.PACKET_CRC_OFFSET])) = crc
//    *(*[data.PACKET_PAYLOAD_LENGTH]byte)(unsafe.Pointer(&UARTregister[data.PACKET_PAYLOAD_OFFSET])) = [data.PACKET_PAYLOAD_LENGTH]byte{}

    conn, err := net.Dial("tcp", ":2000")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Sending to localhost:2000")
    generatePackets(conn)
}
