package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"net/http"
	"urpage/jwt_api"
	"urpage/redis_api"
	"urpage/session"
	"urpage/storage"
)

func CreateDoLogin(conn *pgx.Conn, rdb *redis.Client) {

	doLogin := func(writer http.ResponseWriter, request *http.Request) {
		var (
			CSRFToken  string
			jsonAnswer []byte
			user       storage.User
			err        error
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

		{ // work with form
			CSRFTokenForm := request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}

			email := request.FormValue("email")
			password := request.FormValue("password")

			user, err = storage.GetUserByEmailAndPassword(conn, email, password)

			if err != nil {
				if err == storage.ErrWrongPassword {
					jsonAnswer, _ = json.Marshal(Answer{Err: "wrong-password"})
					return
				}
				if err == pgx.ErrNoRows {
					jsonAnswer, _ = json.Marshal(Answer{Err: "user-not-exist"})
					return
				} else {
					jsonAnswer, _ = json.Marshal(Answer{Err: "other-error"})
					return
				}
			}
		}

		{ // work with JWT
			payload, token, tokenExpireDate, err := jwt_api.GenerateJWTToken(writer, user.UserId)

			if err != nil {
				http.Error(writer, "error generating token", http.StatusInternalServerError)
				return
			}

			refreshToken, refreshExpireDate, err := jwt_api.GenerateRefreshToken(writer, payload)

			if err != nil {
				http.Error(writer, "error generating token", http.StatusInternalServerError)
				return
			}

			err = redis_api.SetJWSToken(rdb, payload.PayloadId, payload.UserId, token, tokenExpireDate)

			if err != nil {
				http.Error(writer, "error setting token", http.StatusInternalServerError)
				return
			}

			err = redis_api.SetRefreshToken(rdb, payload.PayloadId, payload.UserId, refreshToken, refreshExpireDate)

			if err != nil {
				http.Error(writer, "error setting token", http.StatusInternalServerError)
				return
			}
		}

		jsonAnswer, err = json.Marshal(Answer{Err: ""})
	}

	http.HandleFunc("/do/login", doLogin)
}
