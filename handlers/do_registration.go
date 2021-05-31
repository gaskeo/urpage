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
		http.Redirect(writer, request, "/reg", http.StatusSeeOther)
		return
	}

	if storage.CheckEmailExistInDB(email) {
		http.Redirect(writer, request, "/reg", http.StatusSeeOther)
		return
	}

	storage.AddUser(username, passwordHashed, email, "", "")
	http.Redirect(writer, request, "/login", http.StatusSeeOther)
}
