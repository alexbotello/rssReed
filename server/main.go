package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type jsonHandler struct {
	items []Item
	s     *Stream
}

func (j *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j.items = getAllRecords()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j.items)
}

type exitHandler struct{}

func (e *exitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := os.Remove("rss.db")
	if err != nil {
		log.Fatal("Failed to remove database upon close")
	}
	os.Exit(1)
}

func main() {
	stream := newStream()
	verifyDatabase()

	http.Handle("/", &jsonHandler{s: stream})
	http.Handle("/stream", stream)
	http.Handle("/exit", &exitHandler{})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
