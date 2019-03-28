package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	socketBuffersize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBuffersize, WriteBufferSize: socketBuffersize}

// Stream represents a websocket connection between client and server
type Stream struct {
	client *electron
}

func newStream() *Stream {
	return &Stream{}
}

func (s *Stream) streamToSocket() {
	for {
		gatherFeeds(s)
		time.Sleep(30 * time.Second)
	}
}

func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP Socket: ", err)
		return
	}
	electron := &electron{
		conn: socket,
		send: make(chan *Item, messageBufferSize),
	}
	s.client = electron

	go s.streamToSocket()
	electron.write()
}
