package main

import (
	"log"
	"net/http"
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
