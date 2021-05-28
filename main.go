package main

import (
	"html/template"
	"log"
	"net/http"
)

func pageHandler(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/page.html")
	title := request.URL.Path[1:]
	page := &Page{Title: title}

	err := t.Execute(writer, page)
	if err != nil {
		log.Fatal(err)
	}
}

type Page struct {
	Title string
}

func main() {
	http.HandleFunc("/", pageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
