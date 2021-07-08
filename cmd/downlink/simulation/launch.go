package main

import (
	"fmt"
	"net"
	"os"

	"github.com/crmathieu/daq/packages/data"
)

func deg2rad(degree float64) float64 {
	return M_PI * degree / 180
}

func rad2deg(rad float64) float64 {
	return 180 * rad / M_PI
}

const REALTIME_SIM = true
const CALCULATED_SIM = false

var simulation = CALCULATED_SIM

//var simulation = REALTIME_SIM

func main() {

	var vehicle = NewVehicle()
	if events = vehicle.InitGuidance("profile.yml"); events == nil {
		os.Exit(-1)
	}
	if simulation == REALTIME_SIM {
		//	ctx, cancel := context.WithTimeout(context.Background(), 10000 * time.Millisecond)
		//ctx, cancel := context.WithCancel(context.Background())
		conn, err := net.Dial("tcp", data.DOWNLINK_SERVER+":"+data.DOWNLINK_PORT)
		if err != nil {
			fmt.Println(err)
			return
		}
		go vehicle.launch(simulation)
		vehicle.StreamData(conn)

	} else {
		vehicle.launch(simulation)
	}
}
