package handlers

import (
	"go-site/constants"
	"go-site/jwt"
	"go-site/storage"
	"go-site/utils"
	"net/http"
	"strconv"
	"time"
)

func addJWTCookie(user storage.User, writer http.ResponseWriter, request *http.Request) {
	TokenExpireDate := time.Now().Add(constants.JWTExpireTime)
	RefreshExpireDate := TokenExpireDate.Add(constants.RefreshTokenExpireTime)

	token, err := jwt.GenerateJWTToken(user.UserId, TokenExpireDate, jwt.SecretKey)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	refreshToken := jwt.GenerateRefreshToken()

	cookieToken := http.Cookie{
		Name:    "JWT",
		Value:   token,
		Expires: TokenExpireDate,
		Path:    "/",
	}

	cookieRefresh := http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshToken,
		Expires: RefreshExpireDate,
		Path:    "/",
	}

	http.SetCookie(writer, &cookieToken)
	http.SetCookie(writer, &cookieRefresh)
}

func DoLogin(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	email := request.FormValue("email")
	password := request.FormValue("password")

	passwordHashed, err := utils.HashPassword(password)

	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	user := storage.GetUserByEmailAndPassword(email, passwordHashed)

	if user.UserId == 0 {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	addJWTCookie(*user, writer, request)

	http.Redirect(writer, request, "/id/"+strconv.Itoa(user.UserId), http.StatusSeeOther)
}
