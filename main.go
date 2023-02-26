package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

//Creates a HTTP server that will check for requests and redirect it to the WS client.

var address = flag.String("addr", ":8080", "HTTP Server")

func main() {
	flag.Parse()
	wss := newWebSocketServer()
	go wss.ServerLoop()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWS(wss, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I'm connected")
	})

	log.Fatal(http.ListenAndServe(*address, nil))
}
