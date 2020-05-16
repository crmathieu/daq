package main

import (
	"github.com/crmathieu/daq/data"
//	"github.com/crmathieu/daq/downlinksim/instruments"
//	"unsafe"
	"fmt"
	"net"
//	"context"
//	"time"
)

const (
	F9_S1_Thrust 		 = 7686000 		// 845000 * 9	// newtons
	F9_S1_dryMass 		 = 25600		// kg
	F9_S1_PropellantMass = 395700		// kg
	F9_S1_burntime 		 = 162			// seconds
	F9_S1_ISP 			 = 311			// seconds
	F9_S1_DIAMETER 		 = 3.66			// meters

	F9_S2_Thrust 		 = 934000		// Newtons
	F9_S2_dryMass 		 = 3900			// kg
	F9_S2_PropellantMass = 92670		// kg
	F9_S2_burntime 		 = 397			// seconds
	F9_S2_ISP 			 = 348			// seconds
	F9_S2_DIAMETER 		 = 3.66			// meters
)

var vehicle *VEHICLE 

func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
	//ctx, cancel := context.WithCancel(context.Background())
	vehicle = NewVehicle()
    conn, err := net.Dial("tcp", data.DOWNLINK_SERVER)
    if err != nil {
		fmt.Println(err)
		return
	}
	//go vehicle.RunInstrumentsUpdate()
	go vehicle.launch()
	vehicle.StreamData(conn)
}
