package handlers

import (
	"go-site/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func PageHandler(writer http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("templates/page.html")

	if err != nil {
		log.Println(err)
	}

	userIdStr := request.URL.Path[len("/id/"):]

	if len(userIdStr) == 0 {
		MainPageHandler(writer, request)
		return
	}

	userId, err := strconv.Atoi(userIdStr)

	if err != nil {
		log.Println(err)
	}

	user := storage.GetUserViaId(userId)

	if user.UserId == 0 {
		ErrorHandler(writer, request, http.StatusNotFound)
		return
	}

	err = t.Execute(writer, user)

	if err != nil {
		log.Println(err)
	}
}
