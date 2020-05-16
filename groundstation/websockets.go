package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"github.com/crmathieu/daq/data"
	"github.com/go-redis/redis"
	"html/template"	
	"unsafe"
	"io"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,	
		WriteBufferSize: 1024,
	}

	Rclient *redis.Client

	LaunchHUB   *Hub
	gsid  string
	gsEnv string
)

// NewLaunchClient ------------------------------------------------------------
// establishes a websocket connection with a client
// ----------------------------------------------------------------------------
func NewLaunchClient(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		fmt.Println("Upgrader failed :(")
		return
	}

	// build an OBS client and register it to the connection hub
	client := CLIENT{
		Cursor:		  DACQ.iQue.Latest(), // obtain cursor of latest position from DACQ
		Socket:       ws,
		ClientToken:  r.URL.Path[len("/ws/"):], //clientToken,
		Valid:        true,
		WriteErr:     false,
		ReadErr:      false,
	}

	// register CLIENT pointer with hub
	LaunchHUB.register <- &client

	// now starts a go routine that will handle the connection writes
	go WriteLaunchTelemetry(&client)
	// and a go routine that will detect when a connection drops
	go DetectClientDisconnection(&client)
}

// closeConn ------------------------------------------------------------------
// closes the websocket connection associated with a given userid.
// this will work well with a single egomonster instance, but for
// multi-instances and assuming we have a way to target each instance
// individually, we should send the same request to close a websocket
// connection to all egomonster instances. This way, we are sure that
// there will be at least one instance that will succeed.
// ----------------------------------------------------------------------------
func closeConn(w http.ResponseWriter, r *http.Request) {

	clientToken := r.URL.Path[len("/ws/"):]
	client, err := LaunchHUB.GetLaunchClient(clientToken)
	if err != nil {
		// this happens when a userid is not registered with this instance.
		// the same userid could be registered with another instance. The problem
		// is that we don't know which one. THAT'S A MAJOR FLAW FOR CLOSECONN
		fmt.Printf("CloseConn - %s\n", err.Error())
		return
	}
	LaunchHUB.unregister <- client
}

// setClientKey ---------------------------------------------------------------
func setClientKey(ws *websocket.Conn, id string) string {
	return id + "_" + ws.RemoteAddr().String()
}


// WriteLaunchTelemetry --------------------------------------------------------
func WriteLaunchTelemetry(client *CLIENT) {
	var dp data.DataPoint
	var err error
	for {
		if dp, err = client.Cursor.ReadPacket(); err != nil {
			fmt.Println(err)
		} else {
			client.Socket.WriteMessage(websocket.BinaryMessage, (*(*[data.DATAPOINT_SIZE]byte)(unsafe.Pointer(&dp)))[:data.DATAPOINT_SIZE])
		} 
	}
}

// DetectClientDisconnection --------------------------------------------------
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
			LaunchHUB.unregister <- client
		}
	}
}

// serverHome -----------------------------------------------------------------
// called on /stream endpoint. returns a template filled with a launchtoken
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
	fmt.Println("sending template")
	token := r.URL.Path[len("/stream/"):]

	// the token can be used to authorize certain clients only
	if authorized(token) {

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
		io.WriteString(w, fmt.Sprintf("%s: Invalid token\n", token))
	}
}

func authorized(token string) bool {
	// place code here to authorize only certain tokens
	if token != "123" {
		return false
	}
	return true
}