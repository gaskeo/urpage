package jwt_api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
	"urpage/constants"
)

var SecretKey = GenerateKey()

var ErrExpiredToken = errors.New("token has expired")
var ErrInvalidToken = errors.New("invalid token")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

type Payload struct {
	UserId    int
	PayloadId int64
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (payload *Payload) Valid() error {

	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

func GenerateJWTToken(writer http.ResponseWriter, id int) (Payload, string, time.Time, error) {
	tokenExpireDate := time.Now().Add(constants.JWTExpireTime)

	payloadId, err := GenerateId()

	if err != nil {
		return Payload{}, "", time.Time{}, err
	}

	payload := Payload{UserId: id,
		PayloadId: payloadId,
		IssuedAt:  time.Now(),
		ExpiredAt: tokenExpireDate}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	JWTString, err := jwtToken.SignedString([]byte(SecretKey))

	AddJWTCookie(writer, JWTString, tokenExpireDate)

	return payload, JWTString, tokenExpireDate, err
}

func GenerateRefreshToken(writer http.ResponseWriter, payload Payload) (string, time.Time, error) {
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
