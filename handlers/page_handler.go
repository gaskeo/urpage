package handlers

import (
	"go-site/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func PageHandler(writer http.ResponseWriter, request *http.Request) {
	var user storage.User

	{ // get user by user id in path
		userIdStr := request.URL.Path[len("/id/"):]

		if len(userIdStr) == 0 {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}

		userId, err := strconv.Atoi(userIdStr)

		if err != nil {
			log.Println(err)
		}

		user, err = storage.GetUserViaId(userId)

		if err != nil {
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}
	}

	{ // generate template
		t, err := template.ParseFiles("templates/page.html")

		if err != nil {
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}
		err = t.Execute(writer, user)

		if err != nil {
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}
	}
}
