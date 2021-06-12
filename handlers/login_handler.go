package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt_api"
	"go-site/session"
	"go-site/structs"
	"html/template"
	"net/http"
)

func CreateLoginHandler(_ *pgx.Conn, rdb *redis.Client) {

	loginHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			CSRFToken string
			err       error
		)

		{ // check csrf
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)
			if err != nil {
				http.Error(writer, "error session", http.StatusInternalServerError)
				return
			}
		}

		{ // check user authed
			_, err = jwt_api.CheckIfUserAuth(writer, request, rdb)

			if err == nil {
				http.Redirect(writer, request, "/", http.StatusSeeOther)
				return
			}
		}

		{ // generate login page
			t, err := template.ParseFiles("templates/login.html")

			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}

			err = t.Execute(writer, structs.TemplateData{"CSRF": CSRFToken})

			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}
		}
	}

	http.HandleFunc("/login", loginHandler)
}
