package main

import (
	"go-site/handlers"
	"net/http"
)

func generateHandlers() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/favicon.ico", handlers.FaviconHandler)

	http.HandleFunc("/registration", handlers.RegistrationHandler)

	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("/do/registration", handlers.DoRegistration)

	http.HandleFunc("/do/login", handlers.DoLogin)

	http.HandleFunc("/do/edit_main", handlers.DoEditMain)

	http.HandleFunc("/edit/", handlers.EditHandler)

	http.HandleFunc("/id/", handlers.PageHandler)

	http.HandleFunc("/", handlers.MainPageHandler)
}
