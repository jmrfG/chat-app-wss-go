package main

type WSServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newWebSocketServer() *WSServer {
	return &WSServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (wss *WSServer) registerClient(client *Client) {
	wss.clients[client] = true
}

func (wss *WSServer) unregisterClient(client *Client) {
	if _, ok := wss.clients[client]; ok {
		delete(wss.clients, client)
	}
}

func (wss *WSServer) broadcastToClients(msg []byte) {
	for client := range wss.clients {
		client.send <- msg
	}
}

// Runs the server and listens for new client requests
func (wss *WSServer) ServerLoop() {
	for {
		select {
		case client := <-wss.register:
			wss.registerClient(client)
		case client := <-wss.unregister:
			wss.unregisterClient(client)
		case msg := <-wss.broadcast:
			wss.broadcastToClients(msg)
		}
	}
}
