package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(writer http.ResponseWriter, _ *http.Request, status int) {
	writer.WriteHeader(status)

	if status == http.StatusNotFound {

		t, err := template.ParseFiles("templates/error404.html")
		if err != nil {
			log.Println(err)

		}

		err = t.Execute(writer, "")
		if err != nil {
			log.Println(err)
		}

	}
}
