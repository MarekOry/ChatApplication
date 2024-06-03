package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type message struct {
	chatUsername string
	text         string
}

type chatRoom struct {
	chatUsers map[*chatUser]bool
	channel   chan []byte
	enter     chan *chatUser
	leave     chan *chatUser
}

func createChatRoom() *chatRoom {
	return &chatRoom{
		channel:   make(chan []byte),
		enter:     make(chan *chatUser),
		leave:     make(chan *chatUser),
		chatUsers: make(map[*chatUser]bool),
	}
}

func (r *chatRoom) run() {
	for {
		select {
		case chatUser := <-r.enter:
			r.chatUsers[chatUser] = true
		case chatUser := <-r.leave:
			delete(r.chatUsers, chatUser)
			close(chatUser.receive)
		case msg := <-r.channel:
			for client := range r.chatUsers {
				client.receive <- msg
			}
		}
	}
}

var upgrader = websocket.Upgrader{ReadBufferSize: 0, WriteBufferSize: 0}

func (r *chatRoom) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	chatUser := &chatUser{
		conn:    conn,
		receive: make(chan []byte, upgrader.WriteBufferSize),

		chatRoom: r,
	}
	r.enter <- chatUser
	defer func() { r.leave <- chatUser }()
	go chatUser.write()
	chatUser.read()
}
