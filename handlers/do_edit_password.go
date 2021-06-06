package handlers

import (
	"go-site/constants"
	"go-site/storage"
	"go-site/verify_utils"
	"log"
	"net/http"
	"strconv"
)

func DoEditPassword(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var userId int
	var userIdStr, oldPassword, newPassword string
	var err error

	if request.Method != "POST" {
		return
	}

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
	}

	{ // work with form
		userIdStr = request.FormValue("id")
		oldPassword = request.FormValue("old")
		newPassword = request.FormValue("new")
	}

	{ // compare form and authed user
		userIdIntForm, err := strconv.Atoi(userIdStr)

		if err != nil {
			log.Println(err, 13)
		}

		if userId != userIdIntForm {
			http.Redirect(writer, request, "/", http.StatusSeeOther)
			return
		}
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

		if err != nil {
			log.Println(err)
		}
	}

	{ // check old password
		correct, err := verify_utils.CheckPassword(oldPassword, user.Password)
		if err != nil || !correct {
			panic(err)
			return
		}
	}

	{ // set new data
		user.ImagePath = user.ImagePath[len(constants.UserImages):]

		user.Password, err = verify_utils.HashPassword(newPassword)

		if err != nil {
			panic(err)
			return
		}

		err = storage.UpdateUser(user)
		if err != nil {
			panic(err)
			return
		}
	}
}
