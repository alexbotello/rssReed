package main

import (
	"fmt"
	"sort"
	"sync"

	"github.com/mmcdole/gofeed"
)

func main() {
	var items []*gofeed.Item
	var wg sync.WaitGroup
	pipe := make(chan *gofeed.Item)

	wg.Add(len(rssfeeds))

	for _, feed := range rssfeeds {
		go retrieveFeed(&wg, pipe, feed)
	}
	go readFromPipe(&items, pipe)
	wg.Wait()
	sort.Sort(byTime(items))
	for _, item := range items {
		fmt.Println(item.Title)
	}
}
