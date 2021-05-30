package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/main_page.html")
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(writer, "")
}
