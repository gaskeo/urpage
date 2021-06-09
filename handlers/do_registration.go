package handlers

import (
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"go-site/storage"
	"go-site/structs"
	"go-site/verify_utils"
	"net/http"
)

func DoRegistration(writer http.ResponseWriter, request *http.Request) {
	var username, email, password, passwordHashed, CSRFToken, CSRFTokenForm string
	var err error
	var jsonAnswer []byte

	if request.Method != "POST" {
		return
	}

	defer func() {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonAnswer)
	}()

	{ // CSRF check
		_, CSRFToken, err = verify_utils.CheckSessionId(writer, request)

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

		passwordHashed, err = verify_utils.HashPassword(password)

		if err != nil {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}
	}

	{ // email exist check
		userExist, err := storage.CheckEmailExistInDB(email)

		if userExist || err != pgx.ErrNoRows {
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
