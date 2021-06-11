package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"html/template"
	"net/http"
)

func CreateMainPageHandler(conn *pgx.Conn, rdb *redis.Client) {

	mainPageHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			temp, CSRFToken string
			user            structs.User
			err             error
		)

		{ // check csrf
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)
			if err != nil {
				http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
				return
			}
		}

		{ // user auth check
			userId, err := jwt.CheckIfUserAuth(writer, request, rdb)

			if err != nil {
				temp = "templates/index_not_auth.html"

				user = structs.User{}
			} else {
				temp = "templates/index_auth.html"

				user, err = storage.GetUserViaId(conn, userId)
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

	http.HandleFunc("/", mainPageHandler)
}
