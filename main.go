package main

import (
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var conn *pgx.Conn

var userImages = "/static/images/user_images/"

var linkPath = "/static/images/site_images/social_icons/SVG/"
var linksImagesPairs = map[string]string{
	"vk.com":        "vk.svg",
	"facebook.com":  "facebook.svg",
	"habr.com":      "habr.svg",
	"github.com":    "github.svg",
	"instagram.com": "instagram.svg",
	"linkedin.com":  "linkedin.svg",
	"ok.ru":         "ok.svg",
	"pinterest.com": "pinterest.svg",
	"reddit.com":    "reddit.svg",
	"snapchat.com":  "snapchat.svg",
	"tumblr.com":    "tumblr.svg",
	"twitter.com":   "twitter.svg",
	"youtube.com":   "youtube.svg",
	"t.me":          "telegram.svg",
	"other":         "other.svg",
}

func pageHandler(writer http.ResponseWriter, request *http.Request) {
	t, _ := template.ParseFiles("templates/page.html")

	userIdStr := request.URL.Path[1:]
	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		log.Println(err)
	}

	user := getUserViaId(userId)

	err = t.Execute(writer, user)

	if err != nil {
		log.Println(err)
	}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	conn = connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", faviconHandler)

	http.HandleFunc("/", pageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
