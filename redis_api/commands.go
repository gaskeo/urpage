package redis_api

import (
	"github.com/go-redis/redis/v8"
	"go-site/structs"
	"strconv"
	"time"
)

func SetRefreshToken(rdb *redis.Client, payload structs.Payload, refreshToken string, refreshTokenExpireDate time.Time) error {
	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

	return Set(rdb, refreshKey, refreshToken, refreshTokenExpireDate)

}

func SetJWSToken(rdb *redis.Client, payload structs.Payload, JWTToken string, tokenExpireDate time.Time) error {
	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	return Set(rdb, JWTKey, JWTToken, tokenExpireDate)
}

func SetSession(rdb *redis.Client, sessionId, CSRFToken string, expireTime time.Time) error {
	return Set(rdb, sessionId, CSRFToken, expireTime)
}
