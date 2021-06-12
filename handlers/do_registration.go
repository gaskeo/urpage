package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/session"
	"go-site/storage"
	"net/http"
)

func CreateDoRegistration(conn *pgx.Conn, rdb *redis.Client) {

	doRegistration := func(writer http.ResponseWriter, request *http.Request) {
		var (
			username, email, passwordHashed, CSRFToken string
			jsonAnswer                                 []byte
			err                                        error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}
		}

		{
			CSRFTokenForm := request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}

			username = request.FormValue("username")
			email = request.FormValue("email")
			password := request.FormValue("password")

			passwordHashed, err = storage.HashPassword(password)

			if err != nil {
				http.Error(writer, "error hashing password", http.StatusInternalServerError)
				return
			}
		}

		{ // email exist check
			_, err := storage.CheckEmailExistInDB(conn, email)

			if err != pgx.ErrNoRows {
				jsonAnswer, _ = json.Marshal(Answer{Err: "email-exist"})
				return
			}
		}

		{ // add user in db
			_, err = storage.AddUser(conn, username, passwordHashed, email, "", "")

			if err != nil {
				http.Error(writer, "error adding user", http.StatusInternalServerError)
				return
			}
		}

		jsonAnswer, err = json.Marshal(Answer{Err: ""})
	}

	http.HandleFunc("/do/registration", doRegistration)
}
