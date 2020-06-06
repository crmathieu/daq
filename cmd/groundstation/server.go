package main

import (
	"github.com/rs/cors"
	"github.com/crmathieu/daq/packages/data"
	"net/http"
	"fmt"
	"os"
)

const PATH_2_ROOT = "."
// main call as:
// ./groundstation for main station or for a relay
// ./groundstation -r localhost:port -p webport
func main() {

	var grdserver, webport string
	relay := false

	grdserver = data.DOWNLINK_SERVER
	webport = data.DOWNLINK_WEBPORT //"1969"
	args := os.Args[1:]

	for i := range(args) {
		switch (args[i]) {
		case "-r":
			grdserver = args[i+1]
			relay = true
			break
		case "-p":
			webport = args[i+1]
			break;
		default:
			if []byte(args[i])[0] == '-' {
				fmt.Println(args[i],": unknown command switch")
			}
		}
	}


//	if InitGroundStation() {
		LaunchHUB = NewHub()

		// set up a http server to manage data streaming to clients
		mux := http.NewServeMux()
		fs := http.FileServer(http.Dir(PATH_2_ROOT + "/assets"))
		mux.Handle("/", fs)
		mux.HandleFunc("/stream/", serveHome(home))
		mux.HandleFunc("/ws/", NewGroundStationClient)
		mux.HandleFunc("/wr/", NewGroundStationRelay)
		mux.HandleFunc("/wc/", CloseGroundStationClient)

		server := &http.Server{
//			Addr: "0.0.0.0:" + webport, //1969",
			Addr: ":" + webport, //1969",

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
		if !relay {
			fmt.Printf("\n...GroundStation web server is accepting connection on port: %s\n", webport)
		} else {
			fmt.Println(grdserver)
			fmt.Printf("\n...GroundStation RELAY web server is accepting connection on port: %s\n", webport)
		}
		go server.ListenAndServe()	

		DAQ = NewDaq(grdserver, relay)
		//		DAQ = NewDaq(grdserver, relay)
		DAQ.ConnListener()
//	}
}
