package main

import (
    "io"
    "log"
    "net"
	"time"
	"github.com/crmathieu/daq/data"
	"unsafe"
	"fmt"
)

func handle(c net.Conn) {
    // Handle the reads
    start := time.Now()
    tbuf := make([]byte, 81920)
    totalBytes := 0

    for {
		n, err := c.Read(tbuf)
		log.Println(tbuf[:n])
		totalBytes += n
//		if n != 256 {
//			log.Println(n)
//		}
        // Was there an error in reading ?
        if err != nil {
            if err != io.EOF {
                log.Printf("Read error: %s", err)
            }
            break
        }
        //log.Println(n)
    }
    log.Printf("%d bytes read in %s", totalBytes, time.Now().Sub(start))
    c.Close()
}

func readPackets(c net.Conn) {
    // Handle the reads
    start := time.Now()
    tbuf := make([]byte, 81920)
    totalBytes := 0

    for {
		n, err := c.Read(tbuf)
		
		/*if n > 81920 {
			log.Printf("BUFFER OVERFLOW %d bytes read in %s", totalBytes, time.Now().Sub(start))
			c.Close()
		}*/

		//log.Println(tbuf[:n])
		detectPacketFrame(&tbuf, n)
		totalBytes += n
//		if n != 256 {
//			log.Println(n)
//		}
        // Was there an error in reading ?
        if err != nil {
            if err != io.EOF {
                log.Printf("Read error: %s", err)
            }
            break
        }
        //log.Println(n)
    }
    log.Printf("%d bytes read in %s", totalBytes, time.Now().Sub(start))
    c.Close()
}

func detectPacketFrame(pk *[]byte, size int) {
	numberDP := *(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET]))
//	fmt.Println(numberDP)
	//dp := (*[]data.SENSgeneric)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET]))
	//dp := (*[]data.SENSgeneric)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET]))

	for k:=byte(0); k<numberDP; k++ {
		dp := (*data.SENSgeneric)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET+k*data.DATAPOINT_SIZE]))
		switch(dp.Id) {
/*		case data.SVELOCITY: 	v := (*data.SENSvelocity)(unsafe.Pointer(dp))
								fmt.Println("Vel:",v.Velocity,"m/s, Acc:", v.Acceleration)
		case data.SPOSITION:	v := (*data.SENSposition)(unsafe.Pointer(dp))
								fmt.Println("Alt:", v.Altitude,"m, Range:", v.Range,"m")
		case data.STHRUST:		v := (*data.SENSthrust)(unsafe.Pointer(dp))
								fmt.Println("Thrust:", v.Thrust/1000,"kN, Stage:", v.Stage)
		case data.STILTANGLE:	v := (*data.SENStiltAngle)(unsafe.Pointer(dp))
								fmt.Println("Gamma:", v.Angle,"deg")*/
		case data.SMASSPROPELLANT:	v := (*data.SENSpropellantMass)(unsafe.Pointer(dp))
								fmt.Println("Mass:", v.Mass,"kg, Mass Flow:", v.Mflow, "kg/s, Mass Ejected:", v.Mejected)

		}
	}
}

//var UARTregister map[byte]*data.GSbuf

func main() {

    srv, err := net.Listen("tcp", data.DOWNLINK_PORT)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Listening on localhost:2000")
    for {
        conn, err := srv.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go readPackets(conn)
    }
}

/*
// ReadUART -------------------------------------------------------------------
// read the UART register or block waiting for data to arrive
// ----------------------------------------------------------------------------
func ReadUART()  {
	defer wg.Done()

	wg.Add(1)
	var item byte
	for {
		items.Wait()
		//cmutex.Enter()
		item = inputb[itail]
		readUARTRegister(itail, item)
		itail = (itail + 1) % BUFFER_SIZE
		//cmutex.Leave()
		spaces.Signal()
	}
	fmt.Printf("Consumer finished!\n")
}

const (
	STATE_INIT = 0
	STATE_SIG = 1
	STATE_PAYLOAD = 2
)

var readState = STATE_INIT
var key byte

func readUARTRegister(index int, item byte) {
	switch(readState) {
		case STATE_INIT: if item == PACKET_START {
							readState = STATE_SIG	
							//fmt.Println("START")
						} // else ignore
						break
		case STATE_SIG: key = item
						UARTregister[key].Index = 0
						readState = STATE_PAYLOAD
						tag++
						break
		case STATE_PAYLOAD: // add payload
						if UARTregister[key].Index < key {
							UARTregister[key].Buffer[UARTregister[key].Index] = item
							UARTregister[key].Index++
							if UARTregister[key].Index >= key {
								readState = STATE_INIT
								// fmt.Println(UARTregister[key].Buffer)
							}
						} else {
							if item == PACKET_START {
								readState = STATE_SIG
							} else {
								readState = STATE_INIT
							}
						}
						break
		default: break
	}
	//fmt.Printf("%s%02d -> %d\n", name, index, item)
}
*/