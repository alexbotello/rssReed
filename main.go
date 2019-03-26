package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	verifyDatabase()
	db, err := gorm.Open("sqlite3", "rss.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()
	db.AutoMigrate(&RssItem{})
	getAllRecords(db)
	runProcess(db)
}
