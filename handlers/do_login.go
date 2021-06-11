package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go-site/errs"
	"go-site/jwt"
	"go-site/redis_api"
	"go-site/session"
	"go-site/storage"
	"go-site/structs"
	"net/http"
	"time"
)

func CreateDoLogin(conn *pgx.Conn, rdb *redis.Client) {

	doLogin := func(writer http.ResponseWriter, request *http.Request) {
		var (
			email, password, token, CSRFToken, CSRFTokenForm, refreshToken string
			jsonAnswer                                                     []byte
			tokenExpireDate, refreshExpireDate                             time.Time
			user                                                           structs.User
			payload                                                        structs.Payload
			err                                                            error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
				return
			}
		}

		{ // work with form
			CSRFTokenForm = request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
				return
			}

			email = request.FormValue("email")
			password = request.FormValue("password")

			user, err = storage.GetUserByEmailAndPassword(conn, email, password)

			if err != nil {
				if err == errs.ErrWrongPassword {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "wrong-password"})
					return
				}
				if err == pgx.ErrNoRows {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "user-not-exist"})
					return
				} else {
					jsonAnswer, _ = json.Marshal(structs.Answer{Err: "other-error"})
					return
				}
			}
		}

		{ // work with JWT
			payload, token, tokenExpireDate, err = jwt.GenerateJWTToken(writer, user.UserId)

			if err != nil {
				jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}

			refreshToken, refreshExpireDate, err = jwt.GenerateRefreshToken(writer, payload)

			if err != nil {
				jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}

			err = redis_api.SetJWSToken(rdb, payload, token, tokenExpireDate)

			if err != nil {
				jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}

			err = redis_api.SetRefreshToken(rdb, payload, refreshToken, refreshExpireDate)

			if err != nil {
				jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
				return
			}
		}

		jsonAnswer, err = json.Marshal(structs.Answer{Err: ""})
	}

	http.HandleFunc("/do/login", doLogin)
}
