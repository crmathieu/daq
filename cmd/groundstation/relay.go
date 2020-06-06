package main
import (
	"fmt"
	"github.com/crmathieu/daq/packages/data"
	"net/url"
	"github.com/gorilla/websocket"
	"unsafe"
)

// RelayListener -------------------------------------------------------------
// relays downlink data from packets read and writing them in streaming queue 
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
		dp := (*data.DataPoint)(unsafe.Pointer(&packet[0]))
		daq.sQue.WritePacket(*dp)
	}
}
