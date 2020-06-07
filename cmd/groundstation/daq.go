package main

import (
	"io" 
	"net"
	"github.com/crmathieu/daq/packages/streamer"
	"github.com/crmathieu/daq/packages/data"
	"time"
	"fmt"
	"unsafe"
)

type Daq struct {
	Addr          	string
	Port			string
	RelayFrom		string
	ConnListener    func()
	VehicleProfile	string
	sQue 			*streamer.Queue
}

var DAQ *Daq 

// NewDaq ---------------------------------------------------------------------
// creates a new Data Acquisition object
// ----------------------------------------------------------------------------
func NewDaq(server, port, relayFrom string) *Daq {
	var Daq = Daq{
		Addr: server,
		Port: port,
		RelayFrom: relayFrom,
		// create a streamer for this downlink
		sQue: streamer.NewQueue(),
	}
	if relayFrom != "" {
		Daq.ConnListener = Daq.RelayListener
	} else {
		Daq.ConnListener = Daq.ListenAndServe
	}
	return &Daq
}


// ListenAndServe -------------------------------------------------------------
// listen to downlink connection and read packets
// ----------------------------------------------------------------------------
func (daq *Daq) ListenAndServe() {
	// established downlink with launch vehicle
	fmt.Println("listening on:", daq.Addr+":"+daq.Port)
	srv, err := net.Listen("tcp", daq.Addr+":"+daq.Port) 
	if err != nil {
		fmt.Println(err)
		return
	}

	// set up telemetry hub
	go GrndStationHUB.AcceptClient()

	fmt.Println("Ground Station Listening for downlink on", daq.Addr+":"+daq.Port) 

	for {
		conn, err := srv.Accept()
		fmt.Println("got connection")
		if err != nil {
			fmt.Println(err)
			return
		}
		go daq.ReadDownlinkPackets(conn)
	}
}


// ReadDownlinkPackets --------------------------------------------------------
// reads one or more data packet from connection
// ----------------------------------------------------------------------------
func (daq *Daq) ReadDownlinkPackets(c net.Conn) {
	const MAXERROR = 5
	start := time.Now()
    tbuf := make([]byte, 81920)
    totalBytes := 0
	consecutiveErr := 0
	
    for {
		n, err := c.Read(tbuf)
        if err != nil {
            if err != io.EOF {
                fmt.Printf("Read error: %s", err)
			}
			consecutiveErr++
			if consecutiveErr > MAXERROR {
				break
			}
            continue 
        }
		consecutiveErr = 0

		/*if n > 81920 {
			log.Printf("BUFFER OVERFLOW %d bytes read in %s", totalBytes, time.Now().Sub(start))
			c.Close()
		}*/

		daq.demuxDataPoints(&tbuf, n)
		totalBytes += n
    }
    fmt.Printf("%d bytes read in %s", totalBytes, time.Now().Sub(start))
    c.Close()
}

// demuxDataPoints ------------------------------------------------------------
// from the number of datapoints value found in the packet header, reads each
// datapoint and save it in DAQ streamer
// ----------------------------------------------------------------------------  
func (daq *Daq) demuxDataPoints(pk *[]byte, size int) {

	numberDP := *(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET]))

	// calculate checksum on datapoints 
	if data.CRC32(0, (*pk)[data.PACKET_PAYLOAD_OFFSET:], int(numberDP) * data.DATAPOINT_SIZE) == *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) {
//		for k:=byte(0); k<numberDP; k++ {
			dp := (*[data.PACKET_GRP]data.DataPoint)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET])) //+k*data.DATAPOINT_SIZE]))
			daq.sQue.WriteGrpPacket(dp)
	//		daq.viewPacket(dp)
//		}
	} else {
		fmt.Println("CRC error encountered...")
	}

}
/*
func (daq *Daq) demuxDataPoints(pk *[]byte, size int) {

	numberDP := *(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET]))

	// calculate checksum on datapoints 
	if data.CRC32(0, (*pk)[data.PACKET_PAYLOAD_OFFSET:], int(numberDP) * data.DATAPOINT_SIZE) == *(*uint32)(unsafe.Pointer(&(*pk)[data.PACKET_CRC_OFFSET])) {
		for k:=byte(0); k<numberDP; k++ {
			dp := (*data.DataPoint)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET+k*data.DATAPOINT_SIZE]))
			daq.sQue.WritePacket(*dp)
	//		daq.viewPacket(dp)
		}
	} else {
		fmt.Println("CRC error encountered...")
	}

}*/

// viewPacket -----------------------------------------------------------------
// reads each datapoint as "Generic dp" and then cast to its appropriate type 
// based on the datapoint Id
// ----------------------------------------------------------------------------
func viewPacket(dp *data.DataPoint) {
		switch(dp.Id) {
		case data.IDVELOCITY: 	v := (*data.SENSvelocity)(unsafe.Pointer(dp))
								fmt.Println("Vel:",v.Velocity,"m/s, Acc:", v.Acceleration)
		case data.IDPOSITION:	v := (*data.SENSposition)(unsafe.Pointer(dp))
								fmt.Println("Alt:", v.Altitude,"m, Range:", v.Range,"m")
		case data.IDTHRUST:		v := (*data.SENSthrust)(unsafe.Pointer(dp))
								fmt.Println("Thrust:", v.Thrust/1000,"kN, Stage:", v.Stage)
//		case data.IDTILTANGLE:	v := (*data.SENStiltAngle)(unsafe.Pointer(dp))
//								fmt.Println("Gamma:", v.Angle,"deg")
		case data.IDMASSPROPELLANT:	v := (*data.SENSpropellantMass)(unsafe.Pointer(dp))
								//fmt.Println("Mass:", v.Mass,"kg, Mass Flow:", v.Mflow, "kg/s, Mass Ejected:", v.Mejected)
								fmt.Println("Mass:", v.Mass,"kg, Mass Flow:", v.Mflow, "kg/s, Stage:", v.Stage)

		}

}
