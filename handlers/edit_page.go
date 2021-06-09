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

func CreateEditHandler(conn *pgx.Conn, rds *redis.Client) {

	editHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			authUserId, requestedId int
			CSRFToken               string
			t                       *template.Template
			authUser                structs.User
			err                     error
		)

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rds)

			if err != nil {
				http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
				return
			}
		}

		{ // user auth check
			authUserId, err = jwt.CheckIfUserAuth(writer, request, rds)

			authUser, err = storage.GetUserViaId(conn, authUserId)

			if err != nil {
				authUser = structs.User{}
			}
		}

		{ // check url
			requestedId, err = strconv.Atoi(request.URL.Path[len("/edit/"):])

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
			t, err = template.ParseFiles("templates/edit_page.html")
			if err != nil {
				log.Println(err)

				ErrorHandler(writer, request, http.StatusNotFound)
				return
			}

			err = t.Execute(writer, structs.TemplateData{"AuthUser": authUser, "CSRF": CSRFToken})

			if err != nil {
				log.Println(err)
				ErrorHandler(writer, request, http.StatusNotFound)
				return
			}
		}
	}

	http.HandleFunc("/edit/", editHandler)
}
