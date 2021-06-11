package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"go-site/errs"
	"go-site/redis_api"
	"go-site/structs"
	"log"
	"net/http"
	"strconv"
)

func VerifyToken(token string) (*structs.Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errs.ErrInvalidToken
		}

		return []byte(SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &structs.Payload{}, keyFunc)

	if err != nil {
		ver, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(ver.Inner, errs.ErrExpiredToken) {
			return nil, errs.ErrExpiredToken
		}

		return nil, errs.ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*structs.Payload)

	if !ok {
		return nil, errs.ErrInvalidToken
	}

	return payload, nil
}

func CheckIfUserAuth(writer http.ResponseWriter, request *http.Request, rds *redis.Client) (int, error) {

	{ // check jwt token block
		JWTToken, err := request.Cookie("JWT")

		if err == nil {
			payload, err := VerifyToken(JWTToken.Value)

			if err == errs.ErrInvalidToken {
				log.Println(err)
				return 0, err
			}

			if err != errs.ErrExpiredToken && payload != nil {

				redisJWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"
				redisJWTValue, err := redis_api.Get(rds, redisJWTKey)

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
		redisRefreshTokenValue, err := redis_api.Get(rds, redisRefreshTokenKey)

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

		newPayload, newToken, newExpireDate, err := GenerateJWTToken(writer, userId)

		if err != nil {
			return 0, err
		}

		log.Println("generate new token")
		err = redis_api.SetJWSToken(rds, newPayload, newToken, newExpireDate)
		AddJWTCookie(writer, newToken, newExpireDate)

		if err != nil {
			return 0, err
		}

		return newPayload.UserId, nil
	}
}
