package handlers

import (
	"go-site/storage"
	"go-site/structs"
	"go-site/verify_utils"
	"html/template"
	"net/http"
)

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	var temp, CSRFToken string
	var userId int
	var err error
	var user structs.User

	{ // check csrf
		_, CSRFToken, err = verify_utils.CheckSessionId(writer, request)
		if err != nil {
			http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
			return
		}
	}

	{ // user auth check
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			temp = "templates/index_not_auth.html"

			user = structs.User{}
		} else {
			temp = "templates/index_auth.html"

			user, err = storage.GetUserViaId(userId)
			if err != nil {
				user = structs.User{}
			}
		}
	}

	{ // generate template
		t, err := template.ParseFiles(temp)

		if err != nil {
			http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
			return
		}

		err = t.Execute(writer, structs.TemplateData{
			"User": user,
			"CSRF": CSRFToken,
		})
	}
}
