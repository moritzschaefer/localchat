package main

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type ConnectionMessage struct {
	connection *Connection
	message    string
}

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
		case c := <-um.unregister:
			if _, ok := um.users[c]; ok {
				delete(um.users, c)
				close(c.send)
			}
		case messageStruct := <-um.broadcast:
			um.sendMessage(messageStruct.connection, messageStruct.message)
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
func (um *UserManager) sendMessage(senderConnection *Connection, message string) {
	// Read Position of user and find other users within range
	for senderConnection := range um.users {
		select {
		case senderConnection.send <- message:
		default:
			delete(um.users, senderConnection)
			close(senderConnection.send)
		}
	}
}
