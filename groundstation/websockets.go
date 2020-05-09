package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"github.com/crmathieu/daq/data"
	//"github.com/crmathieu/daq/utils"
	"github.com/go-redis/redis"
	//"os"
	"encoding/json"
	"time"
	"html/template"	
)

var (
//	addr     = flag.String("addr", ":8088", "http service address")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	Rclient *redis.Client
/*	streamRedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("CORE_CACHE_ADDR"),
		Password: os.Getenv("CORE_CACHE_PWD"),
		DB:       0,
	})*/

	gHub   *Hub
	gsid  string
	gsEnv string
)

const (
	SENDTO_CHANNEL_READ_DELAY = 13
	SENDTO_CHANNEL_SIZE       = 30
	DEFAULT_PLATFORM          = "stage"
)

// establishes a websocket connection with a client----------------------------
func newConn(w http.ResponseWriter, r *http.Request) {

	launchTok := r.URL.Path[len("/ws/"):]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		fmt.Println("Upgrader failed :(")
		return
	}

	// build an OBS client and register it to the connection hub
	streamer := CLIENT{
		Key:          setClientKey(ws, launchTok),
		Conn:         ws,
		SendTo:       make(chan data.SENSgeneric, SENDTO_CHANNEL_SIZE),
		LaunchToken: launchTok,
		Valid:        true,
		WriteErr:     false,
		ReadErr:      false,
	}

	// register CLIENT pointer with hub
	gHub.Register(&streamer)

	// now starts a go routine that will handle the connection writes
	// and a go routine that will detect when a connection drops
	go stickerListener(launchTok, &streamer)
	go readForDisconnect(&streamer)
}

//  closeConn------------------------------------------------------------------
// closes the websocket connection associated with a given userid.
// this will work well with a single egomonster instance, but for
// multi-instances and assuming we have a way to target each instance
// individually, we should send the same request to close a websocket
// connection to all egomonster instances. This way, we are sure that
// there will be at least one instance that will succeed.
// ----------------------------------------------------------------------------
func closeConn(w http.ResponseWriter, r *http.Request) {
	launchToken := r.URL.Path[len("/ws/"):]
	streamer, err := gHub.getGroundStationClient(launchToken)
	if err != nil {
		// this happens when a userid is not registered with this instance.
		// the same userid could be registered with another instance. The problem
		// is that we don't know which one. THAT'S A MAJOR FLAW FOR CLOSECONN
		fmt.Printf("CloseConn - %s\n", err.Error())
		return
	}
	gHub.Unregister(streamer[0])
}

// setClientKey----------------------------------------------------------------
func setClientKey(ws *websocket.Conn, id string) string {
	return id + "_" + ws.RemoteAddr().String()
}

// stickerListener-------------------------------------------------------------
func stickerListener(token string, streamer *CLIENT) {

	defer func() {
		// recover from panic caused by writing to a closed channel
		if r := recover(); r != nil {
			fmt.Printf("\n*** Recovering from panic\n%v\n***\n", r)
		}
	}()

	// wait for an event to happen (either SendTo or Finished channel)
	loop := true
	for loop && streamer.Valid {
		select {
		case pl, ok := <-streamer.SendTo:
			if ok {
				fmt.Println("Receiving payload from sendTo channel...")
				var payloadInBytes []byte
				payloadInBytes, err := json.Marshal(pl)
				if err == nil {
					if streamer.Valid {
						err = streamer.Conn.WriteMessage(websocket.TextMessage, payloadInBytes)
						if err == nil {
							// wait a bit before reading the next payload
							//gHub.totalStickersSent++
							//gHub.totalAmount += pl.Total
							time.Sleep(SENDTO_CHANNEL_READ_DELAY * time.Second)
							continue
						}
						fmt.Printf("stream-token(%s): Writing Error - %s\n", token, err.Error())
						streamer.WriteErr = true
						//Ghub.Unregister(streamer)
					}
				} else {
					fmt.Printf("stream-token(%s): Payload Unmarshalling Error - %s\n", token, err.Error())
					continue
				}
			} else {
				streamer.SendTo = nil
				fmt.Println("Channel closed...")
			}
			loop = false
		}
	}
	if streamer.Valid {
		fmt.Println("Unregistering...")
		gHub.Unregister(streamer)
	}
	fmt.Printf("stream-token(%s): Closing CLIENT\n", streamer.LaunchToken)
}

// readForDisconnect----------------------------------------------------------
func readForDisconnect(streamer *CLIENT) {
	// read connection (we ignore what is being received) until an error happens
	ever := true
	for ever {
		_, _, err := streamer.Conn.ReadMessage()
		if err != nil {
			// The connection dropped
			fmt.Printf("stream-key(%s): Connection dropped  - %s\n", streamer.Key, err.Error())
			streamer.ReadErr = true
			ever = false
			gHub.Unregister(streamer)
		}
	}
}

// called on /stream endpoint. returns a template filled with creatorID--------
func serveHome(page http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
/*		if r.URL.Path[len("/stream/"):] == "" {
			http.Error(w, "Missing Stream Token", http.StatusMethodNotAllowed)
			return
		}*/
		page(w, r)
	})
}

// homeTest--------------------------------------------------------------------
func homeTest(w http.ResponseWriter, r *http.Request) {

	homeTempl, err := template.ParseFiles("./assets/html/index.html")
	if err != nil {
		fmt.Println("Error parsing index template:")
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	launchtoken := r.URL.Path[len("/stream/"):]

	// prepare template data ...
	var v = struct {
		SocketType   string
		Host         string
		LaunchToken string
		SoundOn      string
		NotifVol     string
		Error        string
	}{
		"ws",
		r.Host,
		launchtoken,
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
	launchtoken := "123" //r.URL.Path[len("/stream/"):]

	// decode the token so that we can query the DB
	//launchID := utils.Decode(launchtoken)

	// Unmarshal settings into this struct
	notificationSettings := struct {
		SoundOn  string `json:"soundOn"`
		NotifVol string `json:"notifVol"`
	}{}

	/*
	// using the creatorID - query MySQL DB for settings related to sounds et al
	err = database.Db.QueryRow("SELECT sound_on, volume FROM NotificationSettings WHERE user_id = ?", launchID).Scan(&notificationSettings.SoundOn, &notificationSettings.NotifVol)
	if err != nil {
		fmt.Printf("home - Error selecting notifications settings - user_id = %s : %s\n", launchID, err.Error())
		// get the string of the error so they have some info
		s := err.Error()
		// take note that error happened and show user so they can retry
		notifError = "There was an issue fetching your Alert settings! Please try refreshing or updating your settings! " + s
	}
*/

	// prepare template data ...
	var v = struct {
		SocketType   string
		Host         string
		LaunchToken string
		SoundOn      string
		NotifVol     string
		Error        string
	}{
		"ws",
		r.Host,
		launchtoken,
		notificationSettings.SoundOn,
		notificationSettings.NotifVol,
		notifError,
	}
	homeTempl.Execute(w, &v)
}
