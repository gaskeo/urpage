package verify_utils

import (
	"go-site/constants"
	"go-site/jwt"
	"go-site/redis_api"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CheckIfUserAuth(request *http.Request) (int, error) {

	{ // check jwt token block
		JWTToken, err := request.Cookie("JWT")

		if err == nil {
			payload, err := jwt.VerifyToken(JWTToken.Value, jwt.SecretKey)

			if err == jwt.ErrInvalidToken {
				log.Println(err)
				return 0, err
			}

			if err != jwt.ErrExpiredToken && payload != nil {
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
					return 0, jwt.ErrInvalidToken
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
			return 0, jwt.ErrInvalidRefreshToken
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

		if err != nil {
			return 0, err
		}

		return newPayload.UserId, nil
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

func GenerateJWTToken(userId int) (string, jwt.Payload, time.Time, error) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)
	refreshExpireDate := tokenExpireDate.Add(constants.RefreshTokenExpireTime)

	payload, token, err := jwt.GenerateJWTToken(userId, tokenExpireDate, jwt.SecretKey)

	if err != nil {
		return "", jwt.Payload{}, time.Time{}, err
	}

	return token, payload, refreshExpireDate, nil
}

func AddJWTCookie(token,
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

func AddJWSTokenInRedis(payload jwt.Payload, JWTToken string, tokenExpireDate time.Time) error {
	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	err := redis_api.Set(JWTKey, JWTToken, tokenExpireDate)

	if err != nil {
		return err
	}

	return nil
}

func AddRefreshTokenInRedis(payload jwt.Payload, refreshToken string, refreshTokenExpireDate time.Time) error {
	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

	err := redis_api.Set(refreshKey, refreshToken, refreshTokenExpireDate)

	if err != nil {
		return err
	}

	return nil
}
