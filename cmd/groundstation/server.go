package main

import (
	"github.com/rs/cors"
	"daq/packages/data"
	"net/http"
	"fmt"
	"os"
)

const PATH_2_ROOT = "."

// main call as:
// ----------------------------------------------------------------------------
// for ground station 
// ./groundstation 
// default values can be overwritten as well
// ./groundstation -s localhost -p port -wp webport
// ----------------------------------------------------------------------------
// for groundstation relay
// ./groundstation -r localhost -wp webport
// ----------------------------------------------------------------------------
func main() {

	relayFrom := ""

	// set up server as groundstation using default values
	grdserver 	:= data.DOWNLINK_SERVER 	// localhost
	port 		:= data.DOWNLINK_PORT 		// 2000
	webport 	:= data.DOWNLINK_WEBPORT	// 1969
	args 		:= os.Args[1:]

	// to start a grounstation with other server values than the default, type:
	// ./groundstation -s <server> -p <port> -wp <webport>
	// to start a groundstation relay, type on the command line
	// ./groundstation -r <server> -wp <webport>
	// <webport> must be different than the groundstation webport (1969)

	for i := range(args) {
		switch (args[i]) {
		case "-r":
			relayFrom = grdserver + ":" + webport
			grdserver = args[i+1]
			break
		case "-s":
			grdserver = args[i+1]
			break
		case "-p":
			port = args[i+1]
			break;
		case "-wp":
			webport = args[i+1]
			break;
		default:
			if []byte(args[i])[0] == '-' {
				fmt.Println(args[i],": unknown command switch")
			}
		}
	}
	GrndStationHUB = NewHub()

	// set up a http server to manage data streaming to clients
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir(PATH_2_ROOT + "/assets"))
	mux.Handle("/", fs)
	mux.HandleFunc("/stream/", 	serveHome(home))
	mux.HandleFunc("/test/", 	serveHome(homeTest))
	mux.HandleFunc("/ws/", 		NewGroundStationClient)
	mux.HandleFunc("/wr/", 		NewGroundStationRelay)
	mux.HandleFunc("/wc/", 		CloseGroundStationClient)

	server := &http.Server{
		Addr: grdserver + ":" + webport,
		Handler: cors.New(cors.Options{
			AllowedHeaders: []string{
				"Authorization",
				"Origin",
				"X-Requested-With",
				"Accept",
				"Content-Type"},
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			Debug:            false,
		}).Handler(mux),
	}
	if relayFrom == "" {
		fmt.Printf("\n...GroundStation web server is accepting connection on port: %s\n", webport)
	} else {
		fmt.Printf("\n...GroundStation RELAY web server is accepting connection on port: %s\n", webport)
	}
	go server.ListenAndServe()	

	DAQ = NewDaq(grdserver, port, relayFrom)
	DAQ.ConnListener()
}
