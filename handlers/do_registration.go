package handlers

import (
	"github.com/jackc/pgx/v4"
	"go-site/storage"
	"go-site/verify_utils"
	"log"
	"net/http"
)

func DoRegistration(writer http.ResponseWriter, request *http.Request) {
	var username, email, password, passwordHashed string
	var err error
	if request.Method != "POST" {
		return
	}

	{
		username = request.FormValue("username")
		email = request.FormValue("email")
		password = request.FormValue("password")

		passwordHashed, err = verify_utils.HashPassword(password)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/registration", http.StatusSeeOther)
			return
		}
	}

	{ // email exist check
		userExist, err := storage.CheckEmailExistInDB(email)

		if userExist || err != pgx.ErrNoRows {
			log.Println(err)
			http.Redirect(writer, request, "/registration", http.StatusSeeOther)
			return
		}
	}

	{ // add user in db
		_, err = storage.AddUser(username, passwordHashed, email, "", "")

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/registration", http.StatusSeeOther)
			return
		}
	}

	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}
