package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func (c *Connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		um.broadcast <- ConnectionMessage{c, string(message)}
	}
	c.ws.Close()
}

func (c *Connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &Connection{send: make(chan string, 10), ws: ws}
	um.register <- c
	defer func() { um.unregister <- c }()
	go c.writer()
	c.reader()
}

func main() {
	// Setup HTTP Connections
	go um.run()
	//http.HandleFunc("/", homeHandler)
	http.HandleFunc("/init", wsHandler)
	addr := ":6789"
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
