package handlers

import (
	"go-site/constants"
	"go-site/jwt"
	"go-site/storage"
	"go-site/verify_utils"
	"log"
	"net/http"
	"strconv"
)

func DoLogin(writer http.ResponseWriter, request *http.Request) {
	var user storage.User
	var err error

	if request.Method != "POST" {
		return
	}

	{ // work with form
		email := request.FormValue("email")
		password := request.FormValue("password")

		user, err = storage.GetUserByEmailAndPassword(email, password)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}
	}

	{ // work with JWT
		token, payload, tokenExpireDate, err := verify_utils.GenerateJWTToken(user.UserId)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

		refreshToken, err := jwt.GenerateRefreshToken()

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		verify_utils.AddJWTCookie(token, tokenExpireDate, writer)
		verify_utils.AddRefreshTokenCookie(refreshToken, payload.PayloadId, payload.UserId, refreshExpireDate, writer)

		err = verify_utils.AddJWSTokenInRedis(payload, token, tokenExpireDate)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}

		err = verify_utils.AddRefreshTokenInRedis(payload, refreshToken, refreshExpireDate)

		if err != nil {
			log.Println(err)
			http.Redirect(writer, request, "/login", http.StatusSeeOther)
			return
		}
	}

	http.Redirect(writer, request, "/id/"+strconv.Itoa(user.UserId), http.StatusSeeOther)
}
