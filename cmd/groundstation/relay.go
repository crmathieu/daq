package main
import (
//	"net"
	"fmt"
//	"time"
//	"io"
	"github.com/crmathieu/daq/packages/data"
	"net/url"
	"github.com/gorilla/websocket"
	"unsafe"
 	//"encoding/binary"
)

// RelayListener -------------------------------------------------------------
// relays downlink data from packets read and writing them in streaming queue 
// ----------------------------------------------------------------------------
func (daq *Daq) RelayListener() {

//	u := url.URL{Scheme: "ws", Host: data.DOWNLINK_SERVER, Path: "/ws/"+CreateLaunchSessionToken()}
	u := url.URL{Scheme: "ws", Host: "localhost:1969", Path: "/wr/"+CreateLaunchSessionToken()}
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
