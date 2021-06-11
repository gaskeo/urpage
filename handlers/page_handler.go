package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func CreatePageHandler(conn *pgx.Conn, rdb *redis.Client) {

	pageHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId, authUserId int
			userIdStr          string
			t                  *template.Template
			templateData       structs.TemplateData
			user               structs.User
			authUser           structs.User
			err                error
		)

		{ // CSRF check
			_, _, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
				return
			}
		}

		{ // get user by user id in path
			userIdStr = request.URL.Path[len("/id/"):]

			if len(userIdStr) == 0 {
				http.Redirect(writer, request, "/", http.StatusSeeOther)
				return
			}

			userId, err = strconv.Atoi(userIdStr)

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
			authUserId, err = jwt.CheckIfUserAuth(writer, request, rdb)

			authUser, err = storage.GetUserViaId(conn, authUserId)

			if err != nil {
				authUser = structs.User{}
			}
		}

		{ // generate template
			t, err = template.ParseFiles("templates/page.html")

			if err != nil {
				ErrorHandler(writer, request, http.StatusNotFound)
				return
			}

			templateData = structs.TemplateData{"User": user, "AuthUser": authUser}

			err = t.Execute(writer, templateData)

			if err != nil {
				ErrorHandler(writer, request, http.StatusNotFound)
				return
			}
		}
	}

	http.HandleFunc("/id/", pageHandler)
}
