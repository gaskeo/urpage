package handlers

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/jwt"
	"go-site/session"
	"go-site/structs"
	"html/template"
	"log"
	"net/http"
)

func CreateRegistrationHandler(_ *pgx.Conn, rdb *redis.Client) {
	registrationHandler := func(writer http.ResponseWriter, request *http.Request) {
		var (
			CSRFToken string
			err       error
		)

		{ // check csrf
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)
			if err != nil {
				http.Error(writer, "что-то пошло не так...", http.StatusInternalServerError)
				return
			}
		}

		{ // user auth check
			_, err = jwt.CheckIfUserAuth(writer, request, rdb)

			if err == nil {
				http.Redirect(writer, request, "/", http.StatusSeeOther)
				return
			}
		}

		{ // generate template
			t, err := template.ParseFiles("templates/registration.html")

			if err != nil {
				log.Println(err)
			}

			err = t.Execute(writer, structs.TemplateData{"CSRF": CSRFToken})

			if err != nil {
				log.Println(err)
			}
		}
	}

	http.HandleFunc("/registration", registrationHandler)
}
