package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
The Upgrader in Gorilla WebSocket is a struct that defines the parameters for upgrading an HTTP connection to a WebSocket connection.
It contains several fields that determine how the upgrade process should be handled, such as subprotocols,
compression options, and buffer sizes.
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		//Only for testing, refactor later
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	wss  *WSServer
	send chan []byte
}

func newClient(conn *websocket.Conn, wss *WSServer) *Client {
	return &Client{
		conn: conn,
		wss:  wss,
		send: make(chan []byte),
	}
}

// Administer WebSocket requests
func ServeWS(wss *WSServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("[Connection Error]", err)
		return
	}

	client := newClient(conn, wss)

	//TODO
	go client.writePump()
	go client.readPump()
	wss.register <- client
	fmt.Println("Client:", client, " has joined the chat!")
}
