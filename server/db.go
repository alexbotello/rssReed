package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
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
	db.AutoMigrate(&Item{}, &Feed{})
	return
}

func addFeedToDB(feed string) {
	log.Printf("Adding %s into db", feed)
	var f Feed
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	if query := db.Where(&Feed{URL: feed}).First(&f); query.Error != nil {
		f = Feed{URL: feed}
		db.NewRecord(f)
		db.Create(&f)
		log.Println("Adding Feed URL into DB")
		return
	}
	log.Println("Feed could not be added to the database")
}

func addItemToDB(result *Result, s *Stream) {
	var rI Item
	var img string
	source := result.source
	item := result.item

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
	if query := db.Where(&Item{Title: item.Title}).First(&rI); query.Error != nil {
		rI = Item{
			Source: source,
			Title:  item.Title,
			Link:   item.Link,
			Desc:   item.Description,
			Date:   item.PublishedParsed,
			Image:  img,
		}
		db.NewRecord(rI)
		db.Create(&rI)
		s.client.send <- &rI
		log.Println("Adding Item Into DB")
	}
}

func getAllFeeds() []Feed {
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var feeds []Feed
	db.Find(&feeds)
	return feeds
}

func getAllRecords() []Item {
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var items []Item
	db.Find(&items)
	return items
}
