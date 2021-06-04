package handlers

import (
	"go-site/jwt"
	"go-site/redis_api"
	"go-site/storage"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func checkIfUserAuth(request *http.Request) int {
	JWTToken, err := request.Cookie("JWT")
	if err == nil {
		payload, err := jwt.VerifyToken(JWTToken.Value, jwt.SecretKey)

		if err == jwt.ErrInvalidToken {
			log.Println(err)
			return 0
		}

		if err != jwt.ErrExpiredToken && payload != nil {
			redisJWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"
			redisJWTValue, err := redis_api.Get(redisJWTKey)
			if err != nil {
				log.Println(err)
				return 0
			}

			if redisJWTValue == JWTToken.Value {
				return payload.UserId
			} else {
				log.Println("invalid token")
				return 0
			}
		}
	}

	log.Println("no jwtToken")

	refreshToken, err := request.Cookie("RefreshToken")
	if err != nil {
		log.Println(err)
		return 0
	}

	refreshTokenId, err := request.Cookie("RefreshTokenId")
	if err != nil {
		log.Println(err)
		return 0
	}

	refreshTokenUserId, err := request.Cookie("RefreshTokenUserId")
	if err != nil {
		log.Println(err)
		return 0
	}

	redisRefreshTokenKey := refreshTokenId.Value + refreshTokenUserId.Value + "Refresh"
	redisRefreshTokenValue, err := redis_api.Get(redisRefreshTokenKey)
	if err != nil {
		log.Println(err)
		return 0
	}

	if refreshToken.Value != redisRefreshTokenValue {
		log.Println("invalid refresh token")
		return 0
	}

	userId, err := strconv.Atoi(refreshTokenUserId.Value)
	if err != nil {
		log.Println(err)
		return 0
	}

	newToken, newPayload, newExpireDate, err := GenerateJWTToken(userId)

	if err != nil {
		return 0
	}

	log.Println("generate new token")
	err = AddJWSTokenInRedis(newPayload, newToken, newExpireDate)

	if err != nil {
		return 0
	}

	return userId
}

func MainPageHandler(writer http.ResponseWriter, request *http.Request) {
	var temp string
	var user storage.User

	userId := checkIfUserAuth(request)
	if userId == 0 {
		user = storage.User{}
	} else {
		user, _ = storage.GetUserViaId(userId)
	}

	if user.UserId == 0 {
		temp = "templates/index_not_auth.html"
	} else {
		temp = "templates/index_auth.html"
	}
	t, err := template.ParseFiles(temp)
	if err != nil {
		log.Println(err)
	}

	err = t.Execute(writer, user)
}
