package verify_utils

import (
	"go-site/constants"
	"go-site/errs"
	"go-site/jwt"
	"go-site/redis_api"
	"go-site/session"
	"go-site/structs"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CheckIfUserAuth(writer http.ResponseWriter, request *http.Request) (int, error) {

	{ // check jwt token block
		JWTToken, err := request.Cookie("JWT")

		if err == nil {
			payload, err := jwt.VerifyToken(JWTToken.Value, jwt.SecretKey)

			if err == errs.ErrInvalidToken {
				log.Println(err)
				return 0, err
			}

			if err != errs.ErrExpiredToken && payload != nil {

				redisJWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"
				redisJWTValue, err := redis_api.Get(redisJWTKey)

				if err != nil {
					log.Println(err)
					return 0, err
				}

				if redisJWTValue == JWTToken.Value {
					return payload.UserId, nil

				} else {
					log.Println("invalid token")
					return 0, errs.ErrInvalidToken
				}
			}
		}
	}

	log.Println("no jwtToken")

	{ // check refresh token block
		refreshToken, err := request.Cookie("RefreshToken")

		if err != nil {
			log.Println(err)
			return 0, err
		}

		refreshTokenId, err := request.Cookie("RefreshTokenId")

		if err != nil {
			log.Println(err)
			return 0, err
		}

		refreshTokenUserId, err := request.Cookie("RefreshTokenUserId")

		if err != nil {
			log.Println(err)
			return 0, err
		}

		redisRefreshTokenKey := refreshTokenId.Value + refreshTokenUserId.Value + "Refresh"
		redisRefreshTokenValue, err := redis_api.Get(redisRefreshTokenKey)

		if err != nil {
			log.Println(err)
			return 0, err
		}

		if refreshToken.Value != redisRefreshTokenValue {
			log.Println("invalid refresh token")
			return 0, errs.ErrInvalidRefreshToken
		}

		userId, err := strconv.Atoi(refreshTokenUserId.Value)
		if err != nil {
			log.Println(err)
			return 0, err
		}

		newToken, newPayload, newExpireDate, err := GenerateJWTToken(userId)

		if err != nil {
			return 0, err
		}

		log.Println("generate new token")
		err = AddJWSTokenInRedis(newPayload, newToken, newExpireDate)
		AddJWTCookie(newToken, newExpireDate, writer)

		if err != nil {
			return 0, err
		}

		return newPayload.UserId, nil
	}
}

func CheckSessionId(writer http.ResponseWriter, request *http.Request) (string, string, error) {
	{ // check cookie
		sessionIdCookie, err := request.Cookie("SessionId")
		if err == nil {
			CSRFToken, err := session.GetCSRFBySessionId(sessionIdCookie.Value)
			return sessionIdCookie.Value, CSRFToken, err
		}
		sessionId := session.GenerateSessionId()
		CSRFToken := session.GenerateCSRFToken()
		expireTime := time.Now().Add(constants.SessionIdExpireTime)
		err = session.AddInRedis(sessionId, CSRFToken, expireTime)

		AddSessionId(sessionId, expireTime, writer)

		if err != nil {
			return "", "", err
		}
		return sessionId, CSRFToken, nil
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password string, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func GenerateJWTToken(userId int) (string, structs.Payload, time.Time, error) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)

	payload, token, err := jwt.GenerateJWTToken(userId, tokenExpireDate, jwt.SecretKey)

	if err != nil {
		return "", structs.Payload{}, time.Time{}, err
	}

	return token, payload, tokenExpireDate, nil
}

func AddSessionId(sessionId string, expireDate time.Time, writer http.ResponseWriter) {
	cookieSessionId := http.Cookie{
		Name:     "SessionId",
		Value:    sessionId,
		Expires:  expireDate,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(writer, &cookieSessionId)
}

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

func AddJWSTokenInRedis(payload structs.Payload, JWTToken string, tokenExpireDate time.Time) error {
	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	err := redis_api.Set(JWTKey, JWTToken, tokenExpireDate)

	if err != nil {
		return err
	}

	return nil
}

func AddRefreshTokenInRedis(payload structs.Payload, refreshToken string, refreshTokenExpireDate time.Time) error {
	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

	err := redis_api.Set(refreshKey, refreshToken, refreshTokenExpireDate)

	if err != nil {
		return err
	}

	return nil
}
