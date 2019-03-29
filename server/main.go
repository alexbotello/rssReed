package main

import (
	"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	stream := newStream()
	verifyDatabase()
	http.Handle("/", &jsonHandler{})
	http.Handle("/stream", stream)
	http.Handle("/exit", &exitHandler{})
	http.Handle("/save", &saveHandler{})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
