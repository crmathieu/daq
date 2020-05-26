package main

import (
	"net"
	"github.com/crmathieu/daq/packages/data"
	"fmt"
)

func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
	//ctx, cancel := context.WithCancel(context.Background())
	var vehicle = NewVehicle()
    conn, err := net.Dial("tcp", data.DOWNLINK_SERVER)
    if err != nil {
		fmt.Println(err)
		return
	}
	//go vehicle.RunInstrumentsUpdate()
	go vehicle.launch()
	vehicle.StreamData(conn)
}
