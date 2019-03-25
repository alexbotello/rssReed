package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

var rssfeeds = []string{
	"https://news.ycombinator.com/rss",
	"http://feeds.bbci.co.uk/news/world/rss.xml",
	"http://www.espn.com/espn/rss/news",
	"https://www.kut.org/rss.xml",
}

func retrieveFeed(w *sync.WaitGroup, c chan<- *gofeed.Item, feed string) {
	defer w.Done()
	resp, err := http.Get(feed)
	if err != nil {
		log.Printf("Request failed to find %s feed\n", feed)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fp := gofeed.NewParser()
	data, _ := fp.ParseString(string(body))
	for _, d := range data.Items {
		c <- d
	}
	fmt.Print(feed + " ")
	fmt.Println("Finished at ", time.Now())
}

func readFromPipe(items *[]*gofeed.Item, c <-chan *gofeed.Item) {
	for i := range c {
		*items = append(*items, i)
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
	return t[i].PublishedParsed.After(*t[j].PublishedParsed)
}
