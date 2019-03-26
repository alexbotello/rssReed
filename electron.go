package main

import "github.com/gorilla/websocket"

type electron struct {
	conn *websocket.Conn
	send chan *RssItem
}

func (e *electron) write() {
	defer e.conn.Close()
	for msg := range e.send {
		err := e.conn.WriteJSON(msg.Title)
		if err != nil {
			return
		}
	}
}
