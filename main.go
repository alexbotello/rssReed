package main

import (
	"encoding/json"
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

type jsonHandler struct {
	items []RssItem
}

func (j *jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	j.items = getAllRecords()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j.items)
}

type exitHandler struct{}

func (e *exitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := os.Remove("rss.db")
	if err != nil {
		log.Fatal("Failed to remove database upon close")
	}
	os.Exit(1)
}

func main() {
	stream := newStream()
	verifyDatabase(stream)

	http.Handle("/", &jsonHandler{})
	http.Handle("/stream", stream)
	http.Handle("/exit", &exitHandler{})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
