package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"urpage/jwt_api"
	"urpage/session"
	"urpage/storage"
)

func CreatePageHandler(conn *pgx.Conn, rdb *redis.Client) {

	pageHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			user, authUser storage.User
			err            error
		)

		{ // CSRF check
			_, _, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				http.Error(writer, "error session", http.StatusInternalServerError)
				return
			}
		}

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

			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				ErrorHandler(writer, request, http.StatusNotFound)
				return
			}
		}

		{ // user auth check
			authUserId, err := jwt_api.CheckIfUserAuth(writer, request, rdb)

			if err != nil {
				authUser = storage.User{}
			} else {
				authUser, _ = storage.GetUserViaId(conn, authUserId)
			}

		}

		{ // generate template
			t, err := template.ParseFiles("templates/page.html")

			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}

			templateData := TemplateData{"User": user, "AuthUser": authUser}

			err = t.Execute(writer, templateData)

			if err != nil {
				http.Error(writer, "error creating page", http.StatusInternalServerError)
				return
			}
		}
	}

	http.HandleFunc("/id/", pageHandler)
}
