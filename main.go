package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func pageHandler(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/page.html")
	title := request.URL.Path[1:]
	page := &WebPage{Title: title}

	err := t.Execute(writer, page)
	if err != nil {
		log.Fatal(err)
	}
}

type WebPage struct {
	Title string
}

func main() {
	conn := connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	fmt.Println(getUserViaId(conn, 15)) // test
	http.HandleFunc("/", pageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
