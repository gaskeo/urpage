package handlers

import (
	"go-site/verify_utils"
	"html/template"
	"log"
	"net/http"
)

func RegistrationHandler(writer http.ResponseWriter, request *http.Request) {
	{ // user auth check
		_, err := verify_utils.CheckIfUserAuth(request)

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

		err = t.Execute(writer, "")

		if err != nil {
			log.Println(err)
		}
	}
}
