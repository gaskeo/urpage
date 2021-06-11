package jwt

import (
	"github.com/go-redis/redis/v8"
	"go-site/redis_api"
	"go-site/structs"
	"net/http"
)

func DeleteJWTToken(writer http.ResponseWriter, rdb *redis.Client, payload structs.Payload) error {
	DeleteJWTCookie(writer)
	return redis_api.DeleteJWTToken(rdb, payload)
}

func DeleteRefreshToken(writer http.ResponseWriter, rdb *redis.Client, payload structs.Payload) error {
	DeleteRefreshTokenCookie(writer)
	return redis_api.DeleteRefreshToken(rdb, payload)
}
