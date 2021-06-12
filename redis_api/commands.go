package redis_api

import (
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func SetRefreshToken(
	rdb *redis.Client, payloadId int64, userId int, refreshToken string, refreshTokenExpireDate time.Time) error {
	refreshKey := strconv.FormatInt(payloadId, 10) + strconv.Itoa(userId) + "Refresh"

	return Set(rdb, refreshKey, refreshToken, refreshTokenExpireDate)

}

func SetJWSToken(rdb *redis.Client, payloadId int64, userId int, JWTToken string, tokenExpireDate time.Time) error {
	JWTKey := strconv.FormatInt(payloadId, 10) + strconv.Itoa(userId) + "JWT"

	return Set(rdb, JWTKey, JWTToken, tokenExpireDate)
}

func SetSession(rdb *redis.Client, sessionId, CSRFToken string, expireTime time.Time) error {
	return Set(rdb, sessionId, CSRFToken, expireTime)
}

func DeleteSession(rdb *redis.Client, sessionId string) error {
	return Set(rdb, sessionId, "", time.Now())
}

func DeleteJWTToken(rdb *redis.Client, payloadId int64, userId int) error {
	JWTKey := strconv.FormatInt(payloadId, 10) + strconv.Itoa(userId) + "JWT"

	return Set(rdb, JWTKey, "", time.Now())
}

func DeleteRefreshToken(rdb *redis.Client, payloadId int64, userId int) error {
	refreshKey := strconv.FormatInt(payloadId, 10) + strconv.Itoa(userId) + "Refresh"

	return Set(rdb, refreshKey, "", time.Now())
}
