package jwt

import (
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-site/constants"
	"go-site/structs"
	"net/http"
	"time"
)

var SecretKey = GenerateKey()

func GenerateJWTToken(writer http.ResponseWriter, id int) (structs.Payload, string, time.Time, error) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)

	payloadId, err := GenerateId()

	if err != nil {
		return structs.Payload{}, "", time.Time{}, err
	}

	payload := structs.Payload{UserId: id,
		PayloadId: payloadId,
		IssuedAt:  time.Now(),
		ExpiredAt: tokenExpireDate}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	JWTString, err := jwtToken.SignedString([]byte(SecretKey))

	AddJWTCookie(writer, JWTString, tokenExpireDate)

	return payload, JWTString, tokenExpireDate, err
}

func GenerateRefreshToken(writer http.ResponseWriter, payload structs.Payload) (string, time.Time, error) {
	expireDate := time.Now().Add(constants.RefreshTokenExpireTime)
	refreshTokenByte := make([]byte, 16)
	_, err := rand.Read(refreshTokenByte)

	if err != nil {
		return "", time.Now(), err
	}

	refreshToken := fmt.Sprintf("%x", refreshTokenByte)

	AddRefreshTokenCookie(writer, refreshToken, payload.PayloadId, payload.UserId, expireDate)

	return refreshToken, expireDate, nil
}
