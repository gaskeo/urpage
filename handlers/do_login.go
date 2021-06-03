package handlers

import (
	"go-site/constants"
	"go-site/jwt"
	"go-site/redis_api"
	"go-site/storage"
	"go-site/utils"
	"net/http"
	"strconv"
	"time"
)

func addJWTCookie(token,
	refreshToken string,
	tokenExpireDate,
	refreshExpireDate time.Time,
	writer http.ResponseWriter) {
	cookieToken := http.Cookie{
		Name:    "JWT",
		Value:   token,
		Expires: tokenExpireDate,
		Path:    "/",
	}

	cookieRefresh := http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshToken,
		Expires: refreshExpireDate,
		Path:    "/",
	}

	http.SetCookie(writer, &cookieToken)
	http.SetCookie(writer, &cookieRefresh)
}

func AddTokensInRedis(payload jwt.Payload,
	JWTToken, refreshToken string,
	tokenExpireDate, refreshTokenExpireDate time.Time) {

	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

	redis_api.Set(JWTKey, JWTToken, tokenExpireDate)

	redis_api.Set(refreshKey, refreshToken, refreshTokenExpireDate)
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

	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)
	refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

	payload, token, err := jwt.GenerateJWTToken(user.UserId, tokenExpireDate, jwt.SecretKey)
	if err != nil {
		http.Redirect(writer, request, "/login", http.StatusSeeOther)
		return
	}

	refreshToken := jwt.GenerateRefreshToken()

	addJWTCookie(token, refreshToken, tokenExpireDate, refreshExpireDate, writer)

	AddTokensInRedis(payload, token, refreshToken, tokenExpireDate, refreshExpireDate)

	http.Redirect(writer, request, "/id/"+strconv.Itoa(user.UserId), http.StatusSeeOther)
}
