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

var results []*Result
var wg sync.WaitGroup
var pipe = make(chan *Result)

var rssfeeds = []string{
	"https://news.ycombinator.com/rss",
	"http://feeds.bbci.co.uk/news/world/rss.xml",
	"http://www.espn.com/espn/rss/news",
	"https://www.kut.org/rss.xml",
}

// Feed represents RSS feed url link
type Feed struct {
	gorm.Model
	url string
}

// Result represents a parsed RSS Feed result
type Result struct {
	source string
	item   *gofeed.Item
}

// Item represents the database model of a RSS Item
type Item struct {
	gorm.Model
	Source string
	Title  string
	Link   string
	Desc   string
	Date   *time.Time
	Image  string
}

func gatherFeeds(s *Stream) {
	// feeds := getAllFeeds()
	wg.Add(len(rssfeeds))
	for _, feed := range rssfeeds {
		go retrieve(feed)
	}
	go readFromPipe()
	wg.Wait()

	sort.Sort(byTime(results))
	for _, result := range results {
		addItemToDB(result, s)
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
	source := data.Title
	for idx, item := range data.Items {
		// Exit loop after 10 items
		if idx == 10 {
			break
		}
		pipe <- &Result{source: source, item: item}
	}
	fmt.Print(feed + " ")
	fmt.Println("Finished at ", time.Now())
}

func readFromPipe() {
	for r := range pipe {
		results = append(results, r)
	}
}

// byTime type is used for sorting feed items from newest to oldest
type byTime []*Result

func (t byTime) Len() int {
	return len(t)
}

func (t byTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t byTime) Less(i, j int) bool {
	return t[i].item.PublishedParsed.Before(*t[j].item.PublishedParsed)
}
