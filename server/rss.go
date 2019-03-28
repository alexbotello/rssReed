package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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

	sort.Sort(byTime(items))
	for _, item := range items {
		addItemToDB(item, s)
	}
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
	for idx, data := range data.Items {
		// Exit loop after 10 items
		if idx == 10 {
			break
		}
		pipe <- data
	}
	fmt.Print(feed + " ")
	fmt.Println("Finished at ", time.Now())
}

func readFromPipe() {
	for i := range pipe {
		items = append(items, i)
	}
}

// byTime type is used for sorting feed items from newest to oldest
type byTime []*gofeed.Item

func (t byTime) Len() int {
	return len(t)
}

func (t byTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byTime) Less(i, j int) bool {
	return t[i].PublishedParsed.Before(*t[j].PublishedParsed)
}
