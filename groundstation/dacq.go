package main

import (
//	"bufio"
/*	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/nareix/joy4/utils/bits/pio"
	"github.com/nareix/joy4/av"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/flv/flvio"*/
	"io" 
	"net"
	//"net/http"
//	"github.com/gorilla/websocket"
	"github.com/crmathieu/daq/packages/queue"
	"github.com/crmathieu/daq/data"
/*	"net/url"
	"strings"*/
	"time"
	"fmt"
	"unsafe"
)

type dacq struct {
	Addr          	string
//	HandlePublish func(*net.Conn)
//	HandlePlay    func(*websocket.Conn)
//	HandleConn    func(*websocket.Conn)	
	VehicleProfile	string
	iQue 			*queue.Queue
}

var DACQ *dacq 

// NewDaq ---------------------------------------------------------------------
// creates a new Data Acquisition object
// ----------------------------------------------------------------------------
func NewDaq() *dacq {
	return &dacq{
		// create a queue for this downlink
		iQue: queue.NewQueue(),
	}
}

// ListenAndServer ------------------------------------------------------------
// listen to downlink connection and read packets
// ----------------------------------------------------------------------------
func (dacq *dacq) ListenAndServer() {
	// established downlink with launch vehicle
	srv, err := net.Listen("tcp", data.DOWNLINK_SERVER)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Ground Station Listening for downlink on", data.DOWNLINK_SERVER)

	for {
		conn, err := srv.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go dacq.ReadDownlinkPackets(conn)
	}
}

// ReadDownlinkPackets --------------------------------------------------------
// reads one or more data packet from connection
// ----------------------------------------------------------------------------
func (dacq *dacq) ReadDownlinkPackets(c net.Conn) {
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

		//log.Println(tbuf[:n])
		dacq.demuxDataPoints(&tbuf, n)
		totalBytes += n
//		if n != 256 {
//			log.Println(n)
//		}
        //log.Println(n)
    }
    fmt.Printf("%d bytes read in %s", totalBytes, time.Now().Sub(start))
    c.Close()
}

// demuxDataPoints ------------------------------------------------------------
// from the number of datapoints value found in the packet header, reads each
// datapoint and save it in DAQ queue
// ----------------------------------------------------------------------------  
func (dacq *dacq) demuxDataPoints(pk *[]byte, size int) {
	numberDP := *(*byte)(unsafe.Pointer(&(*pk)[data.PACKET_NDP_OFFSET]))
	for k:=byte(0); k<numberDP; k++ {
		dp := (*data.DataPoint)(unsafe.Pointer(&(*pk)[data.PACKET_PAYLOAD_OFFSET+k*data.DATAPOINT_SIZE]))
		dacq.iQue.WritePacket(*dp)
//		dacq.viewPacket(dp)
	}
}

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
		case data.IDTILTANGLE:	v := (*data.SENStiltAngle)(unsafe.Pointer(dp))
								fmt.Println("Gamma:", v.Angle,"deg")
		case data.IDMASSPROPELLANT:	v := (*data.SENSpropellantMass)(unsafe.Pointer(dp))
								fmt.Println("Mass:", v.Mass,"kg, Mass Flow:", v.Mflow, "kg/s, Mass Ejected:", v.Mejected)

		}

}
