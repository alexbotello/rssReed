package main

import "github.com/gorilla/websocket"

type electron struct {
	conn *websocket.Conn
	send chan *RssItem
}

func (e *electron) write() {
	defer e.conn.Close()
	for item := range e.send {
		err := e.conn.WriteJSON(item)
		if err != nil {
			return
		}
	}
}
