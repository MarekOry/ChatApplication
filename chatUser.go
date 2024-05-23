package main

import (
	"github.com/gorilla/websocket"
)

type chatUser struct {
	nick     string
	chatRoom *chatRoom
	conn     *websocket.Conn
	receive  chan []byte
}

func (c *chatUser) read() {
	defer c.conn.Close()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.chatRoom.channel <- msg
	}
}

func (c *chatUser) write() {
	defer c.conn.Close()
	for msg := range c.receive {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
