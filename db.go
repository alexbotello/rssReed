package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/jinzhu/gorm"
	"github.com/mmcdole/gofeed"
)

func verifyDatabase() {
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

func addItemToDB(item *gofeed.Item) {
	var rI RssItem
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	if result := db.Where(&RssItem{Title: item.Title}).First(&rI); result.Error != nil {
		rI = RssItem{
			Title: item.Title,
			Link:  item.Link,
			Desc:  item.Description,
			Date:  item.PublishedParsed,
		}
		db.NewRecord(rI)
		db.Create(&rI)
		fmt.Println("Adding item into DB")
	}
}

func getAllRecords() {
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	var items []RssItem
	db.Find(&items)
	sort.Sort(byTime(items))
	for _, item := range items {
		// fmt.Println(item.Date.Format("Mon Jan _2 15:04:05 2006"))
		fmt.Println(item.Title)
	}
}
