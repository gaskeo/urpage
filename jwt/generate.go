package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const minSecretKeySize = 32

var SecretKey = GenerateKey()

var ErrExpiredToken = errors.New("token has expired")
var ErrSmallSecretKey = errors.New("small secret key")
var ErrInvalidToken = errors.New("invalid token")

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

func GenerateJWTToken(id int, expiredAt time.Time, secretKey string) (Payload, string, error) {
	payload := Payload{id,
		GenerateId(),
		time.Now(),
		expiredAt}
	fmt.Println(secretKey)

	if len(secretKey) < minSecretKeySize {
		return Payload{}, "", ErrSmallSecretKey
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	JWTString, err := jwtToken.SignedString([]byte(secretKey))
	return payload, JWTString, err
}
