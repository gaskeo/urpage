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

func GenerateJWTToken(userId int) (string, jwt.Payload, time.Time) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)
	refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

	payload, token, err := jwt.GenerateJWTToken(userId, tokenExpireDate, jwt.SecretKey)
	if err != nil {
		return "", jwt.Payload{}, time.Time{}
	}
	return token, payload, refreshExpireDate
}

func addJWTCookie(token,
	refreshToken string,
	tokenId int64,
	userId int,
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

	cookieRefreshId := http.Cookie{
		Name:    "RefreshTokenId",
		Value:   strconv.FormatInt(tokenId, 10),
		Expires: refreshExpireDate,
		Path:    "/",
	}

	cookieRefreshUserId := http.Cookie{
		Name:    "RefreshTokenUserId",
		Value:   strconv.Itoa(userId),
		Expires: refreshExpireDate,
		Path:    "/",
	}

	http.SetCookie(writer, &cookieToken)
	http.SetCookie(writer, &cookieRefresh)
	http.SetCookie(writer, &cookieRefreshId)
	http.SetCookie(writer, &cookieRefreshUserId)
}

func AddJWSTokenInRedis(payload jwt.Payload, JWTToken string, tokenExpireDate time.Time) {

	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	redis_api.Set(JWTKey, JWTToken, tokenExpireDate)
}

func AddRefreshTokenInRedis(payload jwt.Payload, refreshToken string, refreshTokenExpireDate time.Time) {

	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

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

	token, payload, tokenExpireDate := GenerateJWTToken(user.UserId)

	refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

	refreshToken := jwt.GenerateRefreshToken()

	addJWTCookie(token, refreshToken, payload.PayloadId, payload.UserId, tokenExpireDate, refreshExpireDate, writer)

	AddJWSTokenInRedis(payload, token, tokenExpireDate)

	AddRefreshTokenInRedis(payload, refreshToken, refreshExpireDate)

	http.Redirect(writer, request, "/id/"+strconv.Itoa(user.UserId), http.StatusSeeOther)
}
