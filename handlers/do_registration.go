package handlers

import (
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"net/http"
)

func DoRegistration(writer http.ResponseWriter, request *http.Request) {
	var username, email, password, passwordHashed, CSRFToken, CSRFTokenForm string
	var userExist bool

	var jsonAnswer []byte

	var err error

	if request.Method != "POST" {
		return
	}

	defer func() { SendJson(writer, jsonAnswer) }()

	{ // CSRF check
		_, CSRFToken, err = session.CheckSessionId(writer, request)

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}
	}

	{
		CSRFTokenForm = request.FormValue("csrf")

		if CSRFToken != CSRFTokenForm {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}

		username = request.FormValue("username")
		email = request.FormValue("email")
		password = request.FormValue("password")

		passwordHashed, err = storage.HashPassword(password)

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}
	}

	{ // email exist check
		userExist, err = storage.CheckEmailExistInDB(email)

		if err != nil && (userExist || err != pgx.ErrNoRows) {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "email-exist"})
			return
		}
	}

	{ // add user in db
		_, err = storage.AddUser(username, passwordHashed, email, "", "")

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}
	}

	jsonAnswer, err = json.Marshal(structs.Answer{Err: ""})
}
