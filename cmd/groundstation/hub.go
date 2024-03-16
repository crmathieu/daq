package main
import (
	"daq/packages/streamer"
	"github.com/gorilla/websocket"
	"time"
	"sync"
	"math/rand"
	"fmt"
	"errors"
)

// hub maintains the set of active clients
// and broadcasts messages to the clients.

const (
	REG_CHANNELS_SIZE = 64
)

type CLIENT struct {
	Type			uint8				// either "Standard client" (0) or "relay client" (1)
	ClientToken 	string
	Cursor 			*streamer.QueueCursor 	// queue attached to this channel
	Socket         	*websocket.Conn 	// The websocket connection with client
	Valid			bool
	WriteErr		bool
	ReadErr			bool
	Finished     	chan bool			// to close things gracefully
	Ready           chan bool      		// indicate channel is ready for streaming
}

// hub maintains the set of active clients
// and broadcasts messages to the clients.

type Hub struct {
	l 				  *sync.RWMutex           		// handles concurrent access to hubMap
	hubMap            map[string]*CLIENT  			// Registered clients to this launch.
//	launchMap 		  map[string]string         	// clientid -> stream-token
	register   		  chan *CLIENT					// chan *LaunchChannel       // requests from launches to register.
	unregister 		  chan *CLIENT 					// chan *LaunchChannel       // requests from launches to Unregister.
}

func NewHub() *Hub {
	return &Hub{
		l:          &sync.RWMutex{},
		register:   make(chan *CLIENT, REG_CHANNELS_SIZE),
		unregister: make(chan *CLIENT, REG_CHANNELS_SIZE),
		hubMap:     make(map[string]*CLIENT),
//		launchMap:  make(map[string]string),
	}
}

// GetTelemetryClient ---------------------------------------------------------
// determines if a clientToken has already been register in the hub. If
// this is the case, returns the client associated to it
// ---------------------------------------------------------------------------- 
func (h *Hub) GetTelemetryClient(clientToken string) (*CLIENT, error) {
	h.l.RLock()
	client, ok := h.hubMap[clientToken]
	h.l.RUnlock()

	if ok {
		return client, nil
	}
	return nil, errors.New("Hub: Creator Token (" + clientToken + ") has not registered yet!")
}

// FetchToken -----------------------------------------------------------------
// determines if a launch is currently streaming. If this is the case,
// returns the stream-token corresponding to this launch
// ----------------------------------------------------------------------------
/*func (h *Hub) FetchToken(launchid string) (string, error) {
	var token string
	var ok bool
	if token, ok = h.launchMap[launchid]; ok == false {
		return "", errors.New("Launch "+launchid+" is not streaming...")
	}
	return token, nil
}*/


// RandStringBytes ------------------------------------------------------------
// creates a random string
// ----------------------------------------------------------------------------
func RandStringBytes(length int) string {
	const source = "1234567890"
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}

// SetGSInstanceID ------------------------------------------------------------
// sets/creates the goundstation instance ID
// ----------------------------------------------------------------------------
func SetGSInstanceID() {
	gsid = RandStringBytes(3)
}

func CreateLaunchSessionToken() string {
	return RandStringBytes(8)
}

// GetGSInstanceID ------------------------------------------------------------
// gets the goundstation instance ID
// ----------------------------------------------------------------------------
func GetGSInstanceID() string {
	return gsid
}

// AcceptClient ---------------------------------------------------------------
// wait for events (either register or unregister) to add / remove 
// client to the groundstation hub
// ----------------------------------------------------------------------------
func (h *Hub) AcceptClient() {

	for {
		// keep looping waiting on h.register -or- h.unregister channels
		select {

		case client := <-h.register:
			fmt.Println("registering new DAQ connection")
			// a new creator wants to register
			h.l.RLock()
			if _, ok := h.hubMap[client.ClientToken]; ok {
				// this streamer already exists
				fmt.Println("Client: " + client.ClientToken + ": Duplicate connection")
				fmt.Println("Replacing old connection with new...")
				h.hubMap[client.ClientToken].Finished <- true
				delete(h.hubMap, client.ClientToken)
//				delete(h.launchMap, client.clientid)

			}
//			streams, _ := client.Conn.Streams()
//			client.iQue.WriteHeader(streams)
			h.hubMap[client.ClientToken] = client
			//h.launchMap[streamer.Launchid] = streamer.LaunchToken
			client.Ready<- true
			fmt.Println("END -> registering new DAQ connection")

			h.l.RUnlock()


		case client := <-h.unregister:
			fmt.Println("UNregistering DAQ connection")
			h.l.RLock()
			if _, ok := h.hubMap[client.ClientToken]; ok {

				// registering the stream end in temy and removing cache
				//err = StopStream(streamer)
				//if err != nil {
				//	fmt.Println("Unregister: " + err.Error())
				//}

				//if streamer.ReadErr {
				//	h.readErr++
				//}
				//if streamer.WriteErr {
				//	h.writeErr++
				//}
				//h.totalCloseConn++
				delete(h.hubMap, client.ClientToken)
				//delete(h.creatorMap, streamer.Creatorid)
				client.Finished <- true
			}
			fmt.Println("END -> UNregistering DAQ connection")
			h.l.RUnlock()

		}
	}
}

