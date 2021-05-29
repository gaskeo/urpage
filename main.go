package main

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"os"
)

var conn *pgx.Conn

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
	conn = connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	fmt.Println(getUserViaId(1)) // test
	http.HandleFunc("/", pageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
