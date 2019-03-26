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

type stream struct {
	initialLoad bool
	client      *electron
}

func newStream() *stream {
	return &stream{initialLoad: true}
}

func (s *stream) streamToSocket() {
	for {
		time.Sleep(30 * time.Second)
		gatherFeeds(s)
	}
}

func (s *stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("ServeHTTP Socket: ", err)
		return
	}
	electron := &electron{
		conn: socket,
		send: make(chan *RssItem, messageBufferSize),
	}
	s.client = electron

	go s.streamToSocket()
	electron.write()
}
