package handlers

import (
	"go-site/storage"
	"go-site/verify_utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func PageHandler(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var authUser storage.User
	var users storage.SomeUsers

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

	{ // user auth check
		authUserId, err := verify_utils.CheckIfUserAuth(writer, request)

		authUser, err = storage.GetUserViaId(authUserId)

		if err != nil {
			authUser = storage.User{}
		}
	}

	{ // generate template
		t, err := template.ParseFiles("templates/page.html")

		if err != nil {
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}

		users = storage.SomeUsers{"User": user, "AuthUser": authUser}

		err = t.Execute(writer, users)

		if err != nil {
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}
	}
}
