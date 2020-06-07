package main

import (
	"github.com/gorilla/websocket"
	"github.com/crmathieu/daq/packages/data"
	"github.com/go-redis/redis"
	"html/template"	
	"net/http"
	"fmt"
	"unsafe"
	"io"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,	
		WriteBufferSize: 1024,
	}

	Rclient *redis.Client

	GrndStationHUB   *Hub
	gsid  string
	gsEnv string
)

// NewGroundStationClient -----------------------------------------------------
// establishes a websocket connection with a (browser) client
// ----------------------------------------------------------------------------
func NewGroundStationClient(w http.ResponseWriter, r *http.Request) {
	NewClient(w, r, 0)
}

// NewGroundStationRelay ------------------------------------------------------
// establishes a websocket connection with a Relay
// ----------------------------------------------------------------------------
func NewGroundStationRelay(w http.ResponseWriter, r *http.Request) {
	NewClient(w, r, 1)
}

func NewClient(w http.ResponseWriter, r *http.Request, ctype uint8) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		fmt.Println("Upgrader failed :(")
		return
	}

	// build a client ...
	client := CLIENT{
		Type:		  ctype,
		Cursor:		  DAQ.sQue.Latest(), // obtain cursor of latest position from DAQ
		Socket:       ws,
		ClientToken:  r.URL.Path[len("/wr/"):], //clientToken,
		Valid:        true,
		WriteErr:     false,
		ReadErr:      false,
	}

	// ... and register it to the connection hub
	GrndStationHUB.register <- &client

	// start a go routine that will send telemetry data to client
	go WriteLaunchTelemetry(&client)
	// start a go routine that will detect when a connection drops
	go DetectClientDisconnection(&client)
}

// CloseGroundStationClient ---------------------------------------------------
// closes the websocket connection associated to a given clienttoken.
// ----------------------------------------------------------------------------
func CloseGroundStationClient(w http.ResponseWriter, r *http.Request) {

	clientToken := r.URL.Path[len("/ws/"):]
	client, err := GrndStationHUB.GetTelemetryClient(clientToken)
	if err != nil {
		// this happens when a client is not registered with this instance.
		// the same userid could be registered with another instance. The problem
		// is that we don't know which one. THAT'S A MAJOR FLAW FOR CLOSECONN
		fmt.Printf("CloseConn - %s\n", err.Error())
		return
	}
	GrndStationHUB.unregister <- client
}

// setClientKey ---------------------------------------------------------------
func setClientKey(ws *websocket.Conn, id string) string {
	return id + "_" + ws.RemoteAddr().String()
}


// WriteLaunchTelemetry --------------------------------------------------------
// streams the content of telemetry queue to client
// -----------------------------------------------------------------------------
func WriteLaunchTelemetry(client *CLIENT) {
	var dp [data.PACKET_GRP]data.DataPoint
	var err error
	for {
		if dp, err = client.Cursor.ReadGrpPacket(); err != nil {
			fmt.Println(err)
		} else {
			client.Socket.WriteMessage(websocket.BinaryMessage, (*(*[data.DATAPOINT_SIZE * data.PACKET_GRP]byte)(unsafe.Pointer(&dp[0])))[:data.DATAPOINT_SIZE * data.PACKET_GRP])
		}
	}
}

func WriteLaunchTelemetryByDataPoint(client *CLIENT) {
	var dp [data.PACKET_GRP]data.DataPoint
	var err error
	for {
		if dp, err = client.Cursor.ReadGrpPacket(); err != nil {
			fmt.Println(err)
		} else {
			for k := range dp {
				client.Socket.WriteMessage(websocket.BinaryMessage, (*(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(&dp[k])))[:data.DATAPOINT_SIZE])
			}
//			client.Socket.WriteMessage(websocket.BinaryMessage, (*(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(&dp)))[:data.DATAPOINT_SIZE])
		} 
	}
}

// DetectClientDisconnection --------------------------------------------------
// handles disconnections by unregistering client from hub. The client will 
// make a retry to reconnect and re-register with hub
// ----------------------------------------------------------------------------
func DetectClientDisconnection(client *CLIENT) {
	// read connection (we ignore what is being received) until an error happens
	ever := true
	for ever {
		_, _, err := client.Socket.ReadMessage()
		if err != nil {
			// The connection dropped
			fmt.Printf("stream-key(%s): Connection dropped  - %s\n", client.ClientToken, err.Error())
			client.ReadErr = true
			ever = false
			GrndStationHUB.unregister <- client
		}
	}
}

// serverHome -----------------------------------------------------------------
// called on /stream endpoint. returns a template filled with a clientToken
// ----------------------------------------------------------------------------
func serveHome(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		page(w, r)
	})
}

// homeTest -------------------------------------------------------------------
func homeTest(w http.ResponseWriter, r *http.Request) {

	homeTempl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		fmt.Println("Error parsing index template:")
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	clientToken := r.URL.Path[len("/stream/"):]

	// prepare template data ...
	var v = struct {
		SocketType   string
		Host         string
		ClientToken  string
		SoundOn      string
		NotifVol     string
		Error        string
	}{
		"ws",
		r.Host,
		clientToken,
		"1",
		"50",
		"",
	}
	homeTempl.Execute(w, &v)
}

// home------------------------------------------------------------------------
func home(w http.ResponseWriter, r *http.Request) {

	// notifError is to display if there was an issue getting user settings from MySQL
	var notifError string

	homeTempl, err := template.ParseFiles("./assets/index.html")
	if err != nil {
		fmt.Println("Error parsing index template:")
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	authToken := r.URL.Path[len("/stream/"):]

	// the client token can be used to authorize certain clients only
	if authorized(authToken) {

		// Unmarshal settings into this struct
		notificationSettings := struct {
			SoundOn  string `json:"soundOn"`
			NotifVol string `json:"notifVol"`
		}{}

		// prepare template data ...
		var v = struct {
			SocketType   string
			Host         string
			ClientToken  string
			SoundOn      string
			NotifVol     string
			Error        string
		}{
			"ws",
			r.Host,
			CreateLaunchSessionToken(),
			notificationSettings.SoundOn,
			notificationSettings.NotifVol,
			notifError,
		}
		homeTempl.Execute(w, &v)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, fmt.Sprintf("%s: Invalid authToken\n", authToken))
	}
}

func authorized(token string) bool {
	// place code here to authorize only certain tokens
	if token != "123" {
		return false
	}
	return true
}