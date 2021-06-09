package jwt

import (
	"net/http"
	"strconv"
	"time"
)

func AddJWTCookie(token string, tokenExpireDate time.Time, writer http.ResponseWriter) {
	cookieToken := http.Cookie{
		Name:     "JWT",
		Value:    token,
		Expires:  tokenExpireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookieToken)
}

func AddRefreshTokenCookie(refreshToken string, tokenId int64, userId int, refreshExpireDate time.Time, writer http.ResponseWriter) {
	cookieRefresh := http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshToken,
		Expires:  refreshExpireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	cookieRefreshId := http.Cookie{
		Name:     "RefreshTokenId",
		Value:    strconv.FormatInt(tokenId, 10),
		Expires:  refreshExpireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	cookieRefreshUserId := http.Cookie{
		Name:     "RefreshTokenUserId",
		Value:    strconv.Itoa(userId),
		Expires:  refreshExpireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookieRefresh)
	http.SetCookie(writer, &cookieRefreshId)
	http.SetCookie(writer, &cookieRefreshUserId)
}
