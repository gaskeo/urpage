package handlers

import (
	"go-site/verify_utils"
	"html/template"
	"log"
	"net/http"
)

func LoginHandler(writer http.ResponseWriter, request *http.Request) {
	{ // check user authed
		_, err := verify_utils.CheckIfUserAuth(request)

		if err == nil {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
	}

	{ // generate login page
		t, err := template.ParseFiles("templates/login.html")

		if err != nil {
			log.Println(err)
		}

		err = t.Execute(writer, "")

		if err != nil {
			log.Println(err)
		}
	}
}
