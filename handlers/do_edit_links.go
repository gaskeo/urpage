package handlers

import (
	"encoding/json"
	"go-site/constants"
	"go-site/storage"
	"go-site/structs"
	"go-site/utils"
	"go-site/verify_utils"
	"net/http"
	"strings"
)

func DoEditLinks(writer http.ResponseWriter, request *http.Request) {
	var userId int
	var links, CSRFToken, CSRFTokenForm string

	var linksLst []string
	var jsonAnswer []byte

	var user structs.User

	var err error

	if request.Method != "POST" {
		return
	}

	defer SendJson(writer, jsonAnswer)

	{ // check is session has CSRF
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

	{ // work with form and check CSRF
		CSRFTokenForm = request.FormValue("csrf")

		if CSRFToken != CSRFTokenForm {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}

		links = request.FormValue("links")
	}

	{ // get user
		user, err = storage.GetUserViaId(userId)

		if err != nil {
			http.Error(writer, "Ошибка с БД", http.StatusForbidden)
			return
		}
	}

	{ // set new data
		user.ImagePath = user.ImagePath[len(constants.UserImages):]

		linksLst = strings.Split(links, " ")

		user.Links, err = utils.CreateIconLinkPairs(linksLst)

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
