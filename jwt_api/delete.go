package jwt_api

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"urpage/redis_api"
)

func DeleteJWTToken(writer http.ResponseWriter, rdb *redis.Client, payload Payload) error {
	DeleteJWTCookie(writer)

	return redis_api.DeleteJWTToken(rdb, payload.PayloadId, payload.UserId)
}

func DeleteRefreshToken(writer http.ResponseWriter, rdb *redis.Client, payload Payload) error {
	DeleteRefreshTokenCookie(writer)

	return redis_api.DeleteRefreshToken(rdb, payload.PayloadId, payload.UserId)
}
