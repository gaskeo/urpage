package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(writer, "")

	if err != nil {
		log.Println(err)
	}
}
