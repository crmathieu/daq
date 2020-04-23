package main


// testing synchronization of multiple consumers and 1 producer on a shared resource
import (
	"fmt"
	"sync"
	"time"
	sem "github.com/crmathieu/gosem/pkg/semaphore"
)

var items, spaces *sem.Sem
var cmutex *sem.Mutex

const BUFFER_SIZE = 256
var wg sync.WaitGroup

var PACKET_START byte = byte(0xff)
type DP1 struct {
	Size byte
	Payload [8]byte
}

type DP2 struct {
	Size byte
	Payload [16]byte
}

type DP3 struct {
	Size byte
	Payload [32]byte
}

type DP struct {
	D1 DP1
	D2 DP2
	D3 DP3
}

var dataPackets = DP{	D1: DP1{Size: 8,  Payload: [8]byte{'0','1','2','3','4','5','6','7'}}, 
						D2: DP2{Size: 16, Payload: [16]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}}, 
						D3: DP3{Size: 32, Payload: [32]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f','0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}}}

var inputb [BUFFER_SIZE]byte
var ihead, itail = 0, 0

var buffer1 [BUFFER_SIZE]DP1
var head1, tail1 = 0, 0
var buffer2 [BUFFER_SIZE]DP2
var head2, tail2 = 0, 0
var buffer3 [BUFFER_SIZE]DP3
var head3, tail3 = 0, 0

/*
func main2() {

	wg.Add(1)
	go func() int {
		defer wg.Done()
		fmt.Println("starting subroutine")
		var i = 0
		ticker := time.Tick(1 * time.Millisecond)
		for {
			select {
			case <-ticker: fmt.Println("out"); fmt.Println(i);return i
			default: i++ //fmt.Println(i); i++
			}
		}
	}()
	wg.Wait()
	fmt.Println("end...")
}
*/

/*var GSdp1 DP1
var GSdp2 DP2
var GSdp3 DP3
*/
type GSbuf struct {
	Index byte
	Buffer [32]byte
	Ready bool
}

var UARTregister map[byte]*GSbuf

// main------------------------------------------------------------------------
func main() {
	
	cmutex = sem.Createmutex("consumermutex")
	items = sem.Createsem("usedcount", BUFFER_SIZE, 0)
	spaces = sem.Createsem("availablecount", BUFFER_SIZE, BUFFER_SIZE)
	UARTregister = make(map[byte]*GSbuf)
	UARTregister[8]  = &GSbuf{0, [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, false} //DP1
	UARTregister[16] = &GSbuf{0, [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, false} // DP2
	UARTregister[32] = &GSbuf{0, [32]byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}, false} //DP3

	go downlink()
	go ReadUART()

/*	go consumer("B")
	go consumer("C")
	go consumer("D")
	go consumer("E")
	go consumer("F")
	go consumer("G")
	go consumer("H")
	go consumer("I")
*/
	// wait a bit to give a chance to the goroutines to start
	time.Sleep(1 * time.Second)
	wg.Wait()
}

var tag = 0
// downlink --------------------------------------------------------------------
func downlink() {
	defer wg.Done()
	//var tag int
	wg.Add(1)
	ticker := time.Tick(1 * time.Second)
	for {

		select {
		case <-ticker: fmt.Println("out"); fmt.Println(tag);return 
		default:
				//fmt.Println(".")
				produceItem(PACKET_START)
				produceItem(dataPackets.D1.Size)
				for k := range dataPackets.D1.Payload {
					produceItem(dataPackets.D1.Payload[k])
				}
				produceItem(PACKET_START)
				produceItem(dataPackets.D2.Size)
				for k := range dataPackets.D2.Payload {
					produceItem(dataPackets.D2.Payload[k])
				}
				produceItem(PACKET_START)
				produceItem(dataPackets.D3.Size)
				for k := range dataPackets.D3.Payload {
					produceItem(dataPackets.D3.Payload[k])
				}
		}
	}
	fmt.Printf("Producer with tag = %d finished!\n", tag)
	//fmt.Printf("Buffer = %v\n", inputb)
	return
}

// produceItem ----------------------------------------------------------------
// sends a byte to downstream
// ----------------------------------------------------------------------------
func produceItem(item byte) {
	spaces.Wait()
	inputb[ihead] = item
	ihead = (ihead + 1) % BUFFER_SIZE
	items.Signal()
}


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

func processPayload() {}
/*
// consumer--------------------------------------------------------------------
func subchannel(name string)  {
	defer wg.Done()

	wg.Add(1)
	var item int
	for {
		items.Wait()
		cmutex.Enter()
		item = buffer[tail]
		readUARTRegister(name, tail, item)
		tail = (tail + 1) % BUFFER_SIZE
		cmutex.Leave()
		spaces.Signal()
	}
	fmt.Printf("Consumer finished!\n")
}

func consumeSubchannelItem(name string, index int, item int) {
	//fmt.Printf("%s%02d -> %d\n", name, index, item)
}
*/