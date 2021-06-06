package handlers

import (
	"go-site/storage"
	"go-site/verify_utils"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func EditHandler(writer http.ResponseWriter, request *http.Request) {
	var authUser storage.User

	{ // user auth check
		authUserId, err := verify_utils.CheckIfUserAuth(writer, request)

		authUser, err = storage.GetUserViaId(authUserId)

		if err != nil {
			authUser = storage.User{}
		}
	}

	{ // check url
		requestedId, err := strconv.Atoi(request.URL.Path[len("/edit/"):])

		if err != nil {
			log.Println(err, 123)
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}

		if authUser.UserId != requestedId {
			log.Println("wrong id")
			ErrorHandler(writer, request, http.StatusForbidden)
			return
		}
	}

	{ // generate template
		t, err := template.ParseFiles("templates/edit_page.html")
		if err != nil {
			log.Println(err)

			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}

		err = t.Execute(writer, authUser)

		if err != nil {
			log.Println(err)
			ErrorHandler(writer, request, http.StatusNotFound)
			return
		}
	}
}
