package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type jsonHandler struct {
	items []Feed
	s     *Stream
}

func (j *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Feed []Feed
		Item []Item
	}
	feeds := getAllFeeds()
	items := getAllRecords()
	resp := data{Feed: feeds, Item: items}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
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
	var feed Feed
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", 500)
		return
	}
	err = json.Unmarshal(body, &feed)
	if err != nil {
		http.Error(w, "JSON unmarshalling failed", 500)
		return
	}
	data, err := makeRequest(feed.URL)
	if err != nil {
		http.Error(w, "makeRequest failed to retrieve extra data", 500)
		return
	}
	feed.Name = data.Title
	err = addFeedToDB(&feed)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, "New feed %s successfully saved", feed.Name)
}
