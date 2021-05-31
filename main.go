package main

import (
	"go-site/handlers"
	"go-site/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	storage.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	http.HandleFunc("/reg", handlers.RegistrationHandler)

	http.HandleFunc("/id/", handlers.PageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
