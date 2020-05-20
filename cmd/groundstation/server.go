package main

import (
	"github.com/rs/cors"
	"net/http"
	"fmt"
)

const PATH_2_ROOT = "."
func main() {

	if InitGroundStation() {
		LaunchHUB = NewHub()

		// create an http server to serve client assets
		mux := http.NewServeMux()
		fs := http.FileServer(http.Dir(PATH_2_ROOT + "/assets"))
		mux.Handle("/", fs)
		mux.HandleFunc("/stream/", serveHome(home))
		mux.HandleFunc("/ws/", NewGroundStationClient)
		mux.HandleFunc("/wc/", CloseGroundStationClient)

		server := &http.Server{
			Addr: "0.0.0.0:1969",

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
		fmt.Printf("\n...GroundStation %s is running on port: 1969\n", GetGSInstanceID())
		go server.ListenAndServe()	

		DAQ = NewDaq()
		DAQ.ListenAndServer()
	}
}
