package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

var items []*gofeed.Item
var wg sync.WaitGroup
var pipe = make(chan *gofeed.Item)

var rssfeeds = []string{
	"https://news.ycombinator.com/rss",
	"http://feeds.bbci.co.uk/news/world/rss.xml",
	"http://www.espn.com/espn/rss/news",
	"https://www.kut.org/rss.xml",
}

// RssItem represents the database model for an Rss Item
type RssItem struct {
	gorm.Model
	Title string
	Link  string
	Desc  string
	Date  *time.Time
	Image string
}

func gatherFeeds(s *stream) {
	wg.Add(len(rssfeeds))
	for _, feed := range rssfeeds {
		go retrieve(feed)
	}
	go readFromPipe()
	wg.Wait()

	for _, item := range items {
		addItemToDB(item, s)
	}
	s.initialLoad = false
}

func retrieve(feed string) {
	defer wg.Done()
	resp, err := http.Get(feed)
	if err != nil {
		log.Printf("Request failed to find %s feed\n", feed)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fp := gofeed.NewParser()
	data, err := fp.ParseString(string(body))
	if err != nil {
		log.Printf("Parsing response body failed: %s", err)
		return
	}
	for _, d := range data.Items {
		pipe <- d
	}
	fmt.Print(feed + " ")
	fmt.Println("Finished at ", time.Now())
}

func readFromPipe() {
	for i := range pipe {
		items = append(items, i)
	}
}
