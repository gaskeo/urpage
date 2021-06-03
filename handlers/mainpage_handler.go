package handlers

import (
	"go-site/storage"
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	var temp string

	// auth here
	user := storage.User{UserId: 0, Username: "testUser"}

	if user.UserId == 0 {
		temp = "templates/index_not_auth.html"
	} else {
		temp = "templates/index_auth.html"
	}
	t, err := template.ParseFiles(temp)
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(writer, user)
}
