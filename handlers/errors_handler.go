package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ErrorHandler(writer http.ResponseWriter, _ *http.Request, status int) {
	writer.WriteHeader(status)

	switch status {
	case http.StatusNotFound:
		t, err := template.ParseFiles("templates/error404.html")
		if err != nil {
			log.Println(err)

		}

		err = t.Execute(writer, "")
		if err != nil {
			log.Println(err)
		}
	case http.StatusInternalServerError:
		t, err := template.ParseFiles("templates/error500.html")
		if err != nil {
			log.Println(err)

		}

		err = t.Execute(writer, "")
		if err != nil {
			log.Println(err)
		}
	}
}
