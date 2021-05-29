package main

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var conn *pgx.Conn
var userImages string = "/static/images/user_images/"

func pageHandler(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/page.html")
	userIdStr := request.URL.Path[1:]

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println(err)
	}

	user := getUserViaId(userId)

	err = t.Execute(writer, user)
	if err != nil {
		log.Fatal(err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	conn = connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	fmt.Println(getUserViaId(9)) // test

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", faviconHandler)

	http.HandleFunc("/", pageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
