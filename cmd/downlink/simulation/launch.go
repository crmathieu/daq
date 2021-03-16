package main

import (
	"net"
	"github.com/crmathieu/daq/packages/data"
	"fmt"
)

func deg2rad(degree float64) float64 {
	return M_PI * degree / 180 
}

func rad2deg(rad float64) float64 {
	return 180 * rad / M_PI 
}


func main() {
//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
	//ctx, cancel := context.WithCancel(context.Background())
	var vehicle = NewVehicle()
    conn, err := net.Dial("tcp", data.DOWNLINK_SERVER+":"+data.DOWNLINK_PORT)
    if err != nil {
		fmt.Println(err)
		return
	}
	//go vehicle.RunInstrumentsUpdate()
	go vehicle.launch()
	vehicle.StreamData(conn)
}
