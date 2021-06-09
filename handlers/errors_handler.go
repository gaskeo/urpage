package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, status int) {
	var t *template.Template

	var err error

	writer.WriteHeader(status)

	if status == http.StatusNotFound {

		t, err = template.ParseFiles("templates/error404.html")
		if err != nil {
			log.Println(err)

		}

		err = t.Execute(writer, "")
		if err != nil {
			log.Println(err)
		}

	}
}
