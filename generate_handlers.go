package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/handlers"
	"net/http"
)

func generateHandlers(conn *pgx.Conn, rds *redis.Client) {

	handlerNames := []func(conn *pgx.Conn, rds *redis.Client){
		handlers.CreateRegistrationHandler,
		handlers.CreateLoginHandler,
		handlers.CreateDoRegistration,
		handlers.CreateDoLogin,
		handlers.CreateMainPageHandler,
		handlers.CreateDoEditMain,
		handlers.CreateDoEditLinks,
		handlers.CreateDoEditPassword,
		handlers.CreateEditHandler,
		handlers.CreatePageHandler,
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	for _, handler := range handlerNames {
		handler(conn, rds)
	}
}
