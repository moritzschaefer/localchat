package main

import (
	"log"
	"net/http"
	"text/template"

	"go/build"

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
			log.Fatal("Error while sending:", err)
			break
		}
	}
	c.ws.Close()
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
var homeTempl *template.Template

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error while connecting: ", err)
		return
	}
	c := &Connection{send: make(chan string, 10), ws: ws}
	um.register <- c
	defer func() { um.unregister <- c }()
	go c.writer()
	c.reader()
}

func defaultAssetPath() string {
	p, err := build.Default.Import("gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

func homeHandler(c http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(c, req.Host)
}

func main() {
	// Setup HTTP Connections
	go um.run()
	homeTempl = template.Must(template.ParseFiles("../webclient/index.html"))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/init", wsHandler)
	addr := ":6789"
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
