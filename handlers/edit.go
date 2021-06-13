package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"html/template"
	"net/http"
	"strconv"
	"urpage/jwt_api"
	"urpage/session"
	"urpage/storage"
)

func CreateEditHandler(conn *pgx.Conn, rdb *redis.Client) {

	editHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			CSRFToken string
			authUser  storage.User
			err       error
		)

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				http.Error(writer, "error session", http.StatusInternalServerError)
				return
			}
		}

		{ // user auth check
			authUserId, err := jwt_api.CheckIfUserAuth(writer, request, rdb)

			if err != nil {
				http.Error(writer, "no jwt", http.StatusInternalServerError)
				return
			}

			authUser, err = storage.GetUserViaId(conn, authUserId)

			if err != nil {
				http.Error(writer, "error getting user", http.StatusInternalServerError)
				return
			}
		}

		{ // check url
			requestedId, err := strconv.Atoi(request.URL.Path[len("/edit/"):])

			if err != nil {
				http.Error(writer, "error getting id", http.StatusInternalServerError)
				return
			}

			if authUser.UserId != requestedId {
				http.Error(writer, "error checking id", http.StatusInternalServerError)
				return
			}
		}

		{ // generate template
			t, err := template.ParseFiles("templates/edit.html")
			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}

			err = t.Execute(writer, TemplateData{"AuthUser": authUser, "CSRF": CSRFToken})

			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}
		}
	}

	http.HandleFunc("/edit/", editHandler)
}
