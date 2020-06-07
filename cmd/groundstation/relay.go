package main
import (
	"fmt"
	"github.com/crmathieu/daq/packages/data"
	"net/url"
	"github.com/gorilla/websocket"
	"unsafe"
)

// RelayListener -------------------------------------------------------------
// relays downlink data with packets read from a groundstation websocket 
// connection and write them in relay streaming queue.
// ----------------------------------------------------------------------------
func (daq *Daq) RelayListener() {

	u := url.URL{Scheme: "ws", Host: daq.RelayFrom, Path: "/wr/"+CreateLaunchSessionToken()}
	fmt.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()
	for {
		_, packet, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		dp := (*[data.PACKET_GRP]data.DataPoint)(unsafe.Pointer(&packet[0]))
		daq.sQue.WriteGrpPacket(dp)
	}
}

func (daq *Daq) RelayListenerByDataPoint() {

	u := url.URL{Scheme: "ws", Host: daq.RelayFrom, Path: "/wr/"+CreateLaunchSessionToken()}
	fmt.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
		return
	}
	defer c.Close()
	var dpg = make([]data.DataPoint, data.PACKET_GRP)
	var npk = 0
	for {
		_, packet, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		dpg[npk] = *(*data.DataPoint)(unsafe.Pointer(&packet[0]))
		if npk + 1 >= data.PACKET_GRP {
			dp := (*[data.PACKET_GRP]data.DataPoint)(unsafe.Pointer(&dpg[0]))
			daq.sQue.WriteGrpPacket(dp)
			npk = -1
		}
		npk++
	}
}
