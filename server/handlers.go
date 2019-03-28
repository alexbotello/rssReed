package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type jsonHandler struct {
	items []Feed
	s     *Stream
}

func (j *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gatherFeeds(j.s)
	j.items = getAllFeeds()
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

type saveHandler struct{}

func (s *saveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
