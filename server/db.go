package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

func verifyDatabase() {
	// Check if database file exists; if not create it
	if _, err := os.Stat("./rss.db"); os.IsNotExist(err) {
		f, err := os.Create("rss.db")
		if err != nil {
			log.Fatal("Error creating rss.db file")
		}
		defer f.Close()
	}
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	db.AutoMigrate(&RssItem{})
	return
}

func addItemToDB(item *gofeed.Item, s *stream) {
	var rI RssItem
	var img string

	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	// Some RSS feeds may not provide an image
	if item.Extensions == nil {
		img = "None"
	} else {
		img = item.Extensions["media"]["thumbnail"][0].Attrs["url"]
	}
	// Only add RssItems that do not exist in the database
	if result := db.Where(&RssItem{Title: item.Title}).First(&rI); result.Error != nil {
		rI = RssItem{
			Title: item.Title,
			Link:  item.Link,
			Desc:  item.Description,
			Date:  item.PublishedParsed,
			Image: img,
		}
		db.NewRecord(rI)
		db.Create(&rI)
		s.client.send <- &rI
		log.Println("Adding Item Into DB")
	}
}

func getAllRecords() []RssItem {
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var items []RssItem
	db.Find(&items)
	return items
}
