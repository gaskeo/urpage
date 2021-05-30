package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func RegistrationHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/registration.html")
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(writer, "")

	if err != nil {
		log.Println(err)
	}
}
