package main

import (
	"github.com/rs/cors"
	"net/http"
	"fmt"
)

//var iQue *queue.Queue

const PATH_2_ROOT = "."
func main() {

	if InitConfig() {
		LaunchHUB = NewHub()

		// create an http server to serve client assets
		mux := http.NewServeMux()
		fs := http.FileServer(http.Dir(PATH_2_ROOT + "/assets"))
		mux.Handle("/", fs)
		mux.HandleFunc("/stream/", serveHome(home))
		mux.HandleFunc("/ws/", NewLaunchClient)
		mux.HandleFunc("/wc/", closeConn)

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

			//Handler: cors.AllowAll().Handler(mux),
		}
		fmt.Printf("\n...GroundStation %s is running on port: 1969\n", GetGSInstanceID())
		go server.ListenAndServe()	

		DACQ = NewDaq()
		DACQ.ListenAndServer()

		//log.Fatal(server.ListenAndServe())

		//time.Sleep(10 * time.Second)
	}
}
