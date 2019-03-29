package main

import (
	"errors"
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
	URL  string
	Name string
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
	feeds, err := getAllFeeds()
	if err != nil {
		log.Println(err)
		return
	}
	wg.Add(len(feeds))
	for _, feed := range feeds {
		go retrieve(feed.URL)
	}
	go readFromPipe()
	wg.Wait()

	sort.Sort(byResult(results))
	for _, result := range results {
		addItemToDB(result, s)
	}
}

func retrieve(feed string) {
	defer wg.Done()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recover from err: %s", err)
		}
	}()
	data, err := makeRequest(feed)
	if err != nil {
		log.Println(err)
	}
	source := data.Title
	for idx, item := range data.Items {
		// Exit loop after 10 items
		if idx == 10 {
			break
		}
		pipe <- &Result{source: source, item: item}
	}
	addFeedToDB(&Feed{URL: feed, Name: source})
	fmt.Print(feed + " ")
	fmt.Println("Finished at ", time.Now())
}

func readFromPipe() {
	for r := range pipe {
		results = append(results, r)
	}
}

func makeRequest(url string) (*gofeed.Feed, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Request failed to find feed")
		return nil, errors.New("makeRequest failed to complete")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fp := gofeed.NewParser()
	data, err := fp.ParseString(string(body))
	if err != nil {
		log.Printf("Parsing response body failed:%s", err)
		return nil, errors.New("makeRequest failed to complete")
	}
	return data, nil
}
