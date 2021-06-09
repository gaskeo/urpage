package redis_api

import (
	"go-site/structs"
	"strconv"
	"time"
)

func SetRefreshToken(payload structs.Payload, refreshToken string, refreshTokenExpireDate time.Time) error {
	refreshKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "Refresh"

	return Set(refreshKey, refreshToken, refreshTokenExpireDate)

}

func SetJWSToken(payload structs.Payload, JWTToken string, tokenExpireDate time.Time) error {
	JWTKey := strconv.FormatInt(payload.PayloadId, 10) + strconv.Itoa(payload.UserId) + "JWT"

	return Set(JWTKey, JWTToken, tokenExpireDate)
}

func SetSession(sessionId, CSRFToken string, expireTime time.Time) error {
	return Set(sessionId, CSRFToken, expireTime)
}
