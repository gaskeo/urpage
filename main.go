package main

import (
	"go-site/handlers"
	"go-site/redis_api"
	"go-site/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	_, err := storage.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal("database problem", err)
	}

	log.Println("connected to postgres")

	connectToRedis, err := redis_api.Connect(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_PASSWORD"), 0)
	if err != nil || !connectToRedis {
		log.Fatal("redis problem", err)
	}

	log.Println("connected to redis")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	http.HandleFunc("/registration", handlers.RegistrationHandler)

	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("/do/registration", handlers.DoRegistration)

	http.HandleFunc("/do/login", handlers.DoLogin)

	http.HandleFunc("/id/", handlers.PageHandler)

	http.HandleFunc("/", handlers.MainPageHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
