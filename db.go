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
	return
}

func addItemToDB(db *gorm.DB, item *gofeed.Item) {
	var rI RssItem
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

func getAllRecords(db *gorm.DB) {
	var items []RssItem
	db.Find(&items)
	sort.Sort(byTime(items))
	for _, item := range items {
		// fmt.Println(item.Date.Format("Mon Jan _2 15:04:05 2006"))
		fmt.Println(item.Title)
	}
}
