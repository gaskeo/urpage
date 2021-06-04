package handlers

import (
	"go-site/storage"
	"go-site/utils"
	"net/http"
)

func DoRegistration(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")

	passwordHashed, err := utils.HashPassword(password)

	if err != nil {
		http.Redirect(writer, request, "/registration", http.StatusSeeOther)
		return
	}

	userExist, _ := storage.CheckEmailExistInDB(email)

	if userExist {
		http.Redirect(writer, request, "/registration", http.StatusSeeOther)
		return
	}

	_, err = storage.AddUser(username, passwordHashed, email, "", "")

	if err != nil {
		http.Redirect(writer, request, "/registration", http.StatusSeeOther)
		return
	}

	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}
