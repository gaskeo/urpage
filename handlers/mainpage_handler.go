package handlers

import (
	"go-site/storage"
	"go-site/verify_utils"
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	var temp string
	var userId int
	var err error
	var user storage.User

	{ // user auth check
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			temp = "templates/index_not_auth.html"

			user = storage.User{}
		} else {
			temp = "templates/index_auth.html"

			user, err = storage.GetUserViaId(userId)
			if err != nil {
				user = storage.User{}
			}
		}
	}

	{ // generate template
		t, err := template.ParseFiles(temp)

		if err != nil {
			log.Println(err)
		}

		err = t.Execute(writer, user)
	}
}
