package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	ping           = pongWait * 9 / 10
	maxMessageSize = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func (client *Client) readPump() {
	defer func() {
		//TODO: disconnect
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(appData string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("[Request Error]", err)
			break
		}
		//TODO: Broadcast
		client.wss.broadcast <- msg
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(ping)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("[ERR]", err)
				return
			}

			w.Write(msg)
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				log.Println("[CLOSE ERROR]", err)
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wss.unregister <- client
	close(client.send)
	client.conn.Close()
}
