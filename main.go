package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	items := getAllRecords()

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, items)
}

func main() {
	stream := newStream()
	verifyDatabase(stream)

	http.Handle("/", &templateHandler{filename: "rss.html"})
	http.Handle("/stream", stream)

	http.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		err := os.Remove("rss.db")
		if err != nil {
			log.Fatal("Failed to remove database upon close")
		}
		os.Exit(1)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
