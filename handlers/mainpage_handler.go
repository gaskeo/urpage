package handlers

import (
	"go-site/jwt"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"html/template"
	"net/http"
)

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	var userId int
	var temp, CSRFToken string

	var t *template.Template

	var user structs.User

	var err error

	{ // check csrf
		_, CSRFToken, err = session.CheckSessionId(writer, request)
		if err != nil {
			http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
			return
		}
	}

	{ // user auth check
		userId, err = jwt.CheckIfUserAuth(writer, request)

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
		t, err = template.ParseFiles(temp)

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
