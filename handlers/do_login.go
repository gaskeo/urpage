package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go-site/constants"
	"go-site/jwt"
	"go-site/storage"
	"go-site/verify_utils"
	"net/http"
)

func DoLogin(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var err error
	var jsonAnswer []byte

	if request.Method != "POST" {
		return
	}

	defer func() {
		writer.Header().Set("Content-Type", "application/json")
		_, _ = writer.Write(jsonAnswer)
	}()

	{ // work with form
		email := request.FormValue("email")
		password := request.FormValue("password")

		user, err = storage.GetUserByEmailAndPassword(email, password)

		if err != nil {
			fmt.Println(err)
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
		token, payload, tokenExpireDate, err := verify_utils.GenerateJWTToken(user.UserId)

		if err != nil {
			jsonAnswer, err = json.Marshal(Answer{Err: "other-error"})
			return
		}

		refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

		refreshToken, err := jwt.GenerateRefreshToken()

		if err != nil {
			jsonAnswer, err = json.Marshal(Answer{Err: "other-error"})
			return
		}

		verify_utils.AddJWTCookie(token, tokenExpireDate, writer)
		verify_utils.AddRefreshTokenCookie(refreshToken, payload.PayloadId, payload.UserId, refreshExpireDate, writer)

		err = verify_utils.AddJWSTokenInRedis(payload, token, tokenExpireDate)

		if err != nil {
			jsonAnswer, err = json.Marshal(Answer{Err: "other-error"})
			return
		}

		err = verify_utils.AddRefreshTokenInRedis(payload, refreshToken, refreshExpireDate)

		if err != nil {
			jsonAnswer, err = json.Marshal(Answer{Err: "other-error"})
			return
		}
	}

	jsonAnswer, err = json.Marshal(Answer{Err: ""})
}
