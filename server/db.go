package main

import (
	"errors"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

func verifyDatabases() {
	// Check if database file exists; if not create it
	if _, err := os.Stat("./rss.db"); os.IsNotExist(err) {
		f, err := os.Create("rss.db")
		if err != nil {
			log.Fatal("Error creating rss.db file")
		}
		defer f.Close()
	}
	// Check that feeds database exists
	if _, err := os.Stat("./feeds.db"); os.IsNotExist(err) {
		f, err := os.Create("feeds.db")
		if err != nil {
			log.Fatal("Error creating feeds.db file")
		}
		defer f.Close()
		feedsDB, err := gorm.Open("sqlite3", "feeds.db")
		if err != nil {
			panic("failed to connect to database")
		}
		defer feedsDB.Close()
		feedsDB.AutoMigrate(&Feed{})
	}
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	db.AutoMigrate(&Item{})
	return
}

func addFeedToDB(f *Feed) error {
	var none Feed
	db, err := gorm.Open("sqlite3", "feeds.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	if query := db.Where(f).First(&none); query.Error != nil {
		db.NewRecord(&f)
		db.Create(f)
		return nil
	}
	return errors.New("Error saving feed to DB")
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

func getAllFeeds() ([]Feed, error) {
	db, err := gorm.Open("sqlite3", "feeds.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var feeds []Feed
	db.Find(&feeds)

	if len(feeds) < 1 {
		return nil, errors.New("No feeds saved inside database")
	}
	return feeds, nil
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
