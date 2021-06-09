package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go-site/constants"
	"go-site/errs"
	"go-site/jwt"
	"go-site/storage"
	"go-site/structs"
	"go-site/verify_utils"
	"net/http"
)

func DoLogin(writer http.ResponseWriter, request *http.Request) {
	var user structs.User
	var err error
	var CSRFToken, CSRFTokenForm string
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

	{ // work with form
		CSRFTokenForm = request.FormValue("csrf")

		if CSRFToken != CSRFTokenForm {
			jsonAnswer, _ = json.Marshal(structs.Answer{Err: "no-csrf"})
			return
		}

		email := request.FormValue("email")
		password := request.FormValue("password")

		user, err = storage.GetUserByEmailAndPassword(email, password)

		if err != nil {
			fmt.Println(err)
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
		token, payload, tokenExpireDate, err := verify_utils.GenerateJWTToken(user.UserId)

		if err != nil {
			jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}

		refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

		refreshToken, err := jwt.GenerateRefreshToken()

		if err != nil {
			jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}

		verify_utils.AddJWTCookie(token, tokenExpireDate, writer)
		verify_utils.AddRefreshTokenCookie(refreshToken, payload.PayloadId, payload.UserId, refreshExpireDate, writer)

		err = verify_utils.AddJWSTokenInRedis(payload, token, tokenExpireDate)

		if err != nil {
			jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}

		err = verify_utils.AddRefreshTokenInRedis(payload, refreshToken, refreshExpireDate)

		if err != nil {
			jsonAnswer, err = json.Marshal(structs.Answer{Err: "other-error"})
			return
		}
	}

	jsonAnswer, err = json.Marshal(structs.Answer{Err: ""})
}
