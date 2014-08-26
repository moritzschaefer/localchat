package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan string
}

type UserManager struct {
	users map[*Connection]*User
	// Inbound messages from the connections.
	broadcast chan ConnectionMessage

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection
}

var um = UserManager{
	broadcast:  make(chan ConnectionMessage),
	register:   make(chan *Connection),
	unregister: make(chan *Connection),
	users:      make(map[*Connection]*User),
}

func (um *UserManager) run() {
	for {
		select {
		case c := <-um.register:
			um.users[c] = new(User)
			log.Print("User connected")
		case c := <-um.unregister:
			if _, ok := um.users[c]; ok {
				delete(um.users, c)
				close(c.send)
				log.Print("User disconnected")
			}
		case messageStruct := <-um.broadcast:

			log.Printf("A message was received: %v", messageStruct)
			if err := um.processMessage(messageStruct.connection, messageStruct.message); err != nil {
				log.Fatalf("Error processing message: %e", err)
			}
		}
	}
}

func (um *UserManager) addUser(c *Connection) error {
	if _, ok := um.users[c]; ok {
		return errors.New("Connection already exists")
	}
	um.users[c] = new(User)

	return nil
}

func (um *UserManager) updateUser(c *Connection, field string, value string) error {
	// TODO: check if user exists already
	if _, ok := um.users[c]; !ok {
		return errors.New("Connection is not registered. Call addUser first")
	}
	log.Printf("user sent update %s with value %s", field, value)
	if field == "username" {
		um.users[c].username = value
	} else if field == "position" {
		if pos, err := ParsePosition(value); err != nil {
			return err
		} else {
			um.users[c].position = pos
		}
	} else {
		return fmt.Errorf("Field '%s' is not known")
	}
	return nil
}
func (um *UserManager) processMessage(senderConnection *Connection, message string) error {
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(message), &dat); err != nil {
		return err // Write some error info
	}
	if dat["action"] == "update" {
		if err := um.updateUser(senderConnection, dat["field"].(string), dat["value"].(string)); err != nil {
			return err
		}
	} else if dat["action"] == "message" {
		um.sendMessage(senderConnection, message) //send message or dat["value"]
	} else {
		return fmt.Errorf("action %s not supported", dat["action"])
	}
	return nil
}
func (um *UserManager) sendMessage(senderConnection *Connection, message string) {
	// Read Position of user and find other users within range
	senderPos := um.users[senderConnection].position
	for connection, user := range um.users {
		if senderConnection != connection { // don't send to sender
			emptyPos := Position{0, 0}
			if user.username != "" && user.position != emptyPos && user.position.InRadius(senderPos, float64(MAX_DISTANCE)) {
				select {
				case connection.send <- message:
				default:
					delete(um.users, connection)
					close(connection.send)
				}
			}
		}
	}
}
