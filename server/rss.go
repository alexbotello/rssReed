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
var userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36"

var rssfeeds = []string{
	"https://news.ycombinator.com/rss",
	"http://feeds.bbci.co.uk/news/world/rss.xml",
	"http://www.espn.com/espn/rss/news",
	"https://www.kut.org/rss.xml",
	"https://www.reddit.com/r/golang/new/.rss",
}

// Feed represents RSS feed url link
type Feed struct {
	gorm.Model
	URL  string
	Name string
}

// Result represents a parsed RSS Feed result
type Result struct {
	Source string
	Item   *gofeed.Item
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
	results = []*Result{}
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
		pipe <- &Result{Source: source, Item: item}
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
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Request failed to find feed")
		return nil, errors.New("makeRequest failed to complete")
	}

	request.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fp := gofeed.NewParser()
	data, err := fp.ParseString(string(body))
	if err != nil {
		er := fmt.Sprintf("Parsing response body failed: %s", err)
		return nil, errors.New(er)
	}
	return data, nil
}
