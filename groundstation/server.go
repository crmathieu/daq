package main

import (
    "io"
    "log"
    "net"
	"time"
    "github.com/crmathieu/daq/data"
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

var UARTregister map[byte]*data.GSbuf

func main() {

/*	UARTregister = make(map[byte]*data.GSbuf)
	UARTregister[8]  = &data.GSbuf{0, [data.PACKET_PAYLOAD_LENGTH]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},0,0} //, false} //DP1
	UARTregister[16] = &data.GSbuf{0, [data.PACKET_PAYLOAD_LENGTH]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},0,0} //, false} // DP2
	UARTregister[32] = &data.GSbuf{0, [data.PACKET_PAYLOAD_LENGTH]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},0,0} //, false} //DP3
*/
    srv, err := net.Listen("tcp", ":2000")
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Listening on localhost:2000")
    for {
        conn, err := srv.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go handle(conn)
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