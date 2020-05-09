package main
import (
	"github.com/gorilla/websocket"
	"github.com/crmathieu/daq/data"
	"github.com/crmathieu/daq/utils"
	"encoding/json"
	"io/ioutil"

	"time"
	"sync"
	"math/rand"
	"fmt"
	"errors"
	"net/http"
)

type CLIENT struct {
	Key          string          		// corresponds to id_ip_port
	Conn         *websocket.Conn 		// The websocket connection.
	SendTo       chan data.SENSgeneric  // Buffered channel of outbound messages.
	LaunchToken string          		// what identifies this launch
	Valid        bool            		// indicates if we can use this client
	ReadErr      bool
	WriteErr     bool
}

// hub maintains the set of active clients
// and broadcasts messages to the clients.

type Hub struct {
	l                 *sync.Mutex           // handles concurrent access to hubMap
	hubMap            map[string][]*CLIENT  // Registered clients.
	writeErr          int                   // number of write errors encounter
	readErr           int
//	totalStickersSent int
	totalOpenConn     int
	totalCloseConn    int
//	totalAmount       float64
	start             time.Time
	alive             bool
}

//var	gHub   *Hub
//var	gsid  string

func NewHub() *Hub {
	return &Hub{
		l:                 &sync.Mutex{},
		hubMap:            make(map[string][]*CLIENT),
		writeErr:          0,
		readErr:           0,
//		totalStickersSent: 0,
		totalOpenConn:     0,
		totalCloseConn:    0,
//		totalAmount:       0.0,
		start:             time.Now(),
		alive:             true,
	}
}

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

// GetGSInstanceID ------------------------------------------------------------
// gets the goundstation instance ID
// ----------------------------------------------------------------------------
func GetGSInstanceID() string {
	return gsid
}

// Register--------------------------------------------------------------------
func (h *Hub) Register(gsc *CLIENT) {
	gsid := GetGSInstanceID()
	h.l.Lock()
	if _, ok := h.hubMap[gsc.LaunchToken]; ok {
		// this streamer already exists
		ps := h.hubMap[gsc.LaunchToken]
		found := false
		i := 0
		for ; i < len(ps); i++ {
			if ps[i].Key == gsc.Key {
				// duplicate connection
				fmt.Printf("gsInstanceID=%s - stream-token(%s): Hub registration - Duplicate connection for key = %s\n", gsid, gsc.LaunchToken, gsc.Key)
				fmt.Println("Replacing old connection with new...")
				if ps[i].Valid {
					ps[i].Valid = false
					close(ps[i].SendTo)
					ps[i].SendTo = nil
					ps[i].Conn.Close()
					h.totalCloseConn++
				}
				found = true
				break
			}
		}
		if !found {
			h.hubMap[gsc.LaunchToken] = append(h.hubMap[gsc.LaunchToken], gsc)
		} else {
			h.hubMap[gsc.LaunchToken][i] = gsc
		}

	} else {
		// registering the stream start in redis and temy
		err := StartStream(gsc)
		if err != nil {
			fmt.Printf("egoID=%s - stream-token(%s): Hub registration - %s\n", gsid, gsc.LaunchToken, err.Error())
		}

		// new user
		h.hubMap[gsc.LaunchToken] = []*CLIENT{}
		h.hubMap[gsc.LaunchToken] = append(h.hubMap[gsc.LaunchToken], gsc)
	}

	h.totalOpenConn++
	fmt.Printf("stream-token(%s): Hub-registration - Connection successfully registered for key = %s\n", gsc.LaunchToken, gsc.Key)
	h.l.Unlock()
}

// Unregister -----------------------------------------------------------------
func (h *Hub) Unregister(gsc *CLIENT) bool {
	fmt.Printf("Hub Unregistration - Attempting to unregister stream-token %s\n", gsc.LaunchToken)
	h.l.Lock()
	defer h.l.Unlock()

	if _, ok := h.hubMap[gsc.LaunchToken]; ok {
		ps := h.hubMap[gsc.LaunchToken]

		// registering the stream end in temy and removing cache
		err := StopStream(gsc)
		if err != nil {
			fmt.Println("Hub Unregistration - StopStream: " + err.Error())
		}

		i := 0
		found := false
		for ; i < len(ps); i++ {
			if ps[i].Key == gsc.Key {
				// nailed it!
				found = true
				break
			}
		}
		if found {
			h.totalCloseConn++

			if gsc.WriteErr {
				h.writeErr++
			}
			if gsc.ReadErr {
				h.readErr++
			}

			// remove it from array
			fmt.Printf("stream-token(%s): Removing connection for key = %s...\n", gsc.LaunchToken, gsc.Key)
			h.hubMap[gsc.LaunchToken] = append(ps[:i], ps[i+1:]...)
			if gsc.Valid {
				gsc.Valid = false
				close(gsc.SendTo)
				gsc.SendTo = nil
				gsc.Conn.Close()
			}
			if len(h.hubMap[gsc.LaunchToken]) == 0 {
				// remove token from map too
				fmt.Printf("stream-token(%s): No more open connection. Removing token entry from map...\n", gsc.LaunchToken)
				delete(h.hubMap, gsc.LaunchToken)
			}
			fmt.Printf("stream-token(%s): Hub Unregistration - Connection successfully unregistered for key = %s\n", gsc.LaunchToken, gsc.Key)
			return true
		}
	}
	fmt.Printf("stream-token(%s): Hub Unregistration - Could not find connection for key = %s\n", gsc.LaunchToken, gsc.Key)
	return false
}

// getGroundStationClient -----------------------------------------------------
func (h *Hub) getGroundStationClient(launchToken string) ([]*CLIENT, error) {
	h.l.Lock()
	defer h.l.Unlock()

	if connections, ok := h.hubMap[launchToken]; ok {
		return connections, nil
	}
	return nil, errors.New(fmt.Sprintf("stream-token(%s): Hub getGroundStationClient - The token '%s' has not been registered yet.\n", launchToken, launchToken))
}

// StartStream------------------------------------------------------------------
// inserts into redis a key of creator id and value of stream token
func StartStream(client *CLIENT) error {

	launchID := utils.Decode(client.LaunchToken)
	fmt.Printf("StartStream-user(%s): Starting stream session in Redis\n", client.LaunchToken)

	err := Rclient.Set("daq:"+client.LaunchToken, launchID, 0).Err()
	if err != nil {
		fmt.Printf("StartStream-user(%s): Issue setting new Redis session data - %s\n", client.LaunchToken, err.Error())
		return err
	}

	//uts := time.Now().Unix()
	fmt.Printf("StartStream-user(%s): Session Started\n", client.LaunchToken)
	return nil //temy.NotifStarted(launchID, uts)
}

// StopStream-------------------------------------------------------------------
// removes from redis a key of creator id and value of stream token
func StopStream(client *CLIENT) error {

	//launchID := utils.Decode(client.LaunchToken)
	fmt.Printf("StopStream-user(%s) - Removing stream session from Redis\n", client.LaunchToken)
	err := Rclient.Del("daq:" + client.LaunchToken).Err()
	if err != nil {
		fmt.Printf("StopStream-user(%s) - Issue deleting existing Redis session data - %s\n", client.LaunchToken, err.Error())
		return err
	}

	//uts := time.Now().Unix()
	fmt.Printf("StopStream-user(%s) - Session Stopped\n", client.LaunchToken)
	return nil // temy.NotifEnded(launchID, uts)
}

// sendSticker reads what Pin sends it--------------------------------------------
func sendSticker(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gsid := GetGSInstanceID()

	plBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("SendSticker(egoID=%s): Error extracting payload from request body: %s\n", gsid, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pl data.SENSgeneric
	err = json.Unmarshal(plBytes, &pl)
	if err != nil {
		fmt.Printf("SendSticker(egoID=%s): Error unmarshalling payload - %s\nPAYLOAD=%s\n", gsid, err.Error(), string(plBytes))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	/*
	// check whether user is logged-in or not
	_, err = GetSessionCoredata(pl.UserID, &w, r)
	if err != nil {
		fmt.Printf("Medusa-user(%s): access denied! - %s\n", pl.UserID, err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	*/
	// publish this request to all instances that subscribed to the same channel
	fmt.Printf("SendSticker(egoID=%s) - publishing payload: %s\n", GetGSInstanceID(), string(plBytes))
	err = Publish(string(plBytes))
	if err != nil {
		//panic(err)
		fmt.Printf("SendSticker(egoID=%s): error publishing payload: %s\n", GetGSInstanceID(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
