package handlers

import (
	"encoding/json"
	"go-site/constants"
	"go-site/storage"
	"go-site/verify_utils"
	"net/http"
	"strconv"
)

func DoEditPassword(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var userId int
	var userIdStr, oldPassword, newPassword string
	var err error

	var jsonAnswer []byte

	if request.Method != "POST" {
		return
	}

	defer func() {
		if len(jsonAnswer) > 0 {
			writer.Header().Set("Content-Type", "application/json")
			_, _ = writer.Write(jsonAnswer)
		}
	}()

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			http.Error(writer, "У вас нет доступа", http.StatusForbidden)
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
			jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
			return
		}

		if userId != userIdIntForm {
			http.Error(writer, "У вас нет доступа", http.StatusForbidden)
			return
		}
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

		if err != nil {
			http.Error(writer, "Ошибка с БД", http.StatusForbidden)
			return
		}
	}

	{ // check old password
		correct, err := verify_utils.CheckPassword(oldPassword, user.Password)
		if err != nil || !correct {
			jsonAnswer, _ = json.Marshal(Answer{Err: "wrong-password"})
			return
		}
	}

	{ // set new data
		user.ImagePath = user.ImagePath[len(constants.UserImages):]

		user.Password, err = verify_utils.HashPassword(newPassword)

		if err != nil {
			jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
			return
		}

		err = storage.UpdateUser(user)
		if err != nil {
			jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
			return
		}
		jsonAnswer, _ = json.Marshal(Answer{Err: ""})
	}
}
