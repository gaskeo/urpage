package handlers

import (
	"encoding/json"
	"go-site/constants"
	"go-site/storage"
	"go-site/structs"
	"go-site/verify_utils"
	"net/http"
)

func DoEditPassword(writer http.ResponseWriter, request *http.Request) {
	var userId int
	var oldPassword, newPassword, CSRFToken, CSRFTokenForm string
	var correct bool

	var jsonAnswer []byte

	var user structs.User

	var err error

	if request.Method != "POST" {
		return
	}

	defer SendJson(writer, jsonAnswer)

	{ // CSRF check
		_, CSRFToken, err = verify_utils.CheckSessionId(writer, request)

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}
	}

	{ // check user authed
		userId, err = verify_utils.CheckIfUserAuth(writer, request)

		if err != nil {
			http.Error(writer, "У вас нет доступа", http.StatusForbidden)
			return
		}
	}

	{ // work with form
		CSRFTokenForm = request.FormValue("csrf")

		if CSRFToken != CSRFTokenForm {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}

		oldPassword = request.FormValue("old")
		newPassword = request.FormValue("new")
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

		if err != nil {
			http.Error(writer, "Ошибка с БД", http.StatusForbidden)
			return
		}
	}

	{ // check old password
		correct, err = verify_utils.CheckPassword(oldPassword, user.Password)
		if err != nil || !correct {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "wrong-password"})
			return
		}
	}

	{ // set new data
		user.ImagePath = user.ImagePath[len(constants.UserImages):]

		user.Password, err = verify_utils.HashPassword(newPassword)

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}

		err = storage.UpdateUser(user)
		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}
		jsonAnswer, _ = json.Marshal(structs.Answer{Err: ""})
	}
}
