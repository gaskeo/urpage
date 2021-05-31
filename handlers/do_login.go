package handlers

import (
	"go-site/storage"
	"go-site/utils"
	"net/http"
	"strconv"
)

func DoLogin(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	email := request.FormValue("email")
	password := request.FormValue("password")

	passwordHashed, err := utils.HashPassword(password)

	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user := storage.GetUserByEmailAndPassword(email, passwordHashed)

	if user.UserId == 0 {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	http.Redirect(writer, request, "/id/"+strconv.Itoa(user.UserId), http.StatusSeeOther)
}
